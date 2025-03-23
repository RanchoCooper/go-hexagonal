package event

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// MockEvent implements the Event interface for testing
type MockEvent struct {
	name        string
	aggregateID string
	occurredAt  time.Time
	eventID     string
}

func (e MockEvent) EventName() string {
	return e.name
}

func (e MockEvent) AggregateID() string {
	return e.aggregateID
}

func (e MockEvent) OccurredAt() time.Time {
	return e.occurredAt
}

func (e MockEvent) EventID() string {
	return e.eventID
}

// MockHandler implements the EventHandler interface for testing
type MockHandler struct {
	interestedEvents []string
	handleFunc       func(ctx context.Context, event Event) error
	handledEvents    []Event
}

func NewMockHandler(interestedEvents []string, handleFunc func(ctx context.Context, event Event) error) *MockHandler {
	if handleFunc == nil {
		// Default implementation if not provided
		handleFunc = func(ctx context.Context, event Event) error {
			return nil
		}
	}

	return &MockHandler{
		interestedEvents: interestedEvents,
		handleFunc:       handleFunc,
		handledEvents:    make([]Event, 0),
	}
}

func (h *MockHandler) HandleEvent(ctx context.Context, event Event) error {
	h.handledEvents = append(h.handledEvents, event)
	return h.handleFunc(ctx, event)
}

func (h *MockHandler) InterestedIn(eventName string) bool {
	for _, name := range h.interestedEvents {
		if name == eventName {
			return true
		}
	}
	return false
}

func TestNewNoopEventBus(t *testing.T) {
	bus := NewNoopEventBus()
	assert.NotNil(t, bus)
	assert.IsType(t, &NoopEventBus{}, bus)
}

func TestNoopEventBus_Publish(t *testing.T) {
	bus := NewNoopEventBus()

	// Create a mock event
	event := MockEvent{
		name:        "test.event",
		aggregateID: "123",
		occurredAt:  time.Now(),
		eventID:     "event-123",
	}

	// Publish should return nil for NoopEventBus
	err := bus.Publish(context.Background(), event)
	assert.NoError(t, err)
}

func TestNoopEventBus_Subscribe(t *testing.T) {
	bus := NewNoopEventBus()

	// Create a mock handler
	handler := NewMockHandler([]string{"test.event"}, nil)

	// Subscribe should not panic
	assert.NotPanics(t, func() {
		bus.Subscribe(handler)
	})
}

func TestNoopEventBus_Unsubscribe(t *testing.T) {
	bus := NewNoopEventBus()

	// Create a mock handler
	handler := NewMockHandler([]string{"test.event"}, nil)

	// Unsubscribe should not panic
	assert.NotPanics(t, func() {
		bus.Unsubscribe(handler)
	})
}

func TestNewInMemoryEventBus(t *testing.T) {
	bus := NewInMemoryEventBus()
	assert.NotNil(t, bus)
	assert.IsType(t, &InMemoryEventBus{}, bus)
	assert.Empty(t, bus.handlers)
}

func TestInMemoryEventBus_Subscribe(t *testing.T) {
	bus := NewInMemoryEventBus()

	// Create mock handlers
	handler1 := NewMockHandler([]string{"test.event1"}, nil)
	handler2 := NewMockHandler([]string{"test.event2"}, nil)

	// Subscribe handlers
	bus.Subscribe(handler1)
	bus.Subscribe(handler2)

	// Check that handlers were added
	assert.Len(t, bus.handlers, 2)
	assert.Contains(t, bus.handlers, handler1)
	assert.Contains(t, bus.handlers, handler2)
}

func TestInMemoryEventBus_Unsubscribe(t *testing.T) {
	bus := NewInMemoryEventBus()

	// Create mock handlers
	handler1 := NewMockHandler([]string{"test.event1"}, nil)
	handler2 := NewMockHandler([]string{"test.event2"}, nil)

	// Subscribe handlers
	bus.Subscribe(handler1)
	bus.Subscribe(handler2)

	// Verify initial state
	assert.Len(t, bus.handlers, 2)

	// Unsubscribe one handler
	bus.Unsubscribe(handler1)

	// Check that only one handler remains
	assert.Len(t, bus.handlers, 1)
	assert.NotContains(t, bus.handlers, handler1)
	assert.Contains(t, bus.handlers, handler2)

	// Unsubscribe the other handler
	bus.Unsubscribe(handler2)

	// Check that no handlers remain
	assert.Empty(t, bus.handlers)
}

func TestInMemoryEventBus_Publish(t *testing.T) {
	t.Run("No Handlers", func(t *testing.T) {
		bus := NewInMemoryEventBus()

		// Create a mock event
		event := MockEvent{
			name:        "test.event",
			aggregateID: "123",
			occurredAt:  time.Now(),
			eventID:     "event-123",
		}

		// Publish should not error with no handlers
		err := bus.Publish(context.Background(), event)
		assert.NoError(t, err)
	})

	t.Run("Interested Handlers", func(t *testing.T) {
		bus := NewInMemoryEventBus()

		// Create a mock event
		event := MockEvent{
			name:        "test.event",
			aggregateID: "123",
			occurredAt:  time.Now(),
			eventID:     "event-123",
		}

		// Create mock handlers
		interestedHandler := NewMockHandler([]string{"test.event"}, nil)
		notInterestedHandler := NewMockHandler([]string{"other.event"}, nil)

		// Subscribe handlers
		bus.Subscribe(interestedHandler)
		bus.Subscribe(notInterestedHandler)

		// Publish the event
		err := bus.Publish(context.Background(), event)
		assert.NoError(t, err)

		// Interested handler should have received the event
		assert.Len(t, interestedHandler.handledEvents, 1)
		assert.Equal(t, event, interestedHandler.handledEvents[0])

		// Not interested handler should not have received the event
		assert.Empty(t, notInterestedHandler.handledEvents)
	})

	t.Run("Handler Error", func(t *testing.T) {
		bus := NewInMemoryEventBus()

		// Create a mock event
		event := MockEvent{
			name:        "test.event",
			aggregateID: "123",
			occurredAt:  time.Now(),
			eventID:     "event-123",
		}

		// Create a handler that returns an error
		expectedErr := errors.New("handler error")
		errorHandler := NewMockHandler([]string{"test.event"}, func(ctx context.Context, event Event) error {
			return expectedErr
		})

		// Create a handler that would be called after the error handler
		laterHandler := NewMockHandler([]string{"test.event"}, nil)

		// Subscribe handlers
		bus.Subscribe(errorHandler)
		bus.Subscribe(laterHandler)

		// Publish the event
		err := bus.Publish(context.Background(), event)

		// Should return the error from the handler
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)

		// Error handler should have received the event
		assert.Len(t, errorHandler.handledEvents, 1)

		// Later handler should not be called due to the error
		assert.Empty(t, laterHandler.handledEvents)
	})

	t.Run("Multiple Interested Handlers", func(t *testing.T) {
		bus := NewInMemoryEventBus()

		// Create a mock event
		event := MockEvent{
			name:        "test.event",
			aggregateID: "123",
			occurredAt:  time.Now(),
			eventID:     "event-123",
		}

		// Create multiple interested handlers
		handler1 := NewMockHandler([]string{"test.event"}, nil)
		handler2 := NewMockHandler([]string{"test.event", "other.event"}, nil)
		handler3 := NewMockHandler([]string{"other.event"}, nil)

		// Subscribe handlers
		bus.Subscribe(handler1)
		bus.Subscribe(handler2)
		bus.Subscribe(handler3)

		// Publish the event
		err := bus.Publish(context.Background(), event)
		assert.NoError(t, err)

		// Both interested handlers should have received the event
		assert.Len(t, handler1.handledEvents, 1)
		assert.Equal(t, event, handler1.handledEvents[0])

		assert.Len(t, handler2.handledEvents, 1)
		assert.Equal(t, event, handler2.handledEvents[0])

		// Not interested handler should not have received the event
		assert.Empty(t, handler3.handledEvents)
	})
}
