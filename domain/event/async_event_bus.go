package event

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go-hexagonal/util/errors"
	"go-hexagonal/util/log"

	"go.uber.org/zap"
)

// EventStore defines the interface for persisting events
type EventStore interface {
	// SaveEvent persists an event
	SaveEvent(ctx context.Context, event Event) error
	// GetEvents retrieves events by type
	GetEvents(ctx context.Context, eventType string, since time.Time) ([]Event, error)
	// MarkProcessed marks an event as processed
	MarkProcessed(ctx context.Context, eventID string) error
}

// NoopEventStore is a no-operation event store
type NoopEventStore struct{}

// SaveEvent does nothing and returns nil
func (s *NoopEventStore) SaveEvent(ctx context.Context, event Event) error {
	return nil
}

// GetEvents returns an empty slice
func (s *NoopEventStore) GetEvents(ctx context.Context, eventType string, since time.Time) ([]Event, error) {
	return []Event{}, nil
}

// MarkProcessed does nothing and returns nil
func (s *NoopEventStore) MarkProcessed(ctx context.Context, eventID string) error {
	return nil
}

// AsyncEventBus implements an asynchronous event bus
type AsyncEventBus struct {
	handlers   []EventHandler
	store      EventStore
	mu         sync.RWMutex
	eventQueue chan Event
	workerPool chan struct{} // Semaphore for limiting concurrent workers
	quit       chan struct{}
	wg         sync.WaitGroup
}

// AsyncEventBusConfig holds configuration for AsyncEventBus
type AsyncEventBusConfig struct {
	QueueSize     int
	WorkerCount   int
	EventStore    EventStore
	ErrorCallback func(event Event, err error)
}

// DefaultAsyncEventBusConfig returns the default configuration
func DefaultAsyncEventBusConfig() *AsyncEventBusConfig {
	return &AsyncEventBusConfig{
		QueueSize:   100,
		WorkerCount: 5,
		EventStore:  &NoopEventStore{},
		ErrorCallback: func(event Event, err error) {
			logCtx := log.NewLogContext().
				WithComponent("AsyncEventBus").
				WithOperation("HandleEvent")

			logger, _ := log.New(
				log.WithLevel(log.ParseLogLevel("error")),
				log.WithCaller(true),
			)

			if logger != nil {
				logger.ErrorContext(logCtx, "Failed to process event",
					zap.String("event_type", event.EventName()),
					zap.String("event_id", event.EventID()),
					zap.Error(err),
				)
			}
		},
	}
}

// NewAsyncEventBus creates a new asynchronous event bus
func NewAsyncEventBus(config *AsyncEventBusConfig) *AsyncEventBus {
	if config == nil {
		config = DefaultAsyncEventBusConfig()
	}

	bus := &AsyncEventBus{
		handlers:   make([]EventHandler, 0),
		store:      config.EventStore,
		eventQueue: make(chan Event, config.QueueSize),
		workerPool: make(chan struct{}, config.WorkerCount),
		quit:       make(chan struct{}),
	}

	// Start workers
	bus.startWorkers(config.ErrorCallback)

	return bus
}

// startWorkers starts the worker goroutines
func (b *AsyncEventBus) startWorkers(errorCallback func(Event, error)) {
	b.wg.Add(1)
	go func() {
		defer b.wg.Done()
		for {
			select {
			case event := <-b.eventQueue:
				// Acquire semaphore slot
				b.workerPool <- struct{}{}

				// Process event in a new goroutine
				b.wg.Add(1)
				go func(evt Event) {
					defer b.wg.Done()
					defer func() { <-b.workerPool }() // Release semaphore slot

					// Process the event
					ctx := context.Background()
					b.mu.RLock()
					handlers := make([]EventHandler, len(b.handlers))
					copy(handlers, b.handlers) // Create a copy to avoid holding the lock
					b.mu.RUnlock()

					for _, handler := range handlers {
						if handler.InterestedIn(evt.EventName()) {
							if err := handler.HandleEvent(ctx, evt); err != nil {
								if errorCallback != nil {
									errorCallback(evt, err)
								}
							}
						}
					}

					// Mark event as processed in the store
					if b.store != nil {
						if err := b.store.MarkProcessed(ctx, evt.EventID()); err != nil {
							if errorCallback != nil {
								errorCallback(evt, errors.Wrapf(err, errors.ErrorTypePersistence, "failed to mark event as processed: %s", evt.EventID()))
							}
						}
					}
				}(event)

			case <-b.quit:
				return
			}
		}
	}()
}

// Publish publishes an event asynchronously
func (b *AsyncEventBus) Publish(ctx context.Context, event Event) error {
	// Persist the event first
	if b.store != nil {
		if err := b.store.SaveEvent(ctx, event); err != nil {
			return errors.Wrapf(err, errors.ErrorTypePersistence, "failed to save event: %s", event.EventID())
		}
	}

	// Send event to the queue
	select {
	case b.eventQueue <- event:
		return nil
	default:
		return errors.New(errors.ErrorTypeSystem, "event queue is full")
	}
}

// Subscribe registers an event handler
func (b *AsyncEventBus) Subscribe(handler EventHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.handlers = append(b.handlers, handler)
}

// Unsubscribe removes an event handler
func (b *AsyncEventBus) Unsubscribe(handler EventHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for i, h := range b.handlers {
		if h == handler {
			b.handlers = append(b.handlers[:i], b.handlers[i+1:]...)
			break
		}
	}
}

// Close shuts down the event bus gracefully
func (b *AsyncEventBus) Close(timeout time.Duration) error {
	// Signal workers to stop
	close(b.quit)

	// Wait for all workers to finish with timeout
	c := make(chan struct{})
	go func() {
		b.wg.Wait()
		close(c)
	}()

	select {
	case <-c:
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("timeout waiting for event bus to close")
	}
}

// ReplayEvents replays events from the store
func (b *AsyncEventBus) ReplayEvents(ctx context.Context, eventType string, since time.Time) error {
	if b.store == nil {
		return errors.New(errors.ErrorTypeSystem, "no event store configured")
	}

	events, err := b.store.GetEvents(ctx, eventType, since)
	if err != nil {
		return errors.Wrapf(err, errors.ErrorTypePersistence, "failed to get events for replay")
	}

	for _, event := range events {
		if err := b.Publish(ctx, event); err != nil {
			return errors.Wrapf(err, errors.ErrorTypeSystem, "failed to publish event during replay")
		}
	}

	return nil
}
