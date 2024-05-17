package clients

import (
	"github.com/nuvla/api-client-go"
	"github.com/nuvla/api-client-go/clients/resources"
	"github.com/nuvla/api-client-go/types"
)

type UserClient struct {
	api_client_go.NuvlaClient
	UserID    *types.NuvlaID
	SessionID *types.NuvlaID
}

func NewUserClient(endpoint string, insecure bool, debug bool) *UserClient {
	sessionOpts := api_client_go.DefaultSessionOpts()
	sessionOpts.Insecure = insecure
	sessionOpts.Debug = debug
	sessionOpts.Endpoint = endpoint

	return &UserClient{
		NuvlaClient: *api_client_go.NewNuvlaClient(nil, sessionOpts),
	}
}

func (c *UserClient) AddNuvlaEdge(data map[string]interface{}) (*types.NuvlaID, error) {
	return c.Add("nuvlabox", data)
}

func (c *UserClient) GetNuvlaEdge(id string, fields []string) (*types.NuvlaResource, error) {
	return c.Get(id, fields)
}

func (c *UserClient) AddCredential(data map[string]interface{}) (*types.NuvlaID, error) {
	return c.Add("credential", data)
}

func (c *UserClient) GetId() string {
	return c.UserID.Id
}

func (c *UserClient) GetType() resources.NuvlaResourceType {
	return resources.UserType
}

func (c *UserClient) GetResourceMap() (map[string]interface{}, error) {
	return nil, nil
}

func (c *UserClient) UpdateResource() error {
	return nil
}
