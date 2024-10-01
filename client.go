package api_client_go

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nuvla/api-client-go/clients/resources"
	"github.com/nuvla/api-client-go/common"
	"github.com/nuvla/api-client-go/types"
	log "github.com/sirupsen/logrus"
	"github.com/wI2L/jsondiff"
	"io"
	"net/http"
)

type NuvlaClient struct {
	// Session params
	*NuvlaSession
	SessionOpts SessionOptions
	Credentials types.LogInParams
}

func NewNuvlaClient(cred types.LogInParams, opts *SessionOptions) *NuvlaClient {
	nc := &NuvlaClient{
		NuvlaSession: NewNuvlaSession(opts),
		SessionOpts:  *opts,
	}

	if !common.IsNilValueInterface(cred) {
		log.Debug("Logging in with api keys...")
		if err := nc.login(cred); err != nil {
			log.Errorf("Error logging in with api keys: %s.", err)
		} else {
			nc.Credentials = cred
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
	logInParams := types.NewApiKeyLogInParams(key, secret)
	err := nc.login(logInParams)
	if err != nil {
		log.Errorf("Error logging in with api keys: %s", err)
		return err
	}
	// Save login params if successful and Credentials are different from the current ones
	nc.Credentials = logInParams
	return nil
}

func (nc *NuvlaClient) LoginUser(username string, password string) error {
	logInParams := types.NewUserLogInParams(username, password)
	err := nc.login(logInParams)
	if err != nil {
		log.Errorf("Error logging in with user Credentials: %s", err)
		return err
	}
	// Save login params if successful and Credentials are different from the current ones
	nc.Credentials = logInParams
	return nil
}

func (nc *NuvlaClient) Logout() error {
	// Close connections and logout

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

func (nc *NuvlaClient) needsAuthentication(statusCode int, url string) bool {
	matchStatusCode := statusCode == http.StatusUnauthorized || statusCode == http.StatusForbidden
	matchEndpoint := url == fmt.Sprintf("%s/api/session", nc.endpoint)
	return nc.SessionOpts.ReAuthenticate && matchStatusCode && !matchEndpoint
}

func (nc *NuvlaClient) cimiRequest(ctx context.Context, reqInput *types.RequestOpts) (*http.Response, error) {
	// Setup default client headers for all requests
	// TODO: Might be configurable from session
	if reqInput.Headers == nil {
		reqInput.Headers = make(map[string]string)
	}
	reqInput.Headers["Accept"] = "application/json"
	if nc.compress {
		reqInput.Headers["Accept-Encoding"] = "gzip"
	}
	if reqInput.Bulk {
		reqInput.Headers["bulk"] = "true"
	}

	r, err := nc.Request(ctx, reqInput)
	if err != nil {
		if r != nil {
			_ = r.Body.Close()
		}
		return nil, err
	}

	if nc.needsAuthentication(r.StatusCode, reqInput.Endpoint) {
		// Read response body
		b, err := io.ReadAll(r.Body)
		if err != nil {
			_ = r.Body.Close()
			return nil, fmt.Errorf("error reading response body: %s", err)
		}
		log.Infof("Response body: %s", string(b))
		log.Infof("Request: %s-%s", reqInput.Method, reqInput.Endpoint)

		// Request: Unauthorized
		log.Infof("Re-authenticating...")
		if err := nc.login(nc.Credentials); err != nil {
			return nil, fmt.Errorf("error re-authenticating: %s", err)
		}

		// Retry request
		r, err = nc.Request(ctx, reqInput)
		if err != nil {
			if r != nil {
				_ = r.Body.Close()
			}
			return nil, fmt.Errorf("error re-executing request: %s", err)
		}
	}

	return r, nil
}

// Get executes the get http method
// Allow for selective fields to be returned via the selectFields parameter

func (nc *NuvlaClient) Get(ctx context.Context, resourceId string, selectFields []string) (*types.NuvlaResource, error) {
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

	resp, err := nc.cimiRequest(ctx, r)
	if err != nil {
		log.Errorf("Error executing GET request: %s", err)
		return nil, err
	}

	return types.NewResourceFromResponse(resp), nil
}

// Post executes the post http method
// Data can be any type, but it will be marshaled into JSON
func (nc *NuvlaClient) Post(ctx context.Context, endpoint string, data map[string]interface{}) (*http.Response, error) {
	r := &types.RequestOpts{
		Method:   "POST",
		JsonData: data,
		Endpoint: nc.buildUriEndPoint(endpoint),
	}

	resp, err := nc.cimiRequest(ctx, r)
	if err != nil {
		log.Errorf("Error executing POST request: %s", err)
		return nil, err
	}

	return resp, nil
}

func (nc *NuvlaClient) BulkPost(ctx context.Context, endpoint string, data []map[string]interface{}) (*http.Response, error) {
	r := &types.RequestOpts{
		Method:   "POST",
		JsonData: data,
		Endpoint: nc.buildUriEndPoint(endpoint),
		Bulk:     true,
	}

	resp, err := nc.cimiRequest(ctx, r)
	if err != nil {
		log.Errorf("Error executing POST request: %s", err)
		return nil, err
	}

	return resp, nil
}

func (nc *NuvlaClient) Put(ctx context.Context, uri string, data interface{}, selectFields []string) (*http.Response, error) {
	r := &types.RequestOpts{
		Method:   "PUT",
		Endpoint: nc.buildUriEndPoint(uri),
		JsonData: data,
		Params: &types.RequestParams{
			Select: selectFields,
		},
		Headers: make(map[string]string),
	}
	_, isPatch := data.(jsondiff.Patch)
	if isPatch {
		r.Headers["Content-Type"] = "application/json-patch+json"
	}

	resp, err := nc.cimiRequest(ctx, r)
	if err != nil {
		log.Errorf("Error executing PUT request: %s", err)
		return nil, err
	}
	return resp, nil
}

func (nc *NuvlaClient) delete(ctx context.Context, deleteEndpoint string) (*http.Response, error) {
	r := &types.RequestOpts{
		Method:   "DELETE",
		Endpoint: deleteEndpoint,
	}
	resp, err := nc.cimiRequest(ctx, r)
	if err != nil {
		log.Errorf("Error executing PUT request: %s", err)
		return nil, err
	}
	return resp, nil
}

func (nc *NuvlaClient) Operation(ctx context.Context, resourceId, operation string, data map[string]interface{}) (*http.Response, error) {
	return nc.Post(ctx, nc.buildOperationUriEndPoint(resourceId, operation), data)
}

func (nc *NuvlaClient) BulkOperation(ctx context.Context, resourceId string, operation string, data []map[string]interface{}) (*http.Response, error) {
	return nc.BulkPost(ctx, nc.buildOperationUriEndPoint(resourceId, operation), data)
}

func (nc *NuvlaClient) Edit(ctx context.Context, resourceId string, data map[string]interface{}, toSelect []string) (*http.Response, error) {
	return nc.Put(ctx, resourceId, data, toSelect)
}

func (nc *NuvlaClient) Delete(ctx context.Context, resourceId string) (*http.Response, error) {
	return nc.delete(ctx, nc.buildOperationUriEndPoint(resourceId, "delete"))
}

type SearchOptions struct {
	First       int      `json:"first"`
	Last        int      `json:"last"`
	Filter      string   `json:"filter"`
	Fields      string   `json:"fields"`
	OrderBy     string   `json:"orderby"`
	Select      []string `json:"select"`
	Aggregation string   `json:"aggregation"`
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

func (nc *NuvlaClient) Search(ctx context.Context, resourceType string, opts *SearchOptions) (*resources.NuvlaResourceCollection, error) {

	r := &types.RequestOpts{
		Method:   "PUT",
		Endpoint: nc.buildUriEndPoint(resourceType),
		Params:   nil,
		JsonData: nil,
		Data:     common.GetCleanMapFromStruct(opts),
	}
	resp, err := nc.cimiRequest(ctx, r)

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

// Add creates a new resource of the given type and returns its ID
func (nc *NuvlaClient) Add(ctx context.Context, resourceType resources.NuvlaResourceType, data map[string]interface{}) (*types.NuvlaID, error) {
	res, err := nc.Post(ctx, string(resourceType), data)
	if err != nil {
		log.Errorf("Error adding %s: %s", resourceType, err)
		return nil, err
	}

	var resData map[string]interface{}

	bodyBytes, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		log.Errorf("Error reading response body, cannot extract ID: %s", err)
		return nil, err
	}

	err = json.Unmarshal(bodyBytes, &resData)
	if err != nil {
		log.Errorf("Error unmarshaling response body, cannot extract ID: %s", err)
		return nil, err
	}

	log.Debugf("ID of new %s: %s", resourceType, resData)

	return types.NewNuvlaIDFromId(resData["resource-id"].(string)), nil
}
