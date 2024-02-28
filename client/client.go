package client

import (
	"api-client-go/client/types"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Client interface {
	// LoginApiKeys Credentials management
	LoginApiKeys(key string, secret string) error
	LoginUser(username string, password string) error
	Logout() error
	IsAuthenticated() bool

	cimiRequest(method string, endpointPath string, data interface{}) (interface{}, error)

	// Nuvla API base methods
	Get(endpointPath string) (interface{}, error)
	Post(data interface{}, endpointPath string) (interface{}, error)
	Put(data interface{}, endPointPath string, toDelete []string) (interface{}, error)

	// Resource management methods
	// Add, Edit, Delete, DeleteBulk, OperationBulk.
	Add(endpointPath string, data interface{}) (interface{}, error)
	Edit(endpointPath string, data interface{}) (interface{}, error)
	Delete(endpointPath string) (interface{}, error)
	DeleteBulk(endpointPath string, data interface{}) (interface{}, error)
	OperationBulk(endpointPath string, data interface{}) (interface{}, error)
	Search(endpointPath string, data interface{}) (interface{}, error)
	Operation(endpointPath string, data interface{}) (interface{}, error)
}

type NuvlaClient struct {
	nuvlaEndpoint    string
	debug            bool
	loginCredentials map[string]string

	// Session params
	session       *NuvlaSession
	nuvlaInsecure bool
	compress      bool
}

func NewNuvlaClient(attrs *types.SessionAttributes) *NuvlaClient {

	nc := &NuvlaClient{
		nuvlaEndpoint: attrs.Endpoint,
		nuvlaInsecure: attrs.Insecure,
		debug:         attrs.Debug,
		session:       NewNuvlaSession(attrs),
		compress:      false,
	}
	return nc
}

func (nc *NuvlaClient) LoginApiKeys(key string, secret string) error {
	err := nc.session.login(types.NewApiKeyLogInParams(key, secret))
	if err != nil {
		log.Errorf("Error logging in with api keys: %s", err)
		return err
	}
	return nil
}

func (nc *NuvlaClient) LoginUser(username string, password string) error {
	err := nc.session.login(types.NewUserLogInParams(username, password))
	if err != nil {
		log.Errorf("Error logging in with user credentials: %s", err)
		return err
	}
	return nil
}

func (nc *NuvlaClient) Logout() error {
	return nc.session.logout()
}

func (nc *NuvlaClient) IsAuthenticated() bool {
	return true
}

func (nc *NuvlaClient) buildUriEndPoint(uriEndpoint string) string {
	return fmt.Sprintf("%s/api/%s", nc.nuvlaEndpoint, uriEndpoint)
}

func (nc *NuvlaClient) buildOperationUriEndPoint(uriEndpoint string, operation string) string {
	return fmt.Sprintf("%s/%s", nc.buildUriEndPoint(uriEndpoint), operation)
}

func (nc *NuvlaClient) cimiRequest(reqInput *types.RequestInput) (*http.Response, error) {
	// Setup default client headers for all requests
	// TODO: Might be configurable from session
	if reqInput.Headers == nil {
		reqInput.Headers = make(map[string]string)
	}
	reqInput.Headers["Accept"] = "application/json"
	reqInput.Headers["Accept-Encoding"] = "gzip"

	if reqInput.JsonData != nil {
		reqInput.Headers["Content-Type"] = "application/json"
	}

	r, _ := nc.session.Request(reqInput)

	return r, nil
}

// Get executes the get http method
// Allow for selective fields to be returned via the selectFields parameter

func (nc *NuvlaClient) Get(resourceId string, selectFields []string) (*types.NuvlaResource, error) {
	// Define request inputs to allow adding select fields
	r := &types.RequestInput{
		Method:   "GET",
		Endpoint: nc.buildUriEndPoint(resourceId),
	}

	// Do not create the request params struct unless we need it to prevent overhead down the line
	if selectFields != nil {
		r.Params = &types.RequestParams{
			Select: selectFields,
		}
	}

	resp, err := nc.cimiRequest(r)
	if err != nil {
		log.Errorf("Error executing GET request: %s", err)
		return nil, err
	}

	return types.NewResourceFromResponse(resp), nil
}

// Post executes the post http method
// Data can be any type, but it will be marshaled into JSON
func (nc *NuvlaClient) Post(endpoint string, data map[string]interface{}) (*http.Response, error) {
	r := &types.RequestInput{
		Method:   "POST",
		JsonData: data,
		Endpoint: nc.buildUriEndPoint(endpoint),
	}

	resp, err := nc.cimiRequest(r)
	if err != nil {
		log.Errorf("Error executing POST request: %s", err)
		return nil, err
	}

	return resp, nil
}

func (nc *NuvlaClient) Put(uri string, data map[string]interface{}, toDelete []string) (*http.Response, error) {
	r := &types.RequestInput{
		Method:   "PUT",
		Endpoint: nc.buildUriEndPoint(uri),
		JsonData: data,
		Params: &types.RequestParams{
			Select: toDelete,
		},
	}

	resp, err := nc.cimiRequest(r)
	if err != nil {
		log.Errorf("Error executing PUT request: %s", err)
		return nil, err
	}
	return resp, nil
}

func (nc *NuvlaClient) delete(deleteEndpoint string) (*http.Response, error) {
	r := &types.RequestInput{
		Method:   "DELETE",
		Endpoint: deleteEndpoint,
	}
	resp, err := nc.cimiRequest(r)
	if err != nil {
		log.Errorf("Error executing PUT request: %s", err)
		return nil, err
	}
	return resp, nil
}

func (nc *NuvlaClient) Operation(resourceId, operation string, data map[string]interface{}) (*http.Response, error) {
	return nc.Post(nc.buildOperationUriEndPoint(resourceId, operation), data)
}

func (nc *NuvlaClient) Edit(resourceId string, data map[string]interface{}, toSelect []string) (*http.Response, error) {
	return nc.Put(nc.buildOperationUriEndPoint(resourceId, "edit"), data, toSelect)
}

func (nc *NuvlaClient) Delete(resourceId string) (*http.Response, error) {
	return nc.delete(nc.buildOperationUriEndPoint(resourceId, "delete"))
}
