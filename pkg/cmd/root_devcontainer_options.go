package cmd

import (
	"github.com/at0x0ft/pathru/pkg/devcontainer"
	"github.com/spf13/pflag"
	"os"
)

type devcontainerOptions struct {
	path *OptionData[string]
}

type devcontainerData struct {
	dockerComposeFile    []string
	service              string
	localWorkspaceFolder string
}

func (opts *devcontainerOptions) set(f *pflag.FlagSet) {
	opts.path = CreateStringPersistentOptionData(
		f,
		"",
		"config-path",
		"c",
		"path to devcontainer.json",
	)
}

func (opts *devcontainerOptions) parse() (*devcontainerData, error) {
	if !opts.path.IsSet() {
		return nil, nil
	}

	config, err := devcontainer.Parse(opts.path.Value())
	if err != nil {
		return nil, err
	}

	localWorkspaceFolder := ""
	if envVarName := config.FindLocalWorkspaceFolderEnvVar(); envVarName != "" {
		localWorkspaceFolder = os.Getenv(envVarName)
	}
	return &devcontainerData{
		dockerComposeFile:    config.DockerComposeFile,
		service:              config.Service,
		localWorkspaceFolder: localWorkspaceFolder,
	}, nil
}
