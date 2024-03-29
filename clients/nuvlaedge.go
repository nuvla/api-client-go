package clients

import (
	"encoding/json"
	"fmt"
	nuvla "github.com/nuvla/api-client-go"
	"github.com/nuvla/api-client-go/types"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type NuvlaEdgeState string

const (
	NuvlaEdgeStateNew             NuvlaEdgeState = "NEW"
	NuvlaEdgeStateActivated       NuvlaEdgeState = "ACTIVATED"
	NuvlaEdgeStateCommissioned    NuvlaEdgeState = "COMMISSIONED"
	NuvlaEdgeStateDecommissioning NuvlaEdgeState = "DECOMMISSIONING"
	NuvlaEdgeStateDecommissioned  NuvlaEdgeState = "DECOMMISSIONED"
	NuvlaEdgeStateError           NuvlaEdgeState = "ERROR"
	NuvlaEdgeStateSuspended       NuvlaEdgeState = "SUSPENDED"
)

type NuvlaEdgeResource struct {
	Resource

	// Required
	State             NuvlaEdgeState `json:"state"`
	RefreshInterval   int            `json:"refresh-interval"`
	HeartbeatInterval int            `json:"heartbeat-interval"`
	Version           int            `json:"version"`
	Owner             string         `json:"owner"`

	// Optional
	Location              []string `json:"location"`
	VPNServerID           string   `json:"vpn-server-id"`
	SSHKeys               []string `json:"ssh-keys"`
	Capabilities          []string `json:"capabilities"`
	Online                bool     `json:"online"`
	InferredLocation      string   `json:"inferred-location"`
	NuvlaBoxEngineVersion string   `json:"nuvlabox-engine-version"`

	NuvlaBoxStatus             string `json:"nuvlabox-status"`
	InfrastructureServiceGroup string `json:"infrastructure-service-group"`
	CredentialApiKey           string `json:"credential-api-key"`
	HostLevelManagementApiKey  string `json:"host-level-management-api-key"`
}

type NuvlaEdgeClient struct {
	*nuvla.NuvlaClient

	NuvlaEdgeId       *types.NuvlaID
	NuvlaEdgeStatusId types.NuvlaID
	CredentialId      *types.NuvlaID

	nuvlaEdgeResource *NuvlaEdgeResource
	Credentials       *types.ApiKeyLogInParams
}

func NewNuvlaEdgeClient(nuvlaEdgeId *types.NuvlaID, credentials *types.ApiKeyLogInParams, opts ...nuvla.SessionOptFunc) *NuvlaEdgeClient {
	sessionOpts := nuvla.DefaultSessionOpts()
	for _, fn := range opts {
		fn(sessionOpts)
	}
	log.Infof("Creating NuvlaEdge client with options: %v", sessionOpts)
	ne := &NuvlaEdgeClient{
		NuvlaClient: nuvla.NewNuvlaClient(credentials, sessionOpts),
		NuvlaEdgeId: nuvlaEdgeId,
		Credentials: credentials,
	}

	if ne.Credentials != nil {
		log.Debug("Logging in with api keys...")
		if err := ne.LoginApiKeys(ne.Credentials.Key, ne.Credentials.Secret); err != nil {
			log.Errorf("Error logging in with api keys: %s. Retry manually...", err)

		} else {
			log.Debug("Logging in with api keys... Success.")
		}
	}
	return ne
}

func extractCredentialsFromActivateResponse(resp *http.Response) (*types.ApiKeyLogInParams, error) {

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("Error reading response body: %s", err)
		return nil, err
	}
	defer resp.Body.Close()

	creds := &types.ApiKeyLogInParams{}
	err = json.Unmarshal(body, creds)
	if err != nil {
		log.Errorf("Error unmarshaling response body: %s", err)
		return nil, err
	}

	return creds, nil
}

// LogIn operation
func (ne *NuvlaEdgeClient) LogIn() error {
	err := ne.LoginApiKeys(ne.Credentials.Key, ne.Credentials.Secret)
	if err != nil {
		log.Errorf("Error logging in with api keys: %s", err)
		return err
	}
	log.Infof("Logging in with api keys... Success.")
	return nil
}

// Activate Operation
func (ne *NuvlaEdgeClient) Activate() error {
	log.Infof("Activating NuvlaEdge...%v", ne.NuvlaEdgeId)
	res, err := ne.Operation(ne.NuvlaEdgeId.String(), "activate", nil)
	if err != nil {
		log.Errorf("Error activating NuvlaEdge: %s", err)
		return err
	}
	log.Infof("Code response from activation: %v", res.StatusCode)
	var b []byte
	_, err = res.Body.Read(b)
	log.Infof("Response from activation: %s", string(b))

	creds, err := extractCredentialsFromActivateResponse(res)
	log.Infof("Credentials received from activation: %v", creds)
	ne.Credentials = creds

	return nil
}

// Commission operations
func (ne *NuvlaEdgeClient) Commission(data map[string]interface{}) error {
	log.Debugf("Commissioning NuvlaEdge with payload %v", data)
	res, err := ne.Operation(ne.NuvlaEdgeId.String(), "commission", data)
	if err != nil {
		log.Errorf("Error commissioning NuvlaEdge: %s", err)
		return err
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Errorf("Error reading response body: %s", err)
		return err
	}
	defer res.Body.Close()

	log.Infof("Commissioning code response: %v", res.StatusCode)
	log.Infof("Commissioning response: %s", string(body))
	if res.StatusCode != 200 && res.StatusCode != 201 {
		return fmt.Errorf("commissioning failed with status code: %v", res.StatusCode)
	}

	log.Infof("Updating NuvlaEdge resource...")
	err = ne.UpdateResource()
	if err != nil {
		log.Errorf("Error getting NuvlaEdge resource: %s", err)
		return err
	}
	log.Infof("Updating NuvlaEdge resource... Success")

	if ne.nuvlaEdgeResource.NuvlaBoxStatus != "" {
		ne.NuvlaEdgeStatusId = *types.NewNuvlaIDFromId(ne.nuvlaEdgeResource.NuvlaBoxStatus)
	} else {
		return fmt.Errorf("nuvlabox-status not found in resource")
	}
	return nil
}

// Telemetry operation
func (ne *NuvlaEdgeClient) Telemetry(data map[string]interface{}, Select []string) (*http.Response, error) {
	log.Debugf("Sending telemetry data to NuvlaEdge with payload %v", data)
	res, err := ne.Put(ne.NuvlaEdgeStatusId.String(), data, Select)
	if err != nil {
		log.Errorf("Error sending telemetry data to Nuvla: %s", err)
		return nil, err
	}
	return res, nil
}

// Heartbeat operation
func (ne *NuvlaEdgeClient) Heartbeat() (*http.Response, error) {
	log.Debug("Sending heartbeat to NuvlaEdge...")

	res, err := ne.Operation(ne.NuvlaEdgeId.String(), "heartbeat", nil)
	if err != nil {
		log.Errorf("Error sending heartbeat to NuvlaEdge: %s", err)
		return nil, err
	}
	log.Infof("Code response from heartbeat: %v", res.StatusCode)
	log.Debug("Sending heartbeat to NuvlaEdge... Success.")
	return res, nil
}

func (ne *NuvlaEdgeClient) GetId() string {
	return ne.NuvlaEdgeId.Id
}

func (ne *NuvlaEdgeClient) GetType() ClientResourceType {
	return NuvlaEdgeType
}

func (ne *NuvlaEdgeClient) GetResourceMap() (map[string]interface{}, error) {
	var mapRes map[string]interface{}
	if err := MarshalResourceIntoMap(ne.nuvlaEdgeResource, mapRes); err != nil {
		log.Errorf("Error marshaling NuvlaEdgeResource to map")
		return nil, err
	}
	return mapRes, nil
}

func (ne *NuvlaEdgeClient) UpdateResource() error {
	res, err := ne.Get(ne.NuvlaEdgeId.Id, nil)
	if err != nil {
		log.Infof("Error updating NuvlaEdge resource %s", ne.NuvlaEdgeId)
		return nil
	}

	if ne.nuvlaEdgeResource == nil {
		ne.nuvlaEdgeResource = &NuvlaEdgeResource{}
	}

	b, err := json.Marshal(res.Data)
	if err != nil {
		log.Errorf("Error marshaling NuvlaEdge resource response data to bytes")
		return err
	}

	err = json.Unmarshal(b, ne.nuvlaEdgeResource)
	if err != nil {
		log.Error("Error unmarshalling response into NuvlaEdgeResource structure")
		return err
	}
	log.Infof("Successfully updated NuvlaEdge resource")
	return nil
}

func (ne *NuvlaEdgeClient) GetNuvlaClient() *nuvla.NuvlaClient {
	return ne.NuvlaClient
}
