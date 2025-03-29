package model

import (
	stderrors "errors"

	"go-hexagonal/util/errors"
)

// Example related domain errors
var (
	// ErrExampleNotFound indicates the requested example was not found
	ErrExampleNotFound = errors.New(errors.ErrorTypeNotFound, "example not found")

	// ErrEmptyExampleName indicates an empty example name was provided
	ErrEmptyExampleName = errors.New(errors.ErrorTypeValidation, "example name cannot be empty")

	// ErrInvalidExampleID indicates an invalid example ID was provided
	ErrInvalidExampleID = errors.New(errors.ErrorTypeValidation, "invalid example ID")

	// ErrExampleNameTaken indicates an example with the given name already exists
	ErrExampleNameTaken = errors.New(errors.ErrorTypeConflict, "example name already taken")

	// ErrExampleInvalidUpdate indicates an attempt to update with invalid data
	ErrExampleInvalidUpdate = errors.New(errors.ErrorTypeValidation, "invalid example update data")

	// ErrExampleModified indicates the example was modified concurrently
	ErrExampleModified = errors.New(errors.ErrorTypeConflict, "example modified by another process")
)

// NewExampleNotFoundWithID creates a not found error with the example ID
func NewExampleNotFoundWithID(id int) *errors.AppError {
	return errors.Newf(errors.ErrorTypeNotFound, "example with ID %d not found", id)
}

// NewExampleNotFoundWithName creates a not found error with the example name
func NewExampleNotFoundWithName(name string) *errors.AppError {
	return errors.Newf(errors.ErrorTypeNotFound, "example with name '%s' not found", name)
}

// NewExampleNameTakenError creates an error indicating the name is already taken
func NewExampleNameTakenError(name string) *errors.AppError {
	return errors.Newf(errors.ErrorTypeConflict, "example with name '%s' already exists", name)
}

// IsExampleNotFoundError checks if the error indicates an example not found condition
func IsExampleNotFoundError(err error) bool {
	return stderrors.Is(err, ErrExampleNotFound) || errors.IsNotFoundError(err)
}

// IsExampleValidationError checks if the error is related to example validation
func IsExampleValidationError(err error) bool {
	return stderrors.Is(err, ErrEmptyExampleName) ||
		stderrors.Is(err, ErrInvalidExampleID) ||
		stderrors.Is(err, ErrExampleInvalidUpdate) ||
		errors.IsValidationError(err)
}

// IsExampleNameTakenError checks if the error indicates a name conflict
func IsExampleNameTakenError(err error) bool {
	return stderrors.Is(err, ErrExampleNameTaken)
}

// IsExampleModifiedError checks if the error indicates a concurrent modification
func IsExampleModifiedError(err error) bool {
	return stderrors.Is(err, ErrExampleModified)
}
