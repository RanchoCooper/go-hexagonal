package event

import (
	"context"
	"encoding/json"

	"go-hexagonal/util/log"

	"go.uber.org/zap"
)

// LoggingEventHandler logs all events it receives
type LoggingEventHandler struct {
	interestedEvents []string
}

// NewLoggingEventHandler creates a new logging event handler
func NewLoggingEventHandler(events ...string) *LoggingEventHandler {
	return &LoggingEventHandler{
		interestedEvents: events,
	}
}

// HandleEvent logs the event details
func (h *LoggingEventHandler) HandleEvent(ctx context.Context, event Event) error {
	eventData, _ := json.Marshal(event)
	log.Logger.Info("Event received",
		zap.String("event_name", event.EventName()),
		zap.String("event_data", string(eventData)))
	return nil
}

// InterestedIn checks if the handler is interested in the event
func (h *LoggingEventHandler) InterestedIn(eventName string) bool {
	if len(h.interestedEvents) == 0 {
		return true
	}
	for _, name := range h.interestedEvents {
		if name == eventName {
			return true
		}
	}
	return false
}

// ExampleEventHandler handles example-related events
type ExampleEventHandler struct {
}

// NewExampleEventHandler creates a new example event handler
func NewExampleEventHandler() *ExampleEventHandler {
	return &ExampleEventHandler{}
}

// HandleEvent handles the example event based on its type
func (h *ExampleEventHandler) HandleEvent(ctx context.Context, event Event) error {
	switch event.EventName() {
	case ExampleCreatedEventName:
		return h.handleExampleCreated(ctx, event)
	case ExampleUpdatedEventName:
		return h.handleExampleUpdated(ctx, event)
	case ExampleDeletedEventName:
		return h.handleExampleDeleted(ctx, event)
	default:
		return nil
	}
}

// InterestedIn checks if the handler is interested in the event
func (h *ExampleEventHandler) InterestedIn(eventName string) bool {
	return eventName == ExampleCreatedEventName ||
		eventName == ExampleUpdatedEventName ||
		eventName == ExampleDeletedEventName
}

// handleExampleCreated handles example creation events
func (h *ExampleEventHandler) handleExampleCreated(ctx context.Context, event Event) error {
	log.Logger.Info("Example created",
		zap.String("id", event.AggregateID()),
		zap.String("event_id", event.EventID()))
	return nil
}

// handleExampleUpdated handles example update events
func (h *ExampleEventHandler) handleExampleUpdated(ctx context.Context, event Event) error {
	log.Logger.Info("Example updated",
		zap.String("id", event.AggregateID()),
		zap.String("event_id", event.EventID()))
	return nil
}

// handleExampleDeleted handles example deletion events
func (h *ExampleEventHandler) handleExampleDeleted(ctx context.Context, event Event) error {
	log.Logger.Info("Example deleted",
		zap.String("id", event.AggregateID()),
		zap.String("event_id", event.EventID()))
	return nil
}
