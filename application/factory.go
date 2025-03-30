package application

import (
	"go-hexagonal/application/example"
	"go-hexagonal/domain/repo"
	"go-hexagonal/domain/service"
)

// Factory provides methods to create application use cases
type Factory struct {
	exampleService service.IExampleService
	txFactory      repo.TransactionFactory
}

// NewFactory creates a new application factory
func NewFactory(
	exampleService service.IExampleService,
	txFactory repo.TransactionFactory,
) *Factory {
	return &Factory{
		exampleService: exampleService,
		txFactory:      txFactory,
	}
}

// CreateExampleUseCase returns a new create example use case
func (f *Factory) CreateExampleUseCase() *example.CreateUseCase {
	return example.NewCreateUseCase(f.exampleService, f.txFactory)
}

// DeleteExampleUseCase returns a new delete example use case
func (f *Factory) DeleteExampleUseCase() *example.DeleteUseCase {
	return example.NewDeleteUseCase(f.exampleService, f.txFactory)
}

// UpdateExampleUseCase returns a new update example use case
func (f *Factory) UpdateExampleUseCase() *example.UpdateUseCase {
	return example.NewUpdateUseCase(f.exampleService, f.txFactory)
}

// GetExampleUseCase returns a new get example use case
func (f *Factory) GetExampleUseCase() *example.GetUseCase {
	return example.NewGetUseCase(f.exampleService, f.txFactory)
}

// FindExampleByNameUseCase returns a new find example by name use case
func (f *Factory) FindExampleByNameUseCase() *example.FindByNameUseCase {
	return example.NewFindByNameUseCase(f.exampleService, f.txFactory)
}

// CreateExampleInput creates a new create example input
func (f *Factory) CreateExampleInput(name, alias string) *example.CreateInput {
	return &example.CreateInput{
		Name:  name,
		Alias: alias,
	}
}

// UpdateExampleInput creates a new update example input
func (f *Factory) UpdateExampleInput(id int, name, alias string) *example.UpdateInput {
	return &example.UpdateInput{
		ID:    id,
		Name:  name,
		Alias: alias,
	}
}

// GetExampleInput creates a new get example input
func (f *Factory) GetExampleInput(id int) *example.GetInput {
	return &example.GetInput{
		ID: id,
	}
}

// DeleteExampleInput creates a new delete example input
func (f *Factory) DeleteExampleInput(id int) *example.DeleteInput {
	return &example.DeleteInput{
		ID: id,
	}
}

// FindExampleByNameInput creates a new find example by name input
func (f *Factory) FindExampleByNameInput(name string) *example.FindByNameInput {
	return &example.FindByNameInput{
		Name: name,
	}
}
