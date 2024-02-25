package cmd

import (
	"github.com/at0x0ft/pathru/pkg/pathru"
	"github.com/docker/compose/v2/cmd/compose"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	COMPOSE_PROJECT_OPTIONS_DEFAULT_CONFIG_PATH = "./docker-compose.yml"
)

func NewRootCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "pathru",
		Short: "Command pass-through helper with path conversion",
		Long: `pathru is a CLI command for help executing command in external container.
Usage: pathru <runtime service name> <execute command> -- [command arguments & options]`,
		RunE: runBody,
	}
}

func runBody(cmd *cobra.Command, args []string) error {
	opts := parseComposeOptions(cmd.Flags())
	// TODO: process args
	return pathru.Process(opts, args)
}

// ref: https://github.com/docker/compose/blob/d10a179f3e451f8b03fd99271f011c34bc31bedb/cmd/compose/compose.go#L157-L167
func parseComposeOptions(f *pflag.FlagSet) *compose.ProjectOptions {
	opts := compose.ProjectOptions{}
	f.StringArrayVar(&opts.Profiles, "profile", []string{}, "Specify a profile to enable")
	f.StringVarP(&opts.ProjectName, "project-name", "p", "", "Project name")
	f.StringArrayVarP(&opts.ConfigPaths, "file", "f", []string{COMPOSE_PROJECT_OPTIONS_DEFAULT_CONFIG_PATH}, "Compose configuration files")
	f.StringArrayVar(&opts.EnvFiles, "env-file", nil, "Specify an alternate environment file")
	f.StringVar(&opts.ProjectDir, "project-directory", "", "Specify an alternate working directory\n(default: the path of the, first specified, Compose file)")
	return &opts
}
