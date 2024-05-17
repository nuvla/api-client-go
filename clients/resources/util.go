package resources

const (
	NuvlaTimeStampFormat = "2006-01-02T15:04:05.000Z07:00"
)

type NuvlaResourceType string

const (
	DeploymentType NuvlaResourceType = "deployment"
	UserType       NuvlaResourceType = "user"
	NuvlaEdgeType  NuvlaResourceType = "nuvlaedge"
	NuvlaBoxType   NuvlaResourceType = "nuvlabox"
	JobType        NuvlaResourceType = "job"
)

func GetResourceFromString(resource string) NuvlaResourceType {
	switch resource {
	case "deployment":
		return DeploymentType
	case "user":
		return UserType
	case "nuvlaedge":
		return NuvlaEdgeType
	case "nuvlabox":
		return NuvlaBoxType
	case "job":
		return JobType
	default:
		return ""
	}
}
