package model

import (
	"time"
)

// Example represents a basic example entity
type Example struct {
	Id        int           `json:"id"`
	Name      string        `json:"name"`
	Alias     string        `json:"alias"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
	events    []DomainEvent // Track domain events
}

// DomainEvent represents a domain event interface
type DomainEvent interface {
	EventType() string
}

// NewExample creates a new Example entity with validation
func NewExample(name, alias string) (*Example, error) {
	if name == "" {
		return nil, ErrEmptyExampleName
	}

	example := &Example{
		Name:      name,
		Alias:     alias,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		events:    make([]DomainEvent, 0),
	}

	// Record creation event
	example.addEvent(NewExampleCreatedEvent(example))

	return example, nil
}

// Validate ensures the Example entity meets domain rules
func (e *Example) Validate() error {
	if e.Name == "" {
		return ErrEmptyExampleName
	}
	if e.Id < 0 {
		return ErrInvalidExampleID
	}
	return nil
}

// Update changes the Example entity with validation
func (e *Example) Update(name, alias string) error {
	if name == "" {
		return ErrEmptyExampleName
	}

	e.Name = name
	e.Alias = alias
	e.UpdatedAt = time.Now()

	// Record update event
	e.addEvent(NewExampleUpdatedEvent(e))

	return nil
}

// MarkDeleted marks the entity as deleted and records a deletion event
func (e *Example) MarkDeleted() {
	e.addEvent(NewExampleDeletedEvent(e))
}

// Events returns all accumulated domain events and clears the event list
func (e *Example) Events() []DomainEvent {
	events := e.events
	e.events = make([]DomainEvent, 0)
	return events
}

// addEvent adds a domain event to the entity
func (e *Example) addEvent(event DomainEvent) {
	e.events = append(e.events, event)
}

// TableName returns the table name for the Example model
// This is kept for persistence adapters but is not part of domain logic
func (e Example) TableName() string {
	return "example"
}

// Domain events for Example entity

// ExampleCreatedEvent represents the creation of an example
type ExampleCreatedEvent struct {
	ExampleID int
	Name      string
	Alias     string
	Timestamp time.Time
}

// EventType returns the event type
func (e ExampleCreatedEvent) EventType() string {
	return "example.created"
}

// NewExampleCreatedEvent creates a new example created event
func NewExampleCreatedEvent(example *Example) ExampleCreatedEvent {
	return ExampleCreatedEvent{
		ExampleID: example.Id,
		Name:      example.Name,
		Alias:     example.Alias,
		Timestamp: time.Now(),
	}
}

// ExampleUpdatedEvent represents an update to an example
type ExampleUpdatedEvent struct {
	ExampleID int
	Name      string
	Alias     string
	Timestamp time.Time
}

// EventType returns the event type
func (e ExampleUpdatedEvent) EventType() string {
	return "example.updated"
}

// NewExampleUpdatedEvent creates a new example updated event
func NewExampleUpdatedEvent(example *Example) ExampleUpdatedEvent {
	return ExampleUpdatedEvent{
		ExampleID: example.Id,
		Name:      example.Name,
		Alias:     example.Alias,
		Timestamp: time.Now(),
	}
}

// ExampleDeletedEvent represents the deletion of an example
type ExampleDeletedEvent struct {
	ExampleID int
	Timestamp time.Time
}

// EventType returns the event type
func (e ExampleDeletedEvent) EventType() string {
	return "example.deleted"
}

// NewExampleDeletedEvent creates a new example deleted event
func NewExampleDeletedEvent(example *Example) ExampleDeletedEvent {
	return ExampleDeletedEvent{
		ExampleID: example.Id,
		Timestamp: time.Now(),
	}
}
