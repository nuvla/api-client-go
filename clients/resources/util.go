package resources

const (
	NuvlaTimeStampFormat = "2006-01-02T15:04:05.000Z07:00"
)

type NuvlaResourceType string

func (n NuvlaResourceType) String() string {
	return string(n)
}

const (
	DeploymentType          NuvlaResourceType = "deployment"
	UserType                NuvlaResourceType = "user"
	NuvlaEdgeType           NuvlaResourceType = "nuvlaedge"
	NuvlaBoxType            NuvlaResourceType = "nuvlabox"
	JobType                 NuvlaResourceType = "job"
	DeploymentParameterType NuvlaResourceType = "deployment-parameter"
)
