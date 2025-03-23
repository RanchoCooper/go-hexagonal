package converter

import (
	"fmt"

	"go-hexagonal/api/dto"
	"go-hexagonal/domain/model"
	"go-hexagonal/domain/service"
)

// ExampleConverter implements the service.Converter interface for Example entities
type ExampleConverter struct{}

// NewExampleConverter creates a new ExampleConverter
func NewExampleConverter() service.Converter {
	return &ExampleConverter{}
}

// ToExampleResponse converts a domain model to a response object
func (c *ExampleConverter) ToExampleResponse(example *model.Example) (any, error) {
	if example == nil {
		return nil, fmt.Errorf("example is nil")
	}

	return &dto.CreateExampleResp{
		Id:        uint(example.Id),
		Name:      example.Name,
		Alias:     example.Alias,
		CreatedAt: example.CreatedAt,
		UpdatedAt: example.UpdatedAt,
	}, nil
}

// FromCreateRequest converts a create request to a domain model
func (c *ExampleConverter) FromCreateRequest(req any) (*model.Example, error) {
	createReq, ok := req.(*dto.CreateExampleReq)
	if !ok {
		return nil, fmt.Errorf("invalid request type: %T", req)
	}

	return &model.Example{
		Name:  createReq.Name,
		Alias: createReq.Alias,
	}, nil
}

// FromUpdateRequest converts an update request to a domain model
func (c *ExampleConverter) FromUpdateRequest(req any) (*model.Example, error) {
	updateReq, ok := req.(*dto.UpdateExampleReq)
	if !ok {
		return nil, fmt.Errorf("invalid request type: %T", req)
	}

	return &model.Example{
		Id:    int(updateReq.Id),
		Name:  updateReq.Name,
		Alias: updateReq.Alias,
	}, nil
}
