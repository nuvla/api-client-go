package errors

import "fmt"

// ResourceNotFoundError represents a custom error when a resource is not found
type ResourceNotFoundError struct {
	ResourceType string
	ResourceID   string
}

// Error method to implement the error interface
func (e *ResourceNotFoundError) Error() string {
	return fmt.Sprintf("%s with ID %s not found", e.ResourceType, e.ResourceID)
}

// NewResourceNotFoundError creates a new ResourceNotFoundError
func NewResourceNotFoundError(resourceType, resourceID string) *ResourceNotFoundError {
	return &ResourceNotFoundError{
		ResourceType: resourceType,
		ResourceID:   resourceID,
	}
}
