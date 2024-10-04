package clients

import (
	"context"
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

func (jc *NuvlaJobClient) UpdateResource(ctx context.Context) error {
	res, err := jc.Get(ctx, jc.jobId.Id, nil)
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

func (jc *NuvlaJobClient) UpdateJobStatus(ctx context.Context, opts JobStatusUpdateOpts) error {
	res, err := jc.Edit(ctx, jc.jobId.Id, opts.GetMap(), nil)
	if err != nil {
		log.Errorf("Error updating job status: %s", err)
		return err
	}
	PrintResponse(res)
	opts.UpdateJobResource(jc.jobResource)
	return nil
}

func (jc *NuvlaJobClient) SetProgress(ctx context.Context, progress int8) error {
	if progress < 0 || progress > 100 {
		log.Errorf("Progress value %d is not valid", progress)
		return nil
	}
	log.Debugf("Setting progress in %s to %d", jc.jobId.Id, progress)
	res, err := jc.Edit(ctx, jc.jobId.Id, map[string]interface{}{"progress": progress}, nil)
	if err != nil {
		log.Errorf("Error setting progress to %d: %s", progress, err)
		return err
	}
	PrintResponse(res)
	jc.jobResource.Progress = progress
	return nil
}

// Set Status message
func (jc *NuvlaJobClient) SetStatusMessage(ctx context.Context, message string) {
	res, err := jc.Edit(ctx, jc.jobId.Id, map[string]interface{}{"status-message": message}, nil)
	if err != nil {
		log.Errorf("Error setting status message %s: %s", message, err)
		return
	}
	PrintResponse(res)
	jc.jobResource.StatusMessage = message
}

// SetState
func (jc *NuvlaJobClient) SetState(ctx context.Context, state resources.JobState) {
	res, err := jc.Edit(ctx, jc.jobId.Id, map[string]interface{}{"state": state}, nil)
	if err != nil {
		log.Error("Error setting state %s", state)
	}
	PrintResponse(res)
	jc.jobResource.State = state
}

// SetInitialState sets both the state to RUNNING and the progress to 10
func (jc *NuvlaJobClient) SetInitialState(ctx context.Context) {
	log.Infof("Setting initial processing state...")
	res, err := jc.Edit(ctx, jc.jobId.Id, map[string]interface{}{"state": resources.StateRUNNING, "progress": 10}, nil)
	if err != nil {
		log.Errorf("Error setting initial state %s", err)
		return
	}
	PrintResponse(res)
	log.Infof("Setting initial processing state... Success.")
}

// SetSuccessState sets the state to SUCCESS and the progress to 100
func (jc *NuvlaJobClient) SetSuccessState(ctx context.Context) {
	log.Debugf("Setting success state...")
	res, err := jc.Edit(ctx, jc.jobId.Id, map[string]interface{}{"state": resources.StateSuccess, "progress": 100}, nil)
	if err != nil {
		log.Errorf("Error setting success state %s", err)
		return
	}
	PrintResponse(res)
	log.Debugf("Setting success state... Success.")
}

// SetFailedState sets the state to FAILED and the progress to 100
func (jc *NuvlaJobClient) SetFailedState(ctx context.Context, errMsg string) {
	log.Debugf("Setting failed state...")
	opts := JobStatusUpdateOpts{
		Progress:      100,
		StatusMessage: errMsg,
		State:         resources.StateFailed,
	}
	err := jc.UpdateJobStatus(ctx, opts)
	if err != nil {
		log.Errorf("Error setting failed state %s", err)
		return
	}
	log.Debugf("Setting failed state... Success.")
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
	defer res.Body.Close()
	if err != nil {
		log.Errorf("Error reading response body: %s", err)
		return
	}

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
	log.Debugf("Setting progress response %s", string(bytes))
}

type JobStatusUpdateOpts struct {
	Progress      int8               `json:"progress,omitempty"`
	StatusMessage string             `json:"status-message,omitempty"`
	State         resources.JobState `json:"state,omitempty"`
}

func (u *JobStatusUpdateOpts) GetMap() map[string]interface{} {
	m := make(map[string]interface{})
	if u.Progress != 0 {
		m["progress"] = u.Progress
	}
	if u.StatusMessage != "" {
		m["status-message"] = u.StatusMessage
	}
	if u.State != "" {
		m["state"] = u.State
	}
	return m
}

func (u *JobStatusUpdateOpts) UpdateJobResource(jr *resources.JobResource) {
	if jr == nil {
		return
	}
	if u.Progress != 0 {
		jr.Progress = u.Progress
	}
	if u.StatusMessage != "" {
		jr.StatusMessage = u.StatusMessage
	}
	if u.State != "" {
		jr.State = u.State
	}
}
