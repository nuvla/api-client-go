package resources

type DeploymentResource struct {
	CommonAttributesResource
	// Required by Nuvla API server
	Module      *ModuleResource `json:"module"`
	State       DeploymentState `json:"state"`
	ApiEndpoint string          `json:"api-endpoint"`

	// Optional
	ApiCredentials struct {
		ApiKey    string `json:"api-key"`
		ApiSecret string `json:"api-secret"`
	} `json:"api-credentials"`

	Data                      map[string]interface{} `json:"data"`
	RegistriesCredentials     []string               `json:"registries-credentials"`
	Owner                     string                 `json:"owner"`
	InfrastructureService     string                 `json:"infrastructure-service"`
	Nuvlabox                  string                 `json:"nuvlabox"`
	ExecutionMode             string                 `json:"execution-mode"`
	CredentialName            string                 `json:"credential-name"`
	InfrastructureServiceName string                 `json:"infrastructure-service-name"`
	Id                        string                 `json:"id"`
}

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
