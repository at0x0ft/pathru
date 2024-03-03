package devcontainer

type DevcontainerConfig struct {
	DockerComposeFile []string `json:"dockerComposeFile"`
	Service           string   `json:"service"`
}
