package parallel

import (
	"sync"
)

type StreamSession[P, R any] struct {
	inputChan       chan<- P
	inputChanClosed bool

	receiveChan       <-chan *StreamPayload[R]
	receiveChanClosed chan struct{}
	receivedPayloads  StreamPayloads[R]

	m sync.Mutex
}

func (s *StreamSession[P, R]) initAutoReceive() {
	go func() {
		for v := range s.receiveChan {
			s.receivedPayloads = append(s.receivedPayloads, v)
		}
		close(s.receiveChanClosed)
	}()
}

// Send an item to workers
func (s *StreamSession[P, R]) Send(item P) {
	if s.inputChanClosed {
		panic("you can not send item to a closed session")
	}
	s.inputChan <- item
}

// Wait all workers completed
func (s *StreamSession[P, R]) Wait() {
	s.CompleteSend()
	for range s.receiveChanClosed {
	}
	return
}

// ReceivedPayloads return all receive payload which are auto received
func (s *StreamSession[P, R]) ReceivedPayloads() StreamPayloads[R] {
	s.Wait()
	return s.receivedPayloads
}

// ReceiveChan return the raw chan, it can use to receive payload by yourself
func (s *StreamSession[P, R]) ReceiveChan() <-chan *StreamPayload[R] {
	return s.receiveChan
}

// CompleteSend indicate the workers that no more item to send
func (s *StreamSession[P, R]) CompleteSend() {
	if s.inputChanClosed {
		return
	}

	s.m.Lock()
	defer s.m.Unlock()

	// double check
	if s.inputChanClosed {
		return
	}
	close(s.inputChan)
	s.inputChanClosed = true
}
