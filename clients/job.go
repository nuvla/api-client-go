package clients

import (
	"encoding/json"
	"errors"
	nuvla "github.com/nuvla/api-client-go"
	"github.com/nuvla/api-client-go/clients/resources"
	"github.com/nuvla/api-client-go/types"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type NuvlaJobClient struct {
	*nuvla.NuvlaClient

	jobId       *types.NuvlaID
	jobResource *resources.JobResource
}

func NewJobClient(jobId string, client *nuvla.NuvlaClient) *NuvlaJobClient {
	if client == nil {
		panic("Client should not be nil")
	}

	log.Infof("Job client Endpoint: %s", client.String())
	return &NuvlaJobClient{
		NuvlaClient: client,
		jobId:       types.NewNuvlaIDFromId(jobId),
		jobResource: &resources.JobResource{},
	}
}

func (jc *NuvlaJobClient) UpdateResource() error {
	res, err := jc.Get(jc.jobId.Id, nil)
	if err != nil {
		log.Infof("Error updating Deployment resource %s", jc.jobId)
		return nil
	}

	if jc.jobResource == nil {
		jc.jobResource = &resources.JobResource{}
	}

	b, err := json.Marshal(res.Data)
	if err != nil {
		log.Errorf("Error marshaling deployment resource response data to bytes")
		return err
	}

	err = json.Unmarshal(b, jc.jobResource)
	s, _ := json.MarshalIndent(res.Data, "", "  ")
	log.Infof("JobResource: %s", string(s))
	if err != nil {
		log.Error("Error unmarshalling response into DeploymentResource structure")
		return err
	}
	log.Infof("Successfully updated deployment resource")
	return nil
}

func (jc *NuvlaJobClient) GetId() string {
	return jc.jobId.Id
}

func (jc *NuvlaJobClient) GetType() resources.NuvlaResourceType {
	return resources.JobType
}

func (jc *NuvlaJobClient) GetResourceMap() (map[string]interface{}, error) {
	var mapRes map[string]interface{}
	if err := MarshalResourceIntoMap(jc.jobResource, mapRes); err != nil {
		log.Errorf("Error marshaling DeploymentResource to map")
		return nil, err
	}
	return mapRes, nil
}

func (jc *NuvlaJobClient) GetActionName() string {
	return jc.jobResource.Action
}

func (jc *NuvlaJobClient) GetResource() *resources.JobResource {
	return jc.jobResource
}

func (jc *NuvlaJobClient) PrintResource() {
	if jc.jobResource == nil {
		log.Debugf("Resource empty")
		return
	}

	p, err := json.MarshalIndent(jc.jobResource, "", "  ")
	if err != nil {
		log.Debugf("Error Marshaling %s resource, cannot print", jc.GetType())
		return
	}

	log.Infof("%s resource: \n %s", jc.GetType(), string(p))
}

func (jc *NuvlaJobClient) SetProgress(progress int8) error {
	if progress < 0 || progress > 100 {
		log.Errorf("Progress value %d is not valid", progress)
		return nil
	}
	log.Debugf("Setting progress in %s to %d", jc.jobId.Id, progress)
	res, err := jc.Edit(jc.jobId.Id, map[string]interface{}{"progress": progress}, nil)
	if err != nil {
		log.Errorf("Error setting progress to %d: %s", progress, err)
		return err
	}
	PrintResponse(res)
	jc.jobResource.Progress = progress
	return nil
}

// Set Status message
func (jc *NuvlaJobClient) SetStatusMessage(message string) {
	res, err := jc.Edit(jc.jobId.Id, map[string]interface{}{"status-message": message}, nil)
	if err != nil {
		log.Errorf("Error setting status message %s: %s", message, err)
		return
	}
	PrintResponse(res)
	jc.jobResource.StatusMessage = message
}

// Set State
func (jc *NuvlaJobClient) SetState(state resources.JobState) {
	res, err := jc.Edit(jc.jobId.Id, map[string]interface{}{"state": state}, nil)
	if err != nil {
		log.Error("Error setting state %s", state)
	}
	PrintResponse(res)
	jc.jobResource.State = state
}

// SetInitialState sets both the state to RUNNING and the progress to 10
func (jc *NuvlaJobClient) SetInitialState() {
	log.Infof("Setting initial processing state...")
	res, err := jc.Edit(jc.jobId.Id, map[string]interface{}{"state": resources.StateRUNNING, "progress": 10}, nil)
	if err != nil {
		log.Errorf("Error setting initial state %s", err)
		return
	}
	PrintResponse(res)
	log.Infof("Setting initial processing state... Success.")
}

// SetSuccessState sets the state to SUCCESS and the progress to 100
func (jc *NuvlaJobClient) SetSuccessState() {
	log.Debugf("Setting success state...")
	res, err := jc.Edit(jc.jobId.Id, map[string]interface{}{"state": resources.StateSuccess, "progress": 100}, nil)
	if err != nil {
		log.Errorf("Error setting success state %s", err)
		return
	}
	PrintResponse(res)
	log.Debugf("Setting success state... Success.")
}

func (jc *NuvlaJobClient) GetCredentials() (string, string, error) {
	creds := jc.Credentials.GetParams()
	k, ok := creds["key"]
	if !ok {
		return "", "", errors.New("key not found in credentials")
	}
	s, ok := creds["secret"]
	if !ok {
		return "", "", errors.New("secret not found in credentials")
	}
	return k, s, nil
}

func PrintResponse(res *http.Response) {

	log.Infof("Processing response with jobs...")
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Errorf("Error reading response body: %s", err)
		return
	}

	defer res.Body.Close()

	var sample struct {
		Message string   `json:"message"`
		Jobs    []string `json:"jobs"`
	}
	err = json.Unmarshal(body, &sample)
	if err != nil {
		log.Errorf("Error unmarshaling response body: %s", err)
		return
	}

	bytes, _ := json.MarshalIndent(sample, "", "  ")
	log.Infof("Setting progress response %s", string(bytes))

}
