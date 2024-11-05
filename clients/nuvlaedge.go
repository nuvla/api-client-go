package clients

import (
	"context"
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
	IrsV2       string                   `json:"irs2"`
	IrsV1       string                   `json:"irs"`

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

	return nil
}

func (sf *NuvlaEdgeSessionFreeze) Save(file string) error {

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
	Irs               string

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
func (ne *NuvlaEdgeClient) LogIn(creds types.ApiKeyLogInParams) error {
	err := ne.LoginApiKeys(creds.Key, creds.Secret)
	if err != nil {
		log.Errorf("Error logging in with api keys: %s", err)
		return err
	}

	log.Infof("Logging in with api keys... Success.")
	return nil
}

// Activate Operation
func (ne *NuvlaEdgeClient) Activate(ctx context.Context) (types.ApiKeyLogInParams, error) {
	log.Infof("Activating NuvlaEdge...%v", ne.NuvlaEdgeId)
	res, err := ne.Operation(ctx, ne.NuvlaEdgeId.String(), "activate", nil)
	if err != nil {
		log.Errorf("Error activating NuvlaEdge: %s", err)
		return types.ApiKeyLogInParams{}, err
	}
	if res.StatusCode != 200 && res.StatusCode != 201 {
		return types.ApiKeyLogInParams{}, fmt.Errorf("activation failed with status code: %v", res.StatusCode)
	}

	creds, err := extractCredentialsFromActivateResponse(res)
	if err != nil {
		return types.ApiKeyLogInParams{}, fmt.Errorf("error extracting credentials from activate response: %s", err)
	}

	return *creds, nil
}

// Commission operations
func (ne *NuvlaEdgeClient) Commission(ctx context.Context, data map[string]interface{}) error {
	log.Debugf("Commissioning NuvlaEdge with payload %v", data)
	res, err := ne.Operation(ctx, ne.NuvlaEdgeId.String(), "commission", data)
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
	err = ne.UpdateResource(ctx)
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
func (ne *NuvlaEdgeClient) Telemetry(ctx context.Context, data interface{}, Select []string) (*http.Response, error) {
	log.Debugf("Sending telemetry data to NuvlaEdge with payload %v", data)
	if ne.nuvlaEdgeResource.NuvlaBoxStatus == "" || ne.NuvlaEdgeStatusId == nil {
		err := ne.UpdateResourceSelect(ctx, []string{"nuvlabox-status"})
		if err != nil {
			log.Errorf("Error sending Telemetry, cannot find NuvlaBoxStatus ID: %s", err)
			return nil, err
		}
		ne.NuvlaEdgeStatusId = types.NewNuvlaIDFromId(ne.nuvlaEdgeResource.NuvlaBoxStatus)
	}

	res, err := ne.Put(ctx, ne.NuvlaEdgeStatusId.String(), data, Select)
	if err != nil {
		log.Errorf("Error sending telemetry data to Nuvla: %s", err)
		return nil, err
	}
	return res, nil
}

// Heartbeat operation
func (ne *NuvlaEdgeClient) Heartbeat(ctx context.Context) (*http.Response, error) {
	log.Debug("Sending heartbeat to NuvlaEdge...")

	res, err := ne.Operation(ctx, ne.NuvlaEdgeId.String(), "heartbeat", nil)
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

func (ne *NuvlaEdgeClient) UpdateResourceSelect(ctx context.Context, selects []string) error {
	res, err := ne.Get(ctx, ne.NuvlaEdgeId.Id, selects)
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

func (ne *NuvlaEdgeClient) UpdateResource(ctx context.Context) error {
	return ne.UpdateResourceSelect(ctx, nil)
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
		IrsV2:          ne.Irs,
		// If this point is reached, NuvlaEdgeID should never be nil or empty so if null pointer exception
		// is raised here, there is another issue
		NuvlaEdgeId: ne.NuvlaEdgeId.String(),
	}

	// Keep credentials if available for backwards compatibility
	c, ok := ne.Credentials.(*types.ApiKeyLogInParams)
	if ok {

		f.Credentials = c
	}

	if ne.nuvlaEdgeResource != nil {
		f.InfraServiceId = ne.nuvlaEdgeResource.InfrastructureServiceGroup
	}

	if ne.NuvlaEdgeStatusId != nil {
		f.NuvlaEdgeStatusId = ne.NuvlaEdgeStatusId.String()
	}

	return f.Save(file)
}
