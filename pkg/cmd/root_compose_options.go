package cmd

import (
	"github.com/docker/compose/v2/cmd/compose"
	"github.com/spf13/pflag"
)

type composeOptions compose.ProjectOptions

var defaultComposeOptions = composeOptions{
	Profiles: []string{},
	ProjectName: "",
	ConfigPaths: []string{"./docker-compose.yml"},
	EnvFiles: nil,
	ProjectDir: "",
}

func createNewComposeOptions() *composeOptions {
	res := defaultComposeOptions
	return &res
}

func (opts *composeOptions) equals(arg *composeOptions) bool {
	if arg == nil {
		return false
	}
	if opts.ProjectName != arg.ProjectName || opts.ProjectDir != arg.ProjectDir {
		return false
	}
	return stringArrayEquals(opts.Profiles, arg.Profiles) &&
		stringArrayEquals(opts.ConfigPaths, arg.ConfigPaths) &&
		stringArrayEquals(opts.EnvFiles, arg.EnvFiles)
}

// ref: https://github.com/docker/compose/blob/d10a179f3e451f8b03fd99271f011c34bc31bedb/cmd/compose/compose.go#L157-L167
func (opts *composeOptions) set(f *pflag.FlagSet) {
	f.StringArrayVar(
		&opts.Profiles,
		"profile",
		defaultComposeOptions.Profiles,
		"Specify a profile to enable",
	)
	f.StringVarP(
		&opts.ProjectName,
		"project-name",
		"p",
		defaultComposeOptions.ProjectName,
		"Project name",
	)
	f.StringArrayVarP(
		&opts.ConfigPaths,
		"file",
		"f",
		defaultComposeOptions.ConfigPaths,
		"Compose configuration files",
	)
	f.StringArrayVar(
		&opts.EnvFiles,
		"env-file",
		defaultComposeOptions.EnvFiles,
		"Specify an alternate environment file",
	)
	f.StringVar(
		&opts.ProjectDir,
		"project-directory",
		defaultComposeOptions.ProjectDir,
		"Specify an alternate working directory\n(default: the path of the, first specified, Compose file)",
	)
}
