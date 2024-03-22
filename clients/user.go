package clients

import (
	"encoding/json"
	"github.com/nuvla/api-client-go"
	"github.com/nuvla/api-client-go/types"
	log "github.com/sirupsen/logrus"
	"io"
)

type UserClient struct {
	Client    *api_client_go.NuvlaClient
	UserID    *types.NuvlaID
	SessionID *types.NuvlaID
}

func NewUserClient(endpoint string, insecure bool, debug bool) *UserClient {
	return nil
}

// Add creates a new resource of the given type and returns its ID
func (c *UserClient) Add(resourceType string, data map[string]interface{}) (*types.NuvlaID, error) {
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

	return types.NewNuvlaIDFromId(resData["resource-id"].(string)), nil
}

type SearchOptions struct {
	First       int      `json:"first"`
	Filter      string   `json:"filter"`
	Fields      string   `json:"fields"`
	Select      []string `json:"select"`
	OrderBy     string   `json:"orderby"`
	Aggregation string   `json:"aggregation"`
}

func (c *UserClient) AddNuvlaEdge(data map[string]interface{}) (*types.NuvlaID, error) {
	return c.Add("nuvlabox", data)
}

func (c *UserClient) GetNuvlaEdge(id string, fields []string) (*types.NuvlaResource, error) {
	return c.Client.Get(id, fields)
}

func (c *UserClient) AddCredential(data map[string]interface{}) (*types.NuvlaID, error) {
	return c.Add("credential", data)
}

func (c *UserClient) GetId() string {
	return c.UserID.Id
}

func (c *UserClient) GetType() ClientResourceType {
	return UserType
}

func (c *UserClient) GetResourceMap() (map[string]interface{}, error) {
	return nil, nil
}

func (c *UserClient) UpdateResource() error {
	return nil
}
