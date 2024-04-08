package clients

import (
	"encoding/json"
	nuvla "github.com/nuvla/api-client-go"
	"github.com/nuvla/api-client-go/types"
	log "github.com/sirupsen/logrus"
)

type DeploymentState string

const (
	StateCreated    DeploymentState = "CREATED"
	StateStarted    DeploymentState = "STARTED"
	StateStarting   DeploymentState = "STARTING"
	StateStopped    DeploymentState = "STOPPED"
	StateStopping   DeploymentState = "STOPPING"
	StatePausing    DeploymentState = "PAUSING"
	StatePaused     DeploymentState = "PAUSED"
	StateSuspending DeploymentState = "SUSPENDING"
	StateSuspended  DeploymentState = "SUSPENDED"
	StateUpdating   DeploymentState = "UPDATING"
	StateUpdated    DeploymentState = "UPDATED"
	StatePending    DeploymentState = "PENDING"
	StateError      DeploymentState = "ERROR"
)

type DeploymentResource struct {
	// Required by Nuvla API server
	Module struct {
		Name          string                   `json:"name"`
		Path          string                   `json:"path"`
		ParentPath    string                   `json:"parent-path"`
		SubType       string                   `json:"subtype"`
		Versions      []map[string]interface{} `json:"versions"`
		Content       map[string]interface{}   `json:"content"`
		Valid         bool                     `json:"valid"`
		Compatibility string                   `json:"compatibility"`
		Href          string                   `json:"href"`
	} `json:"module"`

	State       DeploymentState `json:"state"`
	ApiEndpoint string          `json:"api-endpoint"`

	// Optional
	ApiCredentials struct {
		ApiKey    string `json:"api-key"`
		ApiSecret string `json:"api-secret"`
	} `json:"api-credentials"`

	Data                      map[string]interface{} `json:"data"`
	RegistriesCredentials     []string               `json:"registries-credentials"`
	InfrastructureService     string                 `json:"infrastructure-service"`
	Nuvlabox                  string                 `json:"nuvlabox"`
	ExecutionMode             string                 `json:"execution-mode"`
	CredentialName            string                 `json:"credential-name"`
	InfrastructureServiceName string                 `json:"infrastructure-service-name"`
	Id                        string                 `json:"id"`
}

type NuvlaDeploymentClient struct {
	*nuvla.NuvlaClient

	deploymentId *types.NuvlaID

	deploymentResource *DeploymentResource
}

func NewNuvlaDeploymentClient(deploymentId string, client *nuvla.NuvlaClient) *NuvlaDeploymentClient {
	return &NuvlaDeploymentClient{
		NuvlaClient:  client,
		deploymentId: types.NewNuvlaIDFromId(deploymentId),
	}
}

func (dc *NuvlaDeploymentClient) SetDeploymentState() error {
	return nil
}

func (dc *NuvlaDeploymentClient) UpdateResource() error {
	res, err := dc.Get(dc.deploymentId.Id, nil)
	if err != nil {
		log.Infof("Error updating Deployment resource %s", dc.deploymentId)
		return nil
	}

	if dc.deploymentResource == nil {
		dc.deploymentResource = &DeploymentResource{}
	}

	b, err := json.Marshal(res.Data)
	if err != nil {
		log.Errorf("Error marshaling deployment resource response data to bytes")
		return err
	}

	err = json.Unmarshal(b, dc.deploymentResource)
	if err != nil {
		log.Error("Error unmarshalling response into DeploymentResource structure")
		return err
	}
	log.Infof("Successfully updated deployment resource")
	return nil
}

func (dc *NuvlaDeploymentClient) GetId() string {
	return dc.deploymentId.Id
}

func (dc *NuvlaDeploymentClient) GetType() ClientResourceType {
	return DeploymentType
}

func (dc *NuvlaDeploymentClient) GetResourceMap() (map[string]interface{}, error) {

	var mapRes map[string]interface{}
	if err := MarshalResourceIntoMap(dc.deploymentResource, mapRes); err != nil {
		log.Errorf("Error marshaling DeploymentResource to map")
		return nil, err
	}

	return mapRes, nil
}

func (dc *NuvlaDeploymentClient) GetResource() *DeploymentResource {
	return dc.deploymentResource
}

func (dc *NuvlaDeploymentClient) PrintResource() {
	p, err := json.MarshalIndent(dc.deploymentResource, "", "  ")
	if err != nil {
		log.Debugf("Error Marshaling %s resource, cannot print", dc.GetType())
		return
	}

	log.Infof("%s resource: \n %s", dc.GetType(), string(p))
}

func (dc *NuvlaDeploymentClient) SetState(state DeploymentState) error {
	log.Infof("Setting deployment state %s...", state)
	res, err := dc.Edit(dc.GetId(), map[string]interface{}{"state": state}, nil)
	if err != nil {
		log.Errorf("Error setting deployment state %s: %s", state, err)
		return err
	}
	PrintResponse(res)
	log.Infof("Setting deployment state %s... Success.", state)
	return nil
}

func (dc *NuvlaDeploymentClient) SetStateStarted() error {
	return dc.SetState(StateStarted)
}
