package event

import (
	"context"
	"sync"
)

// EventHandler defines the event handler interface
type EventHandler interface {
	// HandleEvent processes an event
	HandleEvent(ctx context.Context, event Event) error
	// InterestedIn checks if the handler is interested in the event
	InterestedIn(eventName string) bool
}

// EventBus defines the event bus interface
type EventBus interface {
	// Publish publishes an event
	Publish(ctx context.Context, event Event) error
	// Subscribe registers an event handler
	Subscribe(handler EventHandler)
	// Unsubscribe removes an event handler
	Unsubscribe(handler EventHandler)
}

// NoopEventBus implements a no-operation event bus
type NoopEventBus struct{}

// NewNoopEventBus creates a new no-operation event bus
func NewNoopEventBus() *NoopEventBus {
	return &NoopEventBus{}
}

// Publish does nothing and returns nil
func (b *NoopEventBus) Publish(ctx context.Context, event Event) error {
	return nil
}

// Subscribe does nothing
func (b *NoopEventBus) Subscribe(handler EventHandler) {}

// Unsubscribe does nothing
func (b *NoopEventBus) Unsubscribe(handler EventHandler) {}

// InMemoryEventBus implements an in-memory event bus
type InMemoryEventBus struct {
	handlers []EventHandler
	mu       sync.RWMutex
}

// NewInMemoryEventBus creates a new in-memory event bus
func NewInMemoryEventBus() *InMemoryEventBus {
	return &InMemoryEventBus{
		handlers: make([]EventHandler, 0),
	}
}

// Publish publishes an event to all interested handlers
func (b *InMemoryEventBus) Publish(ctx context.Context, event Event) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for _, handler := range b.handlers {
		if handler.InterestedIn(event.EventName()) {
			if err := handler.HandleEvent(ctx, event); err != nil {
				return err
			}
		}
	}
	return nil
}

// Subscribe registers an event handler
func (b *InMemoryEventBus) Subscribe(handler EventHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.handlers = append(b.handlers, handler)
}

// Unsubscribe removes an event handler
func (b *InMemoryEventBus) Unsubscribe(handler EventHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for i, h := range b.handlers {
		if h == handler {
			b.handlers = append(b.handlers[:i], b.handlers[i+1:]...)
			break
		}
	}
}
