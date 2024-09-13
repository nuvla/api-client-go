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

type DeploymentParameterResource struct {
	CommonAttributesResource
	NodeId string                 `json:"node-id,omitempty"`
	Value  string                 `json:"value,omitempty"`
	Acl    map[string]interface{} `json:"acl,omitempty"`
	// TODO: ATM, acl is only used here in the client but it belongs to CommonAttributesResource
}

func DefaultDeploymentParamResource() *DeploymentParameterResource {
	return &DeploymentParameterResource{
		CommonAttributesResource: CommonAttributesResource{},
	}
}

type DeploymentParamOptsFunc func(*DeploymentParameterResource)

func WithName(name string) DeploymentParamOptsFunc {
	return func(dp *DeploymentParameterResource) {
		dp.Name = name
	}
}

func WithDescription(description string) DeploymentParamOptsFunc {
	return func(dp *DeploymentParameterResource) {
		dp.Description = description
	}
}

func WithValue(value string) DeploymentParamOptsFunc {
	return func(dp *DeploymentParameterResource) {
		dp.Value = value
	}
}

func WithParent(parent string) DeploymentParamOptsFunc {
	return func(dp *DeploymentParameterResource) {
		dp.Parent = parent
	}
}

func WithAcl(acl map[string]interface{}) DeploymentParamOptsFunc {
	return func(dp *DeploymentParameterResource) {
		dp.Acl = acl
	}
}

func WithNodeId(nodeId string) DeploymentParamOptsFunc {
	return func(dp *DeploymentParameterResource) {
		dp.NodeId = nodeId
	}
}

func (dp *DeploymentParameterResource) New() NuvlaResource {
	return &DeploymentParameterResource{}
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
