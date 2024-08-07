package clients

import (
	"encoding/json"
	"fmt"
	nuvla "github.com/nuvla/api-client-go"
	"github.com/nuvla/api-client-go/clients/resources"
	"github.com/nuvla/api-client-go/common"
	"github.com/nuvla/api-client-go/types"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type NuvlaEdgeSessionFreeze struct {
	// Session data
	nuvla.SessionOptions

	// Client data
	Credentials *types.ApiKeyLogInParams `json:"credentials"`

	// NuvlaEdge Client
	NuvlaEdgeId       string `json:"nuvlaedge-uuid"`
	NuvlaEdgeStatusId string `json:"nuvlaedge-status-uuid"`
	InfraServiceId    string `json:"infra-service-id"`
	VPNServiceId      string `json:"vpn-service-id"`
}

func (sf *NuvlaEdgeSessionFreeze) Load(file string) error {
	return common.ReadJSONFromFile(file, sf)
}

func (sf *NuvlaEdgeSessionFreeze) Save(file string) error {
	// Write b to file
	err := common.WriteIndentedJSONToFile(sf, file)
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

	nuvlaEdgeResource *resources.NuvlaEdgeResource
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
	defer resp.Body.Close()
	if err != nil {
		log.Errorf("Error reading response body: %s", err)
		return nil, err
	}
	log.Infof("Response body from activation: %s", string(body))
	c := make(map[string]string)
	err = json.Unmarshal(body, &c)
	creds := &types.ApiKeyLogInParams{}
	k, ok := c["api-key"]
	if !ok {
		return nil, fmt.Errorf("api-key not found in response")
	}
	creds.Key = k

	s, ok := c["secret-key"]
	if !ok {
		return nil, fmt.Errorf("secret-key not found in response")
	}
	creds.Secret = s

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
	if res.StatusCode != 200 && res.StatusCode != 201 {
		return fmt.Errorf("activation failed with status code: %v", res.StatusCode)
	}

	creds, err := extractCredentialsFromActivateResponse(res)

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

	err = res.Body.Close()
	if err != nil {
		log.Errorf("Error closing commission response body: %s", err)
	}

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
	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("heartbeat failed with status code: %v", res.StatusCode)
	}

	log.Debug("Sending heartbeat to NuvlaEdge... Success.")
	return res, nil
}

func (ne *NuvlaEdgeClient) GetId() string {
	return ne.NuvlaEdgeId.Id
}

func (ne *NuvlaEdgeClient) GetType() resources.NuvlaResourceType {
	return resources.NuvlaEdgeType
}

func (ne *NuvlaEdgeClient) GetResourceMap() (map[string]interface{}, error) {
	var mapRes map[string]interface{}
	if err := MarshalResourceIntoMap(ne.nuvlaEdgeResource, mapRes); err != nil {
		log.Errorf("Error marshaling NuvlaEdgeResource to map")
		return nil, err
	}
	return mapRes, nil
}

func (ne *NuvlaEdgeClient) UpdateResourceSelect(selects []string) error {
	res, err := ne.Get(ne.NuvlaEdgeId.Id, selects)
	if err != nil {
		log.Infof("Error updating NuvlaEdge resource %s", ne.NuvlaEdgeId)
		return nil
	}

	if ne.nuvlaEdgeResource == nil {
		ne.nuvlaEdgeResource = &resources.NuvlaEdgeResource{}
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

func (ne *NuvlaEdgeClient) UpdateResource() error {
	return ne.UpdateResourceSelect(nil)
}

func (ne *NuvlaEdgeClient) GetNuvlaEdgeResource() resources.NuvlaEdgeResource {
	if ne.nuvlaEdgeResource == nil {
		ne.nuvlaEdgeResource = &resources.NuvlaEdgeResource{}
	}
	return *ne.nuvlaEdgeResource
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
