package clients

import (
	"encoding/json"
	nuvla "github.com/nuvla/api-client-go"
	log "github.com/sirupsen/logrus"
)

type ClientResourceType string

const (
	DeploymentType ClientResourceType = "deployment"
	UserType       ClientResourceType = "user"
	NuvlaEdgeType  ClientResourceType = "nuvlaedge"
	NuvlaBoxType   ClientResourceType = "nuvlabox"
	JobType        ClientResourceType = "job"
)

type ClientResource interface {
	GetId() string
	GetType() ClientResourceType
	GetResourceMap() (map[string]interface{}, error)
	UpdateResource() error
}

func NewClientResource(resourceType string, resourceId string, client *nuvla.NuvlaClient) ClientResource {
	switch ClientResourceType(resourceType) {
	case DeploymentType:
		return NewNuvlaDeploymentClient(resourceId, client)
	case UserType:
		return &UserClient{}
	case NuvlaEdgeType:
		return &NuvlaEdgeClient{}
	case NuvlaBoxType:
		return &NuvlaEdgeClient{}
	case JobType:
		return &NuvlaJobClient{}
	default:
		return nil
	}
}

func MarshalResourceIntoMap(resource interface{}, resMap map[string]interface{}) error {
	var mapRes map[string]interface{}
	b, err := json.Marshal(resource)
	if err != nil {
		log.Error("Error marshaling resource to bytes")
		return err
	}

	err = json.Unmarshal(b, &mapRes)
	if err != nil {
		log.Error("Error unmarshaling resource to map")
		return err
	}

	log.Debug("Successfully marshaled resource into map")
	return nil
}

type Resource struct {
	// Required
	Id           string `json:"id"`
	ResourceType string `json:"resource-type"`
	Created      string `json:"created"`
	Updated      string `json:"updated"`

	// Optional
	Name        string `json:"name"`
	Description string `json:"description"`
}
