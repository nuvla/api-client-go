package resources

import (
	"time"
)

type NuvlaResource interface {
	GetId() string
	GetType() string
}

type CommonAttributesResource struct {
	Id           string    `json:"id"`
	ResourceType string    `json:"resource-type"`
	Created      time.Time `json:"created"`
	Updated      time.Time `json:"updated"`
	// Optional fields
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	Parent      string   `json:"parent"`
}

func (r *CommonAttributesResource) GetId() string {
	return r.Id
}

func (r *CommonAttributesResource) GetType() string {
	return r.ResourceType
}
