package api_client_go

import (
	"fmt"
	"github.com/nuvla/api-client-go/clients/resources"
	"github.com/nuvla/api-client-go/types"
	log "github.com/sirupsen/logrus"
	"net/http"
	"reflect"
	"strconv"
)

type ClientOpts struct {
	*SessionOptions

	Credentials types.LogInParams
}

func NewClientOpts(credentials types.LogInParams, opts ...SessionOptFunc) *ClientOpts {
	sessionOpts := DefaultSessionOpts()
	for _, fn := range opts {
		fn(sessionOpts)
	}
	return &ClientOpts{
		SessionOptions: sessionOpts,
		Credentials:    credentials,
	}
}

type NuvlaClient struct {
	// Session params
	*NuvlaSession
	credentials types.LogInParams
}

func NewNuvlaClient(cred types.LogInParams, opts *SessionOptions) *NuvlaClient {
	nc := &NuvlaClient{
		NuvlaSession: NewNuvlaSession(opts),
		credentials:  cred,
	}

	if nc.credentials != nil && nc.reauthenticate {
		log.Debug("Logging in with api keys...")
		if err := nc.login(nc.credentials); err != nil {
			log.Errorf("Error logging in with api keys: %s.", err)
		}

	}
	return nc
}

func NewNuvlaClientFromOpts(cred types.LogInParams, opts ...SessionOptFunc) *NuvlaClient {
	sessionOpts := DefaultSessionOpts()
	for _, fn := range opts {
		fn(sessionOpts)
	}

	return NewNuvlaClient(cred, sessionOpts)
}

func (nc *NuvlaClient) LoginApiKeys(key string, secret string) error {
	err := nc.login(types.NewApiKeyLogInParams(key, secret))
	if err != nil {
		log.Errorf("Error logging in with api keys: %s", err)
		return err
	}
	return nil
}

func (nc *NuvlaClient) LoginUser(username string, password string) error {
	err := nc.login(types.NewUserLogInParams(username, password))
	if err != nil {
		log.Errorf("Error logging in with user credentials: %s", err)
		return err
	}
	return nil
}

func (nc *NuvlaClient) Logout() error {
	return nc.logout()
}

func (nc *NuvlaClient) IsAuthenticated() bool {
	return true
}

func (nc *NuvlaClient) buildUriEndPoint(uriEndpoint string) string {
	return fmt.Sprintf("%s/api/%s", nc.endpoint, uriEndpoint)
}

func (nc *NuvlaClient) buildOperationUriEndPoint(uriEndpoint string, operation string) string {
	return fmt.Sprintf("%s/%s", uriEndpoint, operation)
}

func (nc *NuvlaClient) cimiRequest(reqInput *types.RequestOpts) (*http.Response, error) {
	// Setup default client headers for all requests
	// TODO: Might be configurable from session
	if reqInput.Headers == nil {
		reqInput.Headers = make(map[string]string)
	}
	reqInput.Headers["Accept"] = "application/json"
	if nc.compress {
		reqInput.Headers["Accept-Encoding"] = "gzip"
	}

	r, _ := nc.Request(reqInput)

	return r, nil
}

// Get executes the get http method
// Allow for selective fields to be returned via the selectFields parameter

func (nc *NuvlaClient) Get(resourceId string, selectFields []string) (*types.NuvlaResource, error) {
	// Define request inputs to allow adding select fields
	r := &types.RequestOpts{
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
	r := &types.RequestOpts{
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

func (nc *NuvlaClient) Put(uri string, data map[string]interface{}, selectFields []string) (*http.Response, error) {
	r := &types.RequestOpts{
		Method:   "PUT",
		Endpoint: nc.buildUriEndPoint(uri),
		JsonData: data,
		Params: &types.RequestParams{
			Select: selectFields,
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
	r := &types.RequestOpts{
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
	return nc.Put(resourceId, data, toSelect)
}

func (nc *NuvlaClient) Delete(resourceId string) (*http.Response, error) {
	return nc.delete(nc.buildOperationUriEndPoint(resourceId, "delete"))
}

type SearchOptions struct {
	First       int    `json:"first"`
	Last        int    `json:"last"`
	Filter      string `json:"filter"`
	Fields      string `json:"fields"`
	OrderBy     string `json:"orderby"`
	Aggregation string `json:"aggregation"`
}

// GetCleanMap returns a map with only the non-nil fields
func (so *SearchOptions) GetCleanMap() map[string]string {
	m := make(map[string]string)
	val := reflect.ValueOf(so).Elem()

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		jsonTag := typeField.Tag.Get("json")

		if !valueField.IsZero() {
			if typeField.Name == "First" || typeField.Name == "Last" {
				m[jsonTag] = strconv.Itoa(int(valueField.Int()))
			} else {
				m[jsonTag] = valueField.String()
			}
		}
	}
	return m
}

func NewDefaultSearchOptions() *SearchOptions {
	return &SearchOptions{
		First:       0,
		Last:        0,
		Filter:      "",
		Fields:      "",
		OrderBy:     "",
		Aggregation: "",
	}
}

func (nc *NuvlaClient) Search(resourceType string, opts *SearchOptions) (*resources.NuvlaResourceCollection, error) {

	r := &types.RequestOpts{
		Method:   "PUT",
		Endpoint: nc.buildUriEndPoint(resourceType),
		Params:   nil,
		JsonData: nil,
		Data:     opts.GetCleanMap(),
	}
	resp, err := nc.cimiRequest(r)
	if err != nil {
		log.Errorf("Error executing GET request: %s", err)
		return nil, err
	}
	collection, err := resources.NewCollectionFromResponse(resp)
	if err != nil {
		log.Errorf("Error creating resource collection: %s", err)
		return nil, err
	}
	return collection, err
}
