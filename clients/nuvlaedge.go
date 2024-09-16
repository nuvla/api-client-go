package clients

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	nuvla "github.com/nuvla/api-client-go"
	"github.com/nuvla/api-client-go/clients/resources"
	"github.com/nuvla/api-client-go/common"
	"github.com/nuvla/api-client-go/types"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strings"
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
	if err := common.ReadJSONFromFile(file, sf); err != nil {
		return fmt.Errorf("error loading NuvlaEdgeSessionFreeze: %s", err)
	}

	if sf.Credentials != nil && sf.Credentials.Key != "" && sf.Credentials.Secret != "" {
		decodeCredentialsIfNeeded(sf.Credentials)
	}

	return nil
}

func (sf *NuvlaEdgeSessionFreeze) Save(file string) error {
	// Write b to file
	credsCopy := *sf.Credentials
	freezeCopy := *sf
	freezeCopy.Credentials = &credsCopy
	encodeCredentialsIfNeeded(&credsCopy)
	err := common.WriteIndentedJSONToFile(freezeCopy, file)
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
	creds, ok := ne.Credentials.(*types.ApiKeyLogInParams)
	if !ok {
		return errors.New("credentials not properly formated, exiting")
	}

	err := ne.LoginApiKeys(creds.Key, creds.Secret)
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

	log.Debug("Updating NuvlaEdge resource...")
	err = ne.UpdateResource()
	if err != nil {
		log.Errorf("Error getting NuvlaEdge resource: %s", err)
		return err
	}
	log.Debug("Updating NuvlaEdge resource... Success")

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
	if ne.nuvlaEdgeResource.NuvlaBoxStatus == "" || ne.NuvlaEdgeStatusId == nil {
		err := ne.UpdateResourceSelect([]string{"nuvlabox-status"})
		if err != nil {
			log.Errorf("Error sending Telemetry, cannot find NuvlaBoxStatus ID: %s", err)
			return nil, err
		}
		ne.NuvlaEdgeStatusId = types.NewNuvlaIDFromId(ne.nuvlaEdgeResource.NuvlaBoxStatus)
	}

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

	credCopy := *ne.Credentials.(*types.ApiKeyLogInParams)
	if ne.Credentials != nil && credCopy.Key != "" && credCopy.Secret != "" {
		encodeCredentialsIfNeeded(&credCopy)
	}

	f := NuvlaEdgeSessionFreeze{
		SessionOptions: ne.GetSessionOpts(),
		Credentials:    &credCopy,
		// If this point is reached, NuvlaEdgeID should never be nil or empty so if null pointer exception
		// is raised here, there is another issue
		NuvlaEdgeId: ne.NuvlaEdgeId.String(),
	}

	if ne.nuvlaEdgeResource != nil {
		f.InfraServiceId = ne.nuvlaEdgeResource.InfrastructureServiceGroup
	}

	if ne.NuvlaEdgeStatusId != nil {
		f.NuvlaEdgeStatusId = ne.NuvlaEdgeStatusId.String()
	}

	return f.Save(file)
}

func encodeCredentialsIfNeeded(creds *types.ApiKeyLogInParams) {
	if !strings.HasPrefix(creds.Key, "credential/") {
		// If the key is not a reference to a credential, we don't need to encode it
		log.Infof("API key is not a reference to a credential, skipping encoding")
		return
	}
	log.Infof("Encoding API key and secret...")
	creds.Key = base64.URLEncoding.EncodeToString([]byte(creds.Key))
	creds.Secret = base64.URLEncoding.EncodeToString([]byte(creds.Secret))
}

func decodeCredentialsIfNeeded(creds *types.ApiKeyLogInParams) {
	if strings.HasPrefix(creds.Key, "credential/") {
		// If the key is a reference to a credential, we don't need to decode it
		log.Infof("API key is a reference to a credential, skipping decoding")
		return
	}
	log.Infof("Decoding API key and secret...")
	key, err := base64.URLEncoding.DecodeString(creds.Key)
	if err != nil {
		log.Errorf("Error decoding api key: %s", err)
	}
	secret, err := base64.URLEncoding.DecodeString(creds.Secret)
	if err != nil {
		log.Errorf("Error decoding secret key: %s", err)
	}
	creds.Key = string(key)
	creds.Secret = string(secret)
}
