package resources

import "fmt"

type ContainerResource struct {
}

type ContainerHostResources struct {
	Memory int     `json:"memory"`
	Cpus   float32 `json:"cpus"`
}

type ContainerPorts struct {
	TargetPort int `json:"target-port"`
	// Optional fields
	Protocol      string `json:"protocol"`
	PublishedPort int    `json:"published-port"`
}

type ContainerMounts struct {
	MountType string `json:"mount-type"`
	Target    string `json:"target"`
	// Optional fields
	Source        string              `json:"source"`
	ReadOnly      bool                `json:"read-only"`
	VolumeOptions []map[string]string `json:"volume-options"`
}

type RestartPolicy struct {
	Condition string `json:"condition"`
	// Optional fields
	Delay       int `json:"delay"`
	MaxAttempts int `json:"max-attempts"`
	Window      int `json:"window"`
}

type EnvironmentVariable struct {
	Name string `json:"name"`
	// Optional fields
	Description string `json:"description"`
	Required    bool   `json:"required"`
	Value       string `json:"value"`
}

func (ev *EnvironmentVariable) GetAsString() string {
	return fmt.Sprintf("%s=%s", ev.Name, ev.Value)
}
