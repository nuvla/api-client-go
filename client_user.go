package api_client_go

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io"
)

type UserClient struct {
	Client    *NuvlaClient
	UserID    *NuvlaID
	SessionID *NuvlaID
}

func NewUserClient(endpoint string, insecure bool, debug bool) *UserClient {
	return &UserClient{
		Client: NewNuvlaClient(endpoint, insecure, debug),
	}
}

// Add creates a new resource of the given type and returns its ID
func (c *UserClient) Add(resourceType string, data map[string]interface{}) (*NuvlaID, error) {
	res, err := c.Client.Post(resourceType, data)
	if err != nil {
		log.Errorf("Error adding %s: %s", resourceType, err)
		return nil, err
	}
	var resData map[string]interface{}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Errorf("Error reading response body, cannot extract ID: %s", err)
		return nil, err
	}
	err = json.Unmarshal(bodyBytes, &resData)
	if err != nil {
		log.Errorf("Error unmarshaling response body, cannot extract ID: %s", err)
		return nil, err
	}
	log.Infof("ID of new %s: %s", resourceType, resData)

	return NewNuvlaIDFromId(resData["resource-id"].(string)), nil
}

type SearchOptions struct {
	First       int      `json:"first"`
	Filter      string   `json:"filter"`
	Fields      string   `json:"fields"`
	Select      []string `json:"select"`
	OrderBy     string   `json:"orderby"`
	Aggregation string   `json:"aggregation"`
}

func (c *UserClient) AddNuvlaEdge(data map[string]interface{}) (*NuvlaID, error) {
	return c.Add("nuvlabox", data)
}

func (c *UserClient) GetNuvlaEdge(id string, fields []string) (*NuvlaResource, error) {
	return c.Client.Get(id, fields)
}

func (c *UserClient) AddCredential(data map[string]interface{}) (*NuvlaID, error) {
	return c.Add("credential", data)
}
