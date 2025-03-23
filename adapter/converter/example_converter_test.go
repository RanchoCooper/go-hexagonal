package converter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go-hexagonal/api/dto"
	"go-hexagonal/domain/model"
)

func TestNewExampleConverter(t *testing.T) {
	// Test that NewExampleConverter returns a non-nil converter
	converter := NewExampleConverter()
	assert.NotNil(t, converter, "Converter should not be nil")
	assert.IsType(t, &ExampleConverter{}, converter, "Converter should be of type *ExampleConverter")
}

func TestExampleConverter_ToExampleResponse(t *testing.T) {
	// Create a converter
	converter := NewExampleConverter()

	t.Run("Valid Conversion", func(t *testing.T) {
		// Create test data
		now := time.Now()
		example := &model.Example{
			Id:        123,
			Name:      "Test Example",
			Alias:     "test",
			CreatedAt: now,
			UpdatedAt: now,
		}

		// Convert to response
		resp, err := converter.ToExampleResponse(example)

		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, resp)

		// Check type and values
		typedResp, ok := resp.(*dto.CreateExampleResp)
		assert.True(t, ok, "Response should be of type *dto.CreateExampleResp")
		assert.Equal(t, uint(123), typedResp.Id)
		assert.Equal(t, "Test Example", typedResp.Name)
		assert.Equal(t, "test", typedResp.Alias)
		assert.Equal(t, now, typedResp.CreatedAt)
		assert.Equal(t, now, typedResp.UpdatedAt)
	})

	t.Run("Nil Example", func(t *testing.T) {
		// Try to convert nil
		resp, err := converter.ToExampleResponse(nil)

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "example is nil")
	})
}

func TestExampleConverter_FromCreateRequest(t *testing.T) {
	// Create a converter
	converter := NewExampleConverter()

	t.Run("Valid Conversion", func(t *testing.T) {
		// Create a create request
		req := &dto.CreateExampleReq{
			Name:  "Test Example",
			Alias: "test",
		}

		// Convert to domain model
		model, err := converter.FromCreateRequest(req)

		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, model)
		assert.Equal(t, "Test Example", model.Name)
		assert.Equal(t, "test", model.Alias)
		assert.Zero(t, model.Id, "Id should not be set")
	})

	t.Run("Invalid Request Type", func(t *testing.T) {
		// Try to convert an invalid type
		model, err := converter.FromCreateRequest("invalid type")

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, model)
		assert.Contains(t, err.Error(), "invalid request type")
	})
}

func TestExampleConverter_FromUpdateRequest(t *testing.T) {
	// Create a converter
	converter := NewExampleConverter()

	t.Run("Valid Conversion", func(t *testing.T) {
		// Create an update request
		req := &dto.UpdateExampleReq{
			Id:    456,
			Name:  "Updated Example",
			Alias: "updated",
		}

		// Convert to domain model
		model, err := converter.FromUpdateRequest(req)

		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, model)
		assert.Equal(t, 456, model.Id)
		assert.Equal(t, "Updated Example", model.Name)
		assert.Equal(t, "updated", model.Alias)
	})

	t.Run("Invalid Request Type", func(t *testing.T) {
		// Try to convert an invalid type
		model, err := converter.FromUpdateRequest("invalid type")

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, model)
		assert.Contains(t, err.Error(), "invalid request type")
	})
}
