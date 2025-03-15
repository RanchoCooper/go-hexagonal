// Package event provides domain event functionality
package event

import (
	"time"

	"github.com/google/uuid"
)

// Event defines the domain event interface
type Event interface {
	// EventName returns the event name
	EventName() string
	// AggregateID returns the aggregate ID
	AggregateID() string
	// OccurredAt returns when the event occurred
	OccurredAt() time.Time
	// EventID returns the event ID
	EventID() string
}

// BaseEvent provides a base implementation for events
type BaseEvent struct {
	ID         string      `json:"id"`
	Name       string      `json:"name"`
	Aggregate  string      `json:"aggregate"`
	OccurredOn time.Time   `json:"occurred_on"`
	Payload    interface{} `json:"payload"`
}

// NewBaseEvent creates a new base event
func NewBaseEvent(name, aggregateID string, payload interface{}) BaseEvent {
	return BaseEvent{
		ID:         uuid.New().String(),
		Name:       name,
		Aggregate:  aggregateID,
		OccurredOn: time.Now(),
		Payload:    payload,
	}
}

// EventName returns the event name
func (e BaseEvent) EventName() string {
	return e.Name
}

// AggregateID returns the aggregate ID
func (e BaseEvent) AggregateID() string {
	return e.Aggregate
}

// OccurredAt returns when the event occurred
func (e BaseEvent) OccurredAt() time.Time {
	return e.OccurredOn
}

// EventID returns the event ID
func (e BaseEvent) EventID() string {
	return e.ID
}
