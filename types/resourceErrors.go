package types

import (
	"fmt"
	"github.com/nuvla/api-client-go/clients/resources"
)

// ResourceNotFoundError represents a custom error when a resource is not found
type ResourceNotFoundError struct {
	ResourceType resources.NuvlaResourceType
	ResourceID   string
}

// Error method to implement the error interface
func (e ResourceNotFoundError) Error() string {
	return fmt.Sprintf("%s with ID %s not found", e.ResourceType, e.ResourceID)
}

// NewResourceNotFoundError creates a new ResourceNotFoundError
func NewResourceNotFoundError(resourceType resources.NuvlaResourceType, resourceID string) *ResourceNotFoundError {
	return &ResourceNotFoundError{
		ResourceType: resourceType,
		ResourceID:   resourceID,
	}
}

type ResourceCreationError struct {
	ResourceType resources.NuvlaResourceType
	ResourceData map[string]interface{}
}

func (e ResourceCreationError) Error() string {
	return fmt.Sprintf("Error creating %s resource with data %v", e.ResourceType, e.ResourceData)
}

func NewResourceCreationError(resourceType resources.NuvlaResourceType, resourceData map[string]interface{}) *ResourceCreationError {
	return &ResourceCreationError{
		ResourceType: resourceType,
		ResourceData: resourceData,
	}
}
