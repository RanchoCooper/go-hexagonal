package application

import (
	"go-hexagonal/application/example"
	"go-hexagonal/domain/repo"
	"go-hexagonal/domain/service"
)

// Factory provides methods to create application use cases
type Factory struct {
	exampleService service.IExampleService
	converter      service.Converter
	txFactory      repo.TransactionFactory
}

// NewFactory creates a new application factory
func NewFactory(
	exampleService service.IExampleService,
	converter service.Converter,
	txFactory repo.TransactionFactory,
) *Factory {
	return &Factory{
		exampleService: exampleService,
		converter:      converter,
		txFactory:      txFactory,
	}
}

// CreateExampleUseCase returns a new create example use case
func (f *Factory) CreateExampleUseCase() *example.CreateUseCase {
	return example.NewCreateUseCase(f.exampleService, f.converter, f.txFactory)
}

// DeleteExampleUseCase returns a new delete example use case
func (f *Factory) DeleteExampleUseCase() *example.DeleteUseCase {
	return example.NewDeleteUseCase(f.exampleService, f.txFactory)
}

// UpdateExampleUseCase returns a new update example use case
func (f *Factory) UpdateExampleUseCase() *example.UpdateUseCase {
	return example.NewUpdateUseCase(f.exampleService, f.converter, f.txFactory)
}

// GetExampleUseCase returns a new get example use case
func (f *Factory) GetExampleUseCase() *example.GetUseCase {
	return example.NewGetUseCase(f.exampleService, f.converter, f.txFactory)
}

// FindExampleByNameUseCase returns a new find example by name use case
func (f *Factory) FindExampleByNameUseCase() *example.FindByNameUseCase {
	return example.NewFindByNameUseCase(f.exampleService, f.converter, f.txFactory)
}
