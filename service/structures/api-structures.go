package structures

type Project struct {
	Name string `json:"name"`
	Path string `json:"path"`

	Description string `json:"description"`

	GitURL    string `json:"git_url"`
	GitBranch string `json:"git_branch"`

	Dockerfile string `json:"dockerfile"`
	ImageName  string `json:"image_name"`
	ImageTag   string `json:"image_tag"`

	DeployKey bool `json:"deploy_key"`

	LastBuild string `json:"last_build"`
	Status    string `json:"status"`

	//RunArgs       string `json:"run_args"`
	//ContainerName string `json:"container_name"`
}

type GenericResponse struct {
	Status string `json:"status"`
}
