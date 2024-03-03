package devcontainer

type DevcontainerConfig struct {
	DockerComposeFile []string          `json:"dockerComposeFile"`
	Service           string            `json:"service"`
	RemoteEnv         map[string]string `json:"remoteEnv"`
}

const LOCAL_WORKSPACE_FOLDER_VALUE = "${localWorkspaceFolder}"

func (c *DevcontainerConfig) FindLocalWorkspaceFolderEnvVar() string {
	for envVarName, value := range c.RemoteEnv {
		if value == LOCAL_WORKSPACE_FOLDER_VALUE {
			return envVarName
		}
	}
	return ""
}
