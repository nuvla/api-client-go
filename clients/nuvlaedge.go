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

type NuvlaEdgeSessionFreeze struct {
	// Session data
	nuvla.SessionOptions

	// Client data
	Credentials *types.ApiKeyLogInParams `json:"credentials"`

	// NuvlaEdge Client
	NuvlaEdgeId       string `json:"nuvlaedge-id"`
	NuvlaEdgeStatusId string `json:"nuvlaedge-status-id"`
	InfraServiceId    string `json:"infra-service-id"`
	VPNServiceId      string `json:"vpn-service-id"`
}

func (sf *NuvlaEdgeSessionFreeze) Load(file string) error {
	return nuvla.ReadJSONFromFile(file, sf)
}

func (sf *NuvlaEdgeSessionFreeze) Save(file string) error {
	// Write b to file
	err := nuvla.WriteIndentedJSONToFile(sf, file)
	if err != nil {
		log.Errorf("Error saving NuvlaEdgeSessionFreeze: %s", err)
		return err
	}

	return nil
}

type NuvlaEdgeClient struct {
	*nuvla.NuvlaClient

	NuvlaEdgeId       *types.NuvlaID
	NuvlaEdgeStatusId *types.NuvlaID
	CredentialId      *types.NuvlaID

	nuvlaEdgeResource *NuvlaEdgeResource
	Credentials       *types.ApiKeyLogInParams
}

func NewNuvlaEdgeClient(nuvlaEdgeId string, credentials *types.ApiKeyLogInParams, opts ...nuvla.SessionOptFunc) *NuvlaEdgeClient {
	sessionOpts := nuvla.DefaultSessionOpts()
	for _, fn := range opts {
		fn(sessionOpts)
	}
	log.Infof("Creating NuvlaEdge client with options: %v", sessionOpts)
	ne := &NuvlaEdgeClient{
		NuvlaClient: nuvla.NewNuvlaClient(credentials, sessionOpts),
		NuvlaEdgeId: types.NewNuvlaIDFromId(nuvlaEdgeId),
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

func NewNuvlaEdgeClientFromSessionFreeze(f *NuvlaEdgeSessionFreeze) *NuvlaEdgeClient {
	log.Infof("Creating NuvlaEdge client from session freeze")
	ne := &NuvlaEdgeClient{}

	ne.NuvlaEdgeId = types.NewNuvlaIDFromId(f.NuvlaEdgeId)
	ne.NuvlaEdgeStatusId = types.NewNuvlaIDFromId(f.NuvlaEdgeStatusId)
	ne.Credentials = f.Credentials

	// Create NuvlaClient
	ne.NuvlaClient = nuvla.NewNuvlaClient(ne.Credentials, &f.SessionOptions)

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
		ne.NuvlaEdgeStatusId = types.NewNuvlaIDFromId(ne.nuvlaEdgeResource.NuvlaBoxStatus)
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

func (ne *NuvlaEdgeClient) Freeze(file string) error {
	log.Infof("Freezing NuvlaEdge client...")
	f := NuvlaEdgeSessionFreeze{
		SessionOptions: ne.GetSessionOpts(),
		Credentials:    ne.Credentials,
		// If this point is reached, NuvlaEdgeID should never be nil or empty so if null pointer exception
		// is raised here, there is another issue
		NuvlaEdgeId:    ne.NuvlaEdgeId.String(),
		InfraServiceId: ne.nuvlaEdgeResource.InfrastructureServiceGroup,
	}

	if ne.NuvlaEdgeStatusId != nil {
		f.NuvlaEdgeStatusId = ne.NuvlaEdgeStatusId.String()
	}

	return f.Save(file)
}
