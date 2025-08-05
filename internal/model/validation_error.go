package model

import (
	"fmt"
	"strings"
)

// ValidationError represents a custom validation error with field-specific messages
type ValidationError struct {
	Field   string      `json:"field"`
	Message string      `json:"message"`
	Value   interface{} `json:"value,omitempty"`
}

type ValidationErrors struct {
	Errors []ValidationError `json:"errors"`
}

func (v ValidationErrors) Error() string {
	var messages []string
	for _, err := range v.Errors {
		messages = append(messages, fmt.Sprintf("%s: %s", err.Field, err.Message))
	}
	return strings.Join(messages, "; ")
}

// NewValidationError creates a new validation error
func NewValidationError(field, message string, value interface{}) ValidationError {
	return ValidationError{
		Field:   field,
		Message: message,
		Value:   value,
	}
}

// NewValidationErrors creates a new validation errors collection
func NewValidationErrors(errors ...ValidationError) ValidationErrors {
	return ValidationErrors{Errors: errors}
}

// AddError adds a validation error to the collection
func (v *ValidationErrors) AddError(field, message string, value interface{}) {
	v.Errors = append(v.Errors, NewValidationError(field, message, value))
}

// HasErrors returns true if there are validation errors
func (v *ValidationErrors) HasErrors() bool {
	return len(v.Errors) > 0
}
