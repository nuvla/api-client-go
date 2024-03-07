package api_client_go

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type NuvlaEdgeClient struct {
	apiClient *NuvlaClient

	NuvlaEdgeId       *NuvlaID
	NuvlaEdgeStatusId *NuvlaID
	CredentialId      *NuvlaID

	resource    map[string]interface{}
	credentials *ApiKeyLogInParams
}

type NuvlaEdgeClientOpts struct {
	NuvlaClientOpts *SessionOptions

	NuvlaEdgeId *NuvlaID
	Credentials *ApiKeyLogInParams
}

func NewNuvlaEdgeClient(nuvlaEdgeClientOpts *NuvlaEdgeClientOpts) *NuvlaEdgeClient {
	s := NewNuvlaClientFromOpts(nuvlaEdgeClientOpts.NuvlaClientOpts)
	ne := &NuvlaEdgeClient{
		apiClient:   s,
		NuvlaEdgeId: nuvlaEdgeClientOpts.NuvlaEdgeId,
		credentials: nuvlaEdgeClientOpts.Credentials,
	}
	if ne.credentials != nil {
		log.Debug("Logging in with api keys...")
		if err := s.LoginApiKeys(ne.credentials.Key, ne.credentials.Secret); err != nil {
			log.Errorf("Error logging in with api keys: %s", err)
		}
		log.Debug("Logging in with api keys... Success.")
	}
	return ne
}

//func NewNuvlaEdgeClient(nuvlaEdgeId *api_client_go.NuvlaID, endpoint string, insecure bool, debug bool) *NuvlaEdgeClient {
//	return &NuvlaEdgeClient{
//		NuvlaEdgeId: nuvlaEdgeId,
//		apiClient: api_client_go.NewNuvlaClient(&api_client_go.SessionOptions{
//			Endpoint: endpoint,
//			Insecure: insecure,
//			Debug:    debug,
//		}),
//	}
//}

func (ne *NuvlaEdgeClient) GetResource(fields []string) error {
	res, err := ne.apiClient.Get(ne.NuvlaEdgeId.String(), fields)
	if err != nil {
		log.Errorf("Error getting NuvlaEdge resource: %s", err)
		return err
	}
	ne.resource = res.Data
	ne.NuvlaEdgeStatusId = NewNuvlaIDFromId(ne.resource["status"].(string))
	ne.CredentialId = NewNuvlaIDFromId(ne.resource["credential"].(string))

	return nil
}

func extractCredentialsFromActivateResponse(resp *http.Response) (*ApiKeyLogInParams, error) {
	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("Error reading response body: %s", err)
		return nil, err
	}
	defer resp.Body.Close()
	creds := &ApiKeyLogInParams{}
	err = json.Unmarshal(body, creds)
	if err != nil {
		log.Errorf("Error unmarshaling response body: %s", err)
		return nil, err
	}

	return creds, nil
}

// Activate Operation
func (ne *NuvlaEdgeClient) Activate() error {
	res, err := ne.apiClient.Operation(ne.NuvlaEdgeId.String(), "activate", nil)
	if err != nil {
		log.Errorf("Error activating NuvlaEdge: %s", err)
		return err
	}
	creds, err := extractCredentialsFromActivateResponse(res)
	log.Infof("Credentials received from activation: %v", creds)
	ne.credentials = creds

	return nil
}

// Commission operations
func (ne *NuvlaEdgeClient) Commission(data map[string]interface{}) error {
	log.Debugf("Commissioning NuvlaEdge with payload %v", data)
	res, err := ne.apiClient.Operation(ne.NuvlaEdgeId.String(), "commission", data)
	if err != nil {
		log.Errorf("Error commissioning NuvlaEdge: %s", err)
		return err
	}
	log.Debug("Commissioning NuvlaEdge... Success.")
	log.Debug("Commissioning response: %v", res.Body)
	return nil
}

// Telemetry operation
func (ne *NuvlaEdgeClient) Telemetry(data map[string]interface{}, Select []string) error {
	log.Debugf("Sending telemetry data to NuvlaEdge with payload %v", data)
	res, err := ne.apiClient.Put(ne.NuvlaEdgeStatusId.String(), data, Select)
	if err != nil {
		log.Errorf("Error sending telemetry data to NuvlaEdge: %s", err)
		return err
	}
	log.Debug("Sending telemetry data to NuvlaEdge... Success.")
	log.Debug("Telemetry response: %v", res.Body)
	return nil
}
