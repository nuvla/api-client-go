package api_client_go

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type NuvlaClient struct {
	nuvlaEndpoint string
	debug         bool

	// Session params
	session *NuvlaSession
}

func NewNuvlaClientFromOpts(attrs *SessionOptions) *NuvlaClient {

	nc := &NuvlaClient{
		nuvlaEndpoint: attrs.Endpoint,
		debug:         attrs.Debug,
		session:       NewNuvlaSession(attrs),
	}
	return nc
}

func NewNuvlaClient(endpoint string, insecure bool, debug bool) *NuvlaClient {
	opts := &SessionOptions{
		Endpoint: endpoint,
		Insecure: insecure,
		Debug:    debug,
	}
	return NewNuvlaClientFromOpts(NewSessionOpts(opts))
}

func (nc *NuvlaClient) LoginApiKeys(key string, secret string) error {
	err := nc.session.login(NewApiKeyLogInParams(key, secret))
	if err != nil {
		log.Errorf("Error logging in with api keys: %s", err)
		return err
	}
	return nil
}

func (nc *NuvlaClient) LoginUser(username string, password string) error {
	err := nc.session.login(NewUserLogInParams(username, password))
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

func (nc *NuvlaClient) cimiRequest(reqInput *RequestOpts) (*http.Response, error) {
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

func (nc *NuvlaClient) Get(resourceId string, selectFields []string) (*NuvlaResource, error) {
	// Define request inputs to allow adding select fields
	r := &RequestOpts{
		Method:   "GET",
		Endpoint: nc.buildUriEndPoint(resourceId),
	}

	// Do not create the request params struct unless we need it to prevent overhead down the line
	if selectFields != nil {
		r.Params = &RequestParams{
			Select: selectFields,
		}
	}

	resp, err := nc.cimiRequest(r)
	if err != nil {
		log.Errorf("Error executing GET request: %s", err)
		return nil, err
	}

	return NewResourceFromResponse(resp), nil
}

// Post executes the post http method
// Data can be any type, but it will be marshaled into JSON
func (nc *NuvlaClient) Post(endpoint string, data map[string]interface{}) (*http.Response, error) {
	r := &RequestOpts{
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
	r := &RequestOpts{
		Method:   "PUT",
		Endpoint: nc.buildUriEndPoint(uri),
		JsonData: data,
		Params: &RequestParams{
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
	r := &RequestOpts{
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
