package parallel

import (
	"runtime"
	"sync"

	"github.com/shiningrush/goext/errx"
	"github.com/twmb/murmur3"
)

type Work func(workerIdx int) error

func Do(work Work, ops ...OptionOp) (errs []error) {
	opt := InitialOption(ops)

	var wg sync.WaitGroup
	wg.Add(opt.workerNumber)

	errCh := make(chan error)
	go func() {
		wg.Wait()
		close(errCh)
	}()
	for i := 0; i < opt.workerNumber; i++ {
		go func(workerIdx int) {
			if err := work(workerIdx); err != nil {
				errCh <- err
			}
			wg.Done()
		}(i)
	}

	for err := range errCh {
		errs = append(errs, err)
	}

	wg.Wait()
	return
}

type StreamPayloads[R any] []*StreamPayload[R]

func (p StreamPayloads[R]) AggregateErr() error {
	e := &errx.BatchErrors{}
	for _, v := range p {
		if v.HasError() {
			e.Append(v.Err)
		}
	}

	if e.HasError() {
		return e
	}
	return nil
}

type StreamPayload[R any] struct {
	Result R
	Err    error
}

func (o *StreamPayload[R]) HasError() bool {
	return o.Err != nil
}

type StreamWork[P, R any] func(workerIdx int, item P) (ret R, err error)

func StreamDo[P, R any](work StreamWork[P, R], ops ...OptionOp) *StreamSession[P, R] {
	opt := InitialOption(ops)

	inputChan := make(chan P)

	workerChanList, receiveChan := initWorkers(work, opt)
	dispatch[P](inputChan, workerChanList)

	s := &StreamSession[P, R]{
		inputChan:         inputChan,
		receiveChan:       receiveChan,
		receiveChanClosed: make(chan struct{}),
	}
	runtime.SetFinalizer(s, ensureFreeSession[P, R])
	if !opt.receiveDataExplicit {
		s.initAutoReceive()
	}
	return s
}
func ensureFreeSession[P, R any](session *StreamSession[P, R]) {
	session.CompleteSend()
}

func dispatch[P any](inputChan <-chan P, workerChanList []chan<- P) {
	go func() {
		// if input could have identity, murmur3 would be better a solution than atomic increment, but we can not require that.
		var autoIncrement uint64
		for input := range inputChan {
			if k, ok := (interface{})(input).(KeyOwner); ok {
				mur := murmur3.New32()
				// write return no error
				_, _ = mur.Write([]byte(k.GetKey()))
				workerChanList[int(mur.Sum32())%len(workerChanList)] <- input
				continue
			}

			workerChanList[autoIncrement%uint64(len(workerChanList))] <- input
			autoIncrement++
		}
		for _, v := range workerChanList {
			close(v)
		}
	}()
}

func initWorkers[P, R any](work StreamWork[P, R], opt *Option) (workerChanList []chan<- P, retChan <-chan *StreamPayload[R]) {
	rawRetChan := make(chan *StreamPayload[R])
	retChan = rawRetChan

	var wg sync.WaitGroup
	wg.Add(opt.workerNumber)
	for i := 0; i < opt.workerNumber; i++ {
		// we make each worker have own queue to instead of sharing same work queue, it can avoid chan race
		workerCh := make(chan P)
		workerChanList = append(workerChanList, workerCh)

		go func(idx int, inputCh <-chan P) {
			for input := range inputCh {
				ret, err := work(idx, input)
				if !opt.ignoreResult {
					rawRetChan <- &StreamPayload[R]{
						Result: ret,
						Err:    err,
					}
				}
			}
			wg.Done()
		}(i, workerCh)
	}
	go func() {
		wg.Wait()
		close(rawRetChan)
	}()
	return
}
