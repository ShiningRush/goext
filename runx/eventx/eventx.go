package eventx

import (
	"context"
	"reflect"
	"sync"
	"sync/atomic"
)

var defaultBus EventBus = NewInMemoryEventBus()

// SetEventBus is used for unit test.
func SetEventBus(bus EventBus) {
	defaultBus = bus
}

// EventBus defines the event bus behavior.
type EventBus interface {
	Subscribe(handler EventHandler) error
	Publish(event Event)
	PublishSync(ctx context.Context, event Event)
	Close()
}

// EventHandler handles subscribed events.
type EventHandler interface {
	Topic() []string
	Handle(ctx context.Context, event Event)
}

// Event defines a publishable event.
type Event interface {
	Topic() []string
}

// NewInMemoryEventBus creates an in-memory event bus.
func NewInMemoryEventBus() *InMemoryEventBus {
	return &InMemoryEventBus{handlers: make(map[string][]EventHandler)}
}

// InMemoryEventBus dispatches events in memory.
type InMemoryEventBus struct {
	handlers map[string][]EventHandler
	mutex    sync.Mutex

	closed atomic.Bool
	wg     sync.WaitGroup
}

// Subscribe is thread-safe, but do not invoke Publish at the same time.
func (b *InMemoryEventBus) Subscribe(handler EventHandler) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	for _, tp := range handler.Topic() {
		b.handlers[tp] = append(b.handlers[tp], handler)
	}
	return nil
}

// Publish dispatches an event and does not wait for handlers to finish.
func (b *InMemoryEventBus) Publish(event Event) {
	if b.closed.Load() {
		panic("event bus is already closed")
	}

	matchedHandlers := b.getMatchHandlers(event.Topic())
	for i := range matchedHandlers {
		b.wg.Add(1)
		go func(idx int) {
			matchedHandlers[idx].Handle(context.TODO(), event)
			b.wg.Done()
		}(i)
	}
}

// PublishSync dispatches an event and waits for handlers to finish.
func (b *InMemoryEventBus) PublishSync(ctx context.Context, event Event) {
	if b.closed.Load() {
		panic("event bus is already closed")
	}

	matchedHandlers := b.getMatchHandlers(event.Topic())
	if len(matchedHandlers) > 0 {
		wg := sync.WaitGroup{}
		wg.Add(len(matchedHandlers))
		for i := range matchedHandlers {
			go func(idx int) {
				matchedHandlers[idx].Handle(ctx, event)
				wg.Done()
			}(i)
		}
		wg.Wait()
	}
}

// Close waits for all async goroutines to finish to prevent lost changes.
func (b *InMemoryEventBus) Close() {
	if !b.closed.Load() {
		b.closed.Store(true)
	}
	b.wg.Wait()
}

func (b *InMemoryEventBus) getMatchHandlers(topics []string) (matchedHandlers []EventHandler) {
	for _, tp := range topics {
		if h, ok := b.handlers[tp]; ok {
			matchedHandlers = append(matchedHandlers, h...)
		}
	}
	return
}

func (b *InMemoryEventBus) getEventTopic(eventType reflect.Type) reflect.Type {
	return eventType
}

// Subscribe is thread-safe, but do not call Publish at the same time.
func Subscribe(handler EventHandler) error {
	return defaultBus.Subscribe(handler)
}

// Publish dispatches an event and does not wait for handlers to finish.
func Publish(event Event) {
	defaultBus.Publish(event)
}

// PublishSync dispatches an event and waits for handlers to finish.
func PublishSync(ctx context.Context, event Event) {
	defaultBus.PublishSync(ctx, event)
}

// Close waits for all async goroutines to finish to prevent lost changes.
func Close() {
	defaultBus.Close()
}
