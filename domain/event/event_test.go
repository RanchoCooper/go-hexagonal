package event

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewBaseEvent(t *testing.T) {
	// Test data
	name := "test.event"
	aggregateID := "123"
	payload := map[string]string{"key": "value"}

	// Create a new base event
	event := NewBaseEvent(name, aggregateID, payload)

	// Assertions
	assert.Equal(t, name, event.Name)
	assert.Equal(t, aggregateID, event.Aggregate)
	assert.Equal(t, payload, event.Payload)
	assert.NotEmpty(t, event.ID)
	assert.NotZero(t, event.OccurredOn)

	// Time should be close to now
	assert.WithinDuration(t, time.Now(), event.OccurredOn, 2*time.Second)
}

func TestBaseEvent_EventName(t *testing.T) {
	// Create a test event
	event := BaseEvent{
		Name: "test.event",
	}

	// Test EventName method
	assert.Equal(t, "test.event", event.EventName())
}

func TestBaseEvent_AggregateID(t *testing.T) {
	// Create a test event
	event := BaseEvent{
		Aggregate: "123",
	}

	// Test AggregateID method
	assert.Equal(t, "123", event.AggregateID())
}

func TestBaseEvent_OccurredAt(t *testing.T) {
	// Create a test event with specific time
	now := time.Now()
	event := BaseEvent{
		OccurredOn: now,
	}

	// Test OccurredAt method
	assert.Equal(t, now, event.OccurredAt())
}

func TestBaseEvent_EventID(t *testing.T) {
	// Create a test event
	event := BaseEvent{
		ID: "event-123",
	}

	// Test EventID method
	assert.Equal(t, "event-123", event.EventID())
}
