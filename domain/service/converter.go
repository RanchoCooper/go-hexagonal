package service

import "go-hexagonal/domain/model"

// Converter defines an interface for converting between domain models and data transfer objects
type Converter interface {
	// ToExampleResponse converts a domain model to a response object
	ToExampleResponse(example *model.Example) (any, error)

	// FromCreateRequest converts a request object to a domain model
	FromCreateRequest(req any) (*model.Example, error)

	// FromUpdateRequest converts an update request to a domain model
	FromUpdateRequest(req any) (*model.Example, error)
}
