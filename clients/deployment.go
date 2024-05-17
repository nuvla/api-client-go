package clients

import (
	"encoding/json"
	"fmt"
	nuvla "github.com/nuvla/api-client-go"
	"github.com/nuvla/api-client-go/clients/resources"
	"github.com/nuvla/api-client-go/types"
	log "github.com/sirupsen/logrus"
)

type NuvlaDeploymentClient struct {
	*nuvla.NuvlaClient

	deploymentId *types.NuvlaID

	deploymentResource *resources.DeploymentResource
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
		dc.deploymentResource = &resources.DeploymentResource{}
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

func (dc *NuvlaDeploymentClient) GetType() resources.NuvlaResourceType {
	return resources.DeploymentType
}

func (dc *NuvlaDeploymentClient) GetResourceMap() (map[string]interface{}, error) {

	var mapRes map[string]interface{}
	if err := MarshalResourceIntoMap(dc.deploymentResource, mapRes); err != nil {
		log.Errorf("Error marshaling DeploymentResource to map")
		return nil, err
	}

	return mapRes, nil
}

func (dc *NuvlaDeploymentClient) GetResource() *resources.DeploymentResource {
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

func (dc *NuvlaDeploymentClient) SetState(state resources.DeploymentState) error {
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
	return dc.SetState(resources.StateStarted)
}

func (dc *NuvlaDeploymentClient) GetParameter(paramId, paramName, nodeId, qSelect string) error {
	filters := fmt.Sprintf("parent='%s' and name='%s'", paramId, paramName)
	if nodeId != "" {
		// Concatenate node-id='nodeId'
		filters = filters + fmt.Sprintf(" and node-id='%s'", nodeId)
	}
	return nil
}
