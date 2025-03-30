package example

import (
	"time"

	"go-hexagonal/application/core"
	"go-hexagonal/domain/model"
)

// Input DTOs

// CreateInput represents input for creating a new example
type CreateInput struct {
	core.BaseInput
	Name  string `json:"name" validate:"required"`
	Alias string `json:"alias"`
}

// Validate validates the create input
func (i *CreateInput) Validate() error {
	if i.Name == "" {
		return core.ValidationError("name is required", map[string]any{
			"name": "required",
		})
	}
	return nil
}

// UpdateInput represents input for updating an example
type UpdateInput struct {
	core.BaseInput
	ID    int    `json:"id" validate:"required"`
	Name  string `json:"name" validate:"required"`
	Alias string `json:"alias"`
}

// Validate validates the update input
func (i *UpdateInput) Validate() error {
	if i.ID <= 0 {
		return core.ValidationError("invalid ID", map[string]any{
			"id": "must be positive",
		})
	}
	if i.Name == "" {
		return core.ValidationError("name is required", map[string]any{
			"name": "required",
		})
	}
	return nil
}

// GetInput represents input for retrieving an example by ID
type GetInput struct {
	core.BaseInput
	ID int `json:"id" validate:"required"`
}

// Validate validates the get input
func (i *GetInput) Validate() error {
	if i.ID <= 0 {
		return core.ValidationError("invalid ID", map[string]any{
			"id": "must be positive",
		})
	}
	return nil
}

// DeleteInput represents input for deleting an example
type DeleteInput struct {
	core.BaseInput
	ID int `json:"id" validate:"required"`
}

// Validate validates the delete input
func (i *DeleteInput) Validate() error {
	if i.ID <= 0 {
		return core.ValidationError("invalid ID", map[string]any{
			"id": "must be positive",
		})
	}
	return nil
}

// FindByNameInput represents input for finding an example by name
type FindByNameInput struct {
	core.BaseInput
	Name string `json:"name" validate:"required"`
}

// Validate validates the find by name input
func (i *FindByNameInput) Validate() error {
	if i.Name == "" {
		return core.ValidationError("name is required", map[string]any{
			"name": "required",
		})
	}
	return nil
}

// Output DTOs

// ExampleOutput represents the output format for example entities
type ExampleOutput struct {
	core.BaseOutput
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Alias     string    `json:"alias"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// FromModel converts a domain model to an output DTO
func (o *ExampleOutput) FromModel(example *model.Example) {
	o.ID = example.Id
	o.Name = example.Name
	o.Alias = example.Alias
	o.CreatedAt = example.CreatedAt
	o.UpdatedAt = example.UpdatedAt
	o.Status = "success"
}

// NewExampleOutput creates a new example output from a model
func NewExampleOutput(example *model.Example) *ExampleOutput {
	output := &ExampleOutput{}
	output.FromModel(example)
	return output
}
