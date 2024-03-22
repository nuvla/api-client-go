package clients

import (
	"encoding/json"
	nuvla "github.com/nuvla/api-client-go"
	"github.com/nuvla/api-client-go/types"
	log "github.com/sirupsen/logrus"
)

type JobState string

const (
	StateQueued   JobState = "QUEUED"
	StateRUNNING  JobState = "RUNNING"
	StateFailed   JobState = "FAILED"
	StateCanceled JobState = "CANCELED"
	StateSuccess  JobState = "SUCCESS"
)

type JobResource struct {
	// Required
	State         JobState `json:"state"`
	Action        string   `json:"action"`
	Progress      int8     `json:"progress"`
	ExecutionMode string   `json:"execution-mode"`

	// Optional
	Version        int8 `json:"version"`
	TargetResource struct {
		Href string `json:"href"`
	} `json:"target-resource"`
	AffectedResources []struct {
		Href string `json:"href"`
	} `json:"affected-resources"`
	ReturnCode         int8   `json:"return-code"`
	StatusMessage      string `json:"status-message"`
	TimeOfStatusChange string `json:"time-of-status-change"` // TODO: Try using timestamps and automatic conversion
	ParentJob          string `json:"parent-job"`
	NestedJobs         string `json:"nested-jobs"`
	Priority           int16  `json:"priority"`
	Started            string `json:"started"` // TODO: Try using timestamps and automatic conversion
	Duration           int16  `json:"duration"`
	Expiry             string `json:"expiry"` // TODO: Try using timestamps and automatic conversion
	Output             string `json:"output"`
	Payload            string `json:"payload"` // JSON-compliant string to be passed to the job, such as execution arguments
}

type NuvlaJobClient struct {
	*nuvla.NuvlaClient

	jobId       *types.NuvlaID
	jobResource *JobResource
}

func NewJobClient(jobId string, client *nuvla.NuvlaClient) *NuvlaJobClient {
	if client == nil {
		panic("Client should not be nil")
	}

	log.Infof("Job client Endpoint: %s", client.String())
	return &NuvlaJobClient{
		NuvlaClient: client,
		jobId:       types.NewNuvlaIDFromId(jobId),
		jobResource: &JobResource{},
	}
}

func (jc *NuvlaJobClient) SetProgress(progress int) error {
	return nil
}

func (jc *NuvlaJobClient) UpdateResource() error {
	res, err := jc.Get(jc.jobId.Id, nil)
	if err != nil {
		log.Infof("Error updating Deployment resource %s", jc.jobId)
		return nil
	}

	if jc.jobResource == nil {
		jc.jobResource = &JobResource{}
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

func (jc *NuvlaJobClient) GetType() ClientResourceType {
	return JobType
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

func (jc *NuvlaJobClient) GetResource() *JobResource {
	return jc.jobResource
}

func (jc *NuvlaJobClient) PrintResource() {
	if jc.jobResource == nil {
		log.Debugf("Cannot pring ")
		return
	}

	p, err := json.MarshalIndent(jc.jobResource, "", "  ")
	if err != nil {
		log.Errorf("Error Marshaling %s resource", jc.GetType())
	}

	log.Infof("%s resource: \n %s", jc.GetType(), string(p))
}
