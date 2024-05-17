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
