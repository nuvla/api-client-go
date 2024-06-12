package resources

type NuvlaEdgeState string

const (
	NuvlaEdgeStateNew             NuvlaEdgeState = "NEW"
	NuvlaEdgeStateActivated       NuvlaEdgeState = "ACTIVATED"
	NuvlaEdgeStateCommissioned    NuvlaEdgeState = "COMMISSIONED"
	NuvlaEdgeStateDecommissioning NuvlaEdgeState = "DECOMMISSIONING"
	NuvlaEdgeStateDecommissioned  NuvlaEdgeState = "DECOMMISSIONED"
	NuvlaEdgeStateError           NuvlaEdgeState = "ERROR"
	NuvlaEdgeStateSuspended       NuvlaEdgeState = "SUSPENDED"
)

type NuvlaEdgeResource struct {
	CommonAttributesResource

	// Required
	State             NuvlaEdgeState `json:"state"`
	RefreshInterval   int            `json:"refresh-interval"`
	HeartbeatInterval int            `json:"heartbeat-interval"`
	Version           int            `json:"version"`
	Owner             string         `json:"owner"`

	// Optional
	Location              []float32 `json:"location"`
	VPNServerID           string    `json:"vpn-server-id"`
	SSHKeys               []string  `json:"ssh-keys"`
	Capabilities          []string  `json:"capabilities"`
	Online                bool      `json:"online"`
	InferredLocation      []float32 `json:"inferred-location"`
	NuvlaBoxEngineVersion string    `json:"nuvlabox-engine-version"`

	NuvlaBoxStatus             string `json:"nuvlabox-status"`
	InfrastructureServiceGroup string `json:"infrastructure-service-group"`
	CredentialApiKey           string `json:"credential-api-key"`
	HostLevelManagementApiKey  string `json:"host-level-management-api-key"`
}
