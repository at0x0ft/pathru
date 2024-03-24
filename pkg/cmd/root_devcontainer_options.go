package cmd

import (
	"github.com/at0x0ft/pathru/pkg/devcontainer"
	"github.com/spf13/pflag"
	"os"
)

type RootDevcontainerOptionsState uint8

const (
	DATA_IS_NOT_PARSED RootDevcontainerOptionsState = iota
	DATA_IS_SET
	DATA_IS_NOT_SET
)

type rootDevcontainerOptions struct {
	state RootDevcontainerOptionsState
	path  *OptionData[string]
	data  *devcontainerData
}

type devcontainerData struct {
	dockerComposeFile    []string
	service              string
	localWorkspaceFolder string
}

func (opts *rootDevcontainerOptions) set(f *pflag.FlagSet) {
	opts.path = CreateStringPersistentOptionData(
		f,
		"",
		"config-path",
		"c",
		"path to devcontainer.json",
	)
}

func (opts *rootDevcontainerOptions) parse() error {
	if !opts.path.IsSet() {
		opts.state = DATA_IS_NOT_SET
		return nil
	}

	config, err := devcontainer.Parse(opts.path.Value())
	if err != nil {
		return err
	}

	localWorkspaceFolder := ""
	if envVarName := config.FindLocalWorkspaceFolderEnvVar(); envVarName != "" {
		localWorkspaceFolder = os.Getenv(envVarName)
	}
	opts.data = &devcontainerData{
		dockerComposeFile:    config.DockerComposeFile,
		service:              config.Service,
		localWorkspaceFolder: localWorkspaceFolder,
	}
	opts.state = DATA_IS_SET
	return nil
}
