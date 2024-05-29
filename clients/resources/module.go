package resources

type ModuleResource struct {
	CommonAttributesResource
	Name          string                     `json:"name"`
	Path          string                     `json:"path"`
	ParentPath    string                     `json:"parent-path"`
	SubType       string                     `json:"subtype"`
	Versions      []map[string]interface{}   `json:"versions"`
	Content       *ModuleApplicationResource `json:"content"`
	Valid         bool                       `json:"valid"`
	Compatibility string                     `json:"compatibility"`
	Href          string                     `json:"href"`
}

type ModuleComponentResource struct {
	CommonAttributesResource
	Author        string                 `json:"author"`
	Architectures []string               `json:"architectures"`
	Image         map[string]interface{} `json:"image"` // Server definition of image is a map with constraint keys, might be useful to have a struct here
	// Optional fields
	Commit                 string                `json:"commit"`
	Memory                 int                   `json:"memory"`
	Cpus                   float32               `json:"cpus"`
	RestartPolicy          RestartPolicy         `json:"restart-policy"`
	Ports                  []ContainerPorts      `json:"ports"`
	Mounts                 []ContainerMounts     `json:"mounts"`
	PrivateRegistries      []string              `json:"private-registries"`
	RegistriesCredentials  []string              `json:"registries-credentials"`
	Urls                   [][]string            `json:"urls"`
	EnvironmentalVariables []EnvironmentVariable `json:"environmental-variables"`
	OutputParameters       []OutputParameter     `json:"output-parameters"`
}

type ModuleApplicationResource struct {
	CommonAttributesResource
	DockerCompose string `json:"docker-compose"`
	Author        string `json:"author"`
	// Optional fields
	Commit                string                  `json:"commit"`
	Urls                  [][]string              `json:"urls"`
	OutputParameters      []OutputParameter       `json:"output-parameters"`
	EnvironmentVariables  []EnvironmentVariable   `json:"environmental-variables"`
	PrivateRegistries     []string                `json:"private-registries"`
	RegistriesCredentials []string                `json:"registries-credentials"`
	Files                 []ModuleApplicationFile `json:"files"`
	RequiresUserRights    bool                    `json:"requires-user-rights"`
}

type ModuleApplicationFile struct {
	FileName    string `json:"file-name"`
	FileContent string `json:"file-content"`
}

type OutputParameter struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
