package parallel

import (
	"sync"
)

type StreamSession struct {
	inputChan       chan<- interface{}
	inputChanClosed bool

	receiveChan       <-chan *StreamPayload
	receiveChanClosed chan struct{}
	receivedPayloads  []*StreamPayload

	m sync.Mutex
}

func (s *StreamSession) initAutoReceive() {
	go func() {
		for v := range s.receiveChan {
			s.receivedPayloads = append(s.receivedPayloads, v)
		}
		close(s.receiveChanClosed)
	}()
}

// Send an item to workers
func (s *StreamSession) Send(item interface{}) {
	if s.inputChanClosed {
		panic("you can not send item to a closed session")
	}
	s.inputChan <- item
}

// Wait all workers completed
func (s *StreamSession) Wait() {
	s.CompleteSend()
	for range s.receiveChanClosed {
	}
	return
}

// ReceivedPayloads return all receive payload which are auto received
func (s *StreamSession) ReceivedPayloads() []*StreamPayload {
	s.Wait()
	return s.receivedPayloads
}

// ReceiveChan return the raw chan, it can use to receive payload by yourself
func (s *StreamSession) ReceiveChan() <-chan *StreamPayload {
	return s.receiveChan
}

// CompleteSend indicate the workers that no more item to send
func (s *StreamSession) CompleteSend() {
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
