package clients

import (
	"encoding/json"
	"errors"
	"fmt"
	nuvla "github.com/nuvla/api-client-go"
	"github.com/nuvla/api-client-go/clients/resources"
	"github.com/nuvla/api-client-go/common"
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

func (dc *NuvlaDeploymentClient) searchParameter(parentId, paramName, nodeId string) (*resources.DeploymentParameterResource, error) {
	filters := fmt.Sprintf("parent='%s' and name='%s'", parentId, paramName)
	if nodeId != "" {
		// Concatenate node-id='nodeId'
		filters = filters + fmt.Sprintf(" and node-id='%s'", nodeId)
	}

	// Search opts
	opts := &nuvla.SearchOptions{
		Filter: filters,
	}
	parameters, err := dc.Search(string(resources.DeploymentParameterType), opts)
	// TODO: See if err != nil should be consider as resource not found error
	if parameters.Count <= 0 {
		return nil, types.NewResourceNotFoundError(resources.DeploymentParameterType, "")
	}

	param := &resources.DeploymentParameterResource{}
	err = resources.NewResourceFromMap(parameters.Resources[0], param)
	if err != nil {
		return nil, err
	}
	return param, nil
}

func (dc *NuvlaDeploymentClient) SearchParameter(parentId, paramName, nodeId string) *resources.DeploymentParameterResource {
	param, err := dc.searchParameter(parentId, paramName, nodeId)
	if err != nil {
		log.Errorf("Error getting parameter %s: %s", paramName, err)
		return nil
	}
	return param
}

func (dc *NuvlaDeploymentClient) GetParameter(paramId string, paramSelect []string) (*resources.DeploymentParameterResource, error) {
	res, err := dc.Get(paramId, paramSelect)
	if err != nil {
		log.Errorf("Error getting parameter %s: %s", paramId, err)
		return nil, err
	}

	param := &resources.DeploymentParameterResource{}
	err = resources.NewResourceFromMap(res.Data, param)
	if err != nil {
		return nil, err
	}
	return param, nil
}

func (dc *NuvlaDeploymentClient) CreateParameter(userId string, opts ...resources.DeploymentParamOptsFunc) error {
	paramOpts := resources.DefaultDeploymentParamResource()
	for _, fn := range opts {
		fn(paramOpts)
	}

	if paramOpts.Parent == "" || paramOpts.Name == "" || userId == "" {
		log.Errorf("Parent, Name and UserId are required to create a parameter")
		jsOpts, _ := json.Marshal(paramOpts)
		var m map[string]interface{}
		_ = json.Unmarshal(jsOpts, &m)
		return types.NewResourceCreationError(resources.DeploymentParameterType, m)
	}

	aclMap := map[string]interface{}{
		"owners":   []string{"group/nuvla-admin"},
		"edit-acl": []string{userId},
	}
	paramOpts.Acl = aclMap

	// Create parameter
	var m map[string]interface{}
	jsOpts, _ := json.Marshal(paramOpts)
	_ = json.Unmarshal(jsOpts, &m)
	//ATM we don't use the ID of the parameter
	_, err := dc.Add(resources.DeploymentParameterType, m)
	return err
}

// UpdateParameter updates a parameter. If the parameter does not exist, it creates it.
func (dc *NuvlaDeploymentClient) UpdateParameter(userId string, opts ...resources.DeploymentParamOptsFunc) error {
	paramOpts := resources.DefaultDeploymentParamResource()
	for _, fn := range opts {
		fn(paramOpts)
	}
	// Try to get the parameter to retrieve the ID, if it does not exist, create it
	paramData, err := dc.searchParameter(paramOpts.Parent, paramOpts.Name, paramOpts.NodeId)
	if err != nil {
		var resourceNotFoundError *types.ResourceNotFoundError
		if errors.As(err, &resourceNotFoundError) {
			return dc.CreateParameter(userId, opts...)
		} else {
			log.Errorf("Error getting parameter %s: %s", paramOpts.Name, err)
		}
	}
	log.Infof("Updating parameter %s...", paramOpts.Name)
	_, err = dc.Edit(paramData.Id, common.GetCleanMapFromStruct(paramOpts), nil)
	if err != nil {
		log.Errorf("Error updating parameter %s: %s", paramOpts.Name, err)
		return err
	}
	log.Debugf("Parameter %s updated", paramOpts.Name)
	return nil
}
