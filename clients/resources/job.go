package resources

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
