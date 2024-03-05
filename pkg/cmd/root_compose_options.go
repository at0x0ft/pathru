package cmd

import "github.com/spf13/pflag"

type composeOptions struct {
	profiles    *OptionData[[]string]
	projectName *OptionData[string]
	configPaths *OptionData[[]string]
	envFiles    *OptionData[[]string]
	projectDir  *OptionData[string]
}

// ref: https://github.com/docker/compose/blob/d10a179f3e451f8b03fd99271f011c34bc31bedb/cmd/compose/compose.go#L157-L167
func (opts *composeOptions) set(f *pflag.FlagSet) {
	opts.profiles = CreateStringArrayOptionData(
		f,
		[]string{},
		"profile",
		"Specify a profile to enable",
	)
	opts.projectName = CreateStringPersistentOptionData(
		f,
		"",
		"project-name",
		"p",
		"Project name",
	)
	opts.configPaths = CreateStringArrayPersistentOptionData(
		f,
		[]string{"./docker-compose.yml"},
		"file",
		"f",
		"Compose configuration files",
	)
	opts.envFiles = CreateStringArrayOptionData(
		f,
		nil,
		"env-file",
		"Specify an alternate environment file",
	)
	opts.projectDir = CreateStringOptionData(
		f,
		"",
		"project-directory",
		"Specify an alternate working directory\n(default: the path of the, first specified, Compose file)",
	)
}
