package parallel

import (
	"runtime"
	"sync"

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

type StreamPayload struct {
	Result interface{}
	Err    error
}

func (o *StreamPayload) HasError() bool {
	return o.Err != nil
}

type StreamWork func(workerIdx int, item interface{}) (ret interface{}, err error)

func StreamDo(work StreamWork, ops ...OptionOp) *StreamSession {
	opt := InitialOption(ops)

	inputChan := make(chan interface{})

	workerChanList, receiveChan := initWorkers(work, opt)
	dispatch(inputChan, workerChanList)

	s := &StreamSession{
		inputChan:         inputChan,
		receiveChan:       receiveChan,
		receiveChanClosed: make(chan struct{}),
	}
	runtime.SetFinalizer(s, ensureFreeSession)
	if !opt.receiveDataExplicit {
		s.initAutoReceive()
	}
	return s
}
func ensureFreeSession(session *StreamSession) {
	session.CompleteSend()
}

func dispatch(inputChan <-chan interface{}, workerChanList []chan<- interface{}) {
	go func() {
		// if input could have identity, murmur3 would be better a solution than atomic increment, but we can not require that.
		var autoIncrement uint64
		for input := range inputChan {
			if k, ok := input.(KeyOwner); ok {
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

func initWorkers(work StreamWork, opt *Option) (workerChanList []chan<- interface{}, retChan <-chan *StreamPayload) {
	rawRetChan := make(chan *StreamPayload)
	retChan = rawRetChan

	var wg sync.WaitGroup
	wg.Add(opt.workerNumber)
	for i := 0; i < opt.workerNumber; i++ {
		// we make each worker have own queue to instead of sharing same work queue, it can avoid chan race
		workerCh := make(chan interface{})
		workerChanList = append(workerChanList, workerCh)

		go func(idx int, inputCh <-chan interface{}) {
			for input := range inputCh {
				ret, err := work(idx, input)
				if !opt.ignoreResult {
					rawRetChan <- &StreamPayload{
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
