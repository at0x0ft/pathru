package cmd

import (
	"github.com/at0x0ft/pathru/pkg/devcontainer"
	"github.com/spf13/pflag"
	"os"
)

type devcontainerOptions struct {
	path                 string
	dockerComposeFile    []string
	service              string
	localWorkspaceFolder string
}

var defaultDevcontainerOptions = devcontainerOptions{}

func createNewDevcontainerOptions() *devcontainerOptions {
	res := defaultDevcontainerOptions
	return &res
}

func (opts *devcontainerOptions) set(f *pflag.FlagSet) {
	f.StringVarP(&opts.path, "config-path", "c", "", "path to devcontainer.json")
}

func (opts *devcontainerOptions) parse() (*devcontainerOptions, error) {
	if opts.path == "" {
		return createNewDevcontainerOptions(), nil
	}

	config, err := devcontainer.Parse(opts.path)
	if err != nil {
		return nil, err
	}

	newOpts := *opts
	newOpts.dockerComposeFile = config.DockerComposeFile
	newOpts.service = config.Service
	envVarName := config.FindLocalWorkspaceFolderEnvVar()
	if envVarName != "" {
		newOpts.localWorkspaceFolder = os.Getenv(envVarName)
	}
	return &newOpts, nil
}
