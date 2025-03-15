package event

import (
	"strconv"
)

const (
	// ExampleCreatedEventName is the name for example creation events
	ExampleCreatedEventName = "example.created"
	// ExampleUpdatedEventName is the name for example update events
	ExampleUpdatedEventName = "example.updated"
	// ExampleDeletedEventName is the name for example deletion events
	ExampleDeletedEventName = "example.deleted"
)

// ExampleCreatedPayload contains data for example creation events
type ExampleCreatedPayload struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Alias string `json:"alias"`
}

// ExampleCreatedEvent represents an example creation event
type ExampleCreatedEvent struct {
	BaseEvent
}

// NewExampleCreatedEvent creates a new example creation event
func NewExampleCreatedEvent(id int, name, alias string) ExampleCreatedEvent {
	payload := ExampleCreatedPayload{
		ID:    id,
		Name:  name,
		Alias: alias,
	}
	return ExampleCreatedEvent{
		BaseEvent: NewBaseEvent(ExampleCreatedEventName, strconv.Itoa(id), payload),
	}
}

// ExampleUpdatedPayload contains data for example update events
type ExampleUpdatedPayload struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Alias string `json:"alias"`
}

// ExampleUpdatedEvent represents an example update event
type ExampleUpdatedEvent struct {
	BaseEvent
}

// NewExampleUpdatedEvent creates a new example update event
func NewExampleUpdatedEvent(id int, name, alias string) ExampleUpdatedEvent {
	payload := ExampleUpdatedPayload{
		ID:    id,
		Name:  name,
		Alias: alias,
	}
	return ExampleUpdatedEvent{
		BaseEvent: NewBaseEvent(ExampleUpdatedEventName, strconv.Itoa(id), payload),
	}
}

// ExampleDeletedPayload contains data for example deletion events
type ExampleDeletedPayload struct {
	ID int `json:"id"`
}

// ExampleDeletedEvent represents an example deletion event
type ExampleDeletedEvent struct {
	BaseEvent
}

// NewExampleDeletedEvent creates a new example deletion event
func NewExampleDeletedEvent(id int) ExampleDeletedEvent {
	payload := ExampleDeletedPayload{
		ID: id,
	}
	return ExampleDeletedEvent{
		BaseEvent: NewBaseEvent(ExampleDeletedEventName, strconv.Itoa(id), payload),
	}
}
