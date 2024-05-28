package resources

import (
	"encoding/json"
	"time"
)

type NuvlaResource interface {
	GetId() string
	GetType() string
	New() NuvlaResource
}

type CommonAttributesResource struct {
	Id           string    `json:"id,omitempty"`
	ResourceType string    `json:"resource-type,omitempty"`
	Created      time.Time `json:"created,omitempty"`
	Updated      time.Time `json:"updated,omitempty"`

	// Optional fields
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	Parent      string   `json:"parent,omitempty"`
}

func (r *CommonAttributesResource) GetId() string {
	return r.Id
}

func (r *CommonAttributesResource) GetType() string {
	return r.ResourceType
}

func NewResourceFromMap(resMap map[string]interface{}, resource NuvlaResource) error {
	// Unmarshal map into resource
	jsonResource, err := json.Marshal(resMap)
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonResource, resource)
	if err != nil {
		return err
	}
	return nil
}
