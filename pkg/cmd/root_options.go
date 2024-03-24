package cmd

import (
	"github.com/at0x0ft/pathru/pkg/domain"
	"github.com/docker/compose/v2/cmd/compose"
	"github.com/spf13/pflag"
)

type rootOptions struct {
	baseService      *OptionData[string]
	composeOpts      *rootComposeOptions
	devcontainerOpts *rootDevcontainerOptions
}

func createRootOptions() *rootOptions {
	return &rootOptions{
		composeOpts:      &rootComposeOptions{},
		devcontainerOpts: &rootDevcontainerOptions{},
	}
}

func (opts *rootOptions) set(f *pflag.FlagSet) {
	opts.baseService = CreateStringPersistentOptionData(
		f,
		domain.HOST_BASE_SERVICE,
		"base-service",
		"b",
		"base current service name",
	)
	opts.composeOpts.set(f)
	opts.devcontainerOpts.set(f)
}

func (opts *rootOptions) parse() error {
	return opts.devcontainerOpts.parse()
}

func (opts *rootOptions) getProjectOptions() *compose.ProjectOptions {
	res := &compose.ProjectOptions{
		Profiles:    opts.composeOpts.profiles.Value(),
		ProjectName: opts.composeOpts.projectName.Value(),
		ConfigPaths: opts.composeOpts.configPaths.Value(),
		EnvFiles:    opts.composeOpts.envFiles.Value(),
		ProjectDir:  opts.composeOpts.projectDir.Value(),
	}
	if !opts.composeOpts.configPaths.IsSet() && opts.devcontainerOpts.state == DATA_IS_SET {
		res.ConfigPaths = opts.devcontainerOpts.data.dockerComposeFile
	}
	if !opts.composeOpts.projectDir.IsSet() && opts.devcontainerOpts.state == DATA_IS_SET {
		res.ProjectDir = opts.devcontainerOpts.data.localWorkspaceFolder
	}
	return res
}

func (opts *rootOptions) getBaseService() string {
	if !opts.baseService.IsSet() && opts.devcontainerOpts.state == DATA_IS_SET {
		return opts.devcontainerOpts.data.service
	}
	return opts.baseService.Value()
}
