package cmd

import (
	"fmt"
	"github.com/at0x0ft/pathru/pkg/pathru"
	"github.com/docker/compose/v2/cmd/compose"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	COMPOSE_PROJECT_OPTIONS_DEFAULT_CONFIG_PATH = "./docker-compose.yml"
	DEFAULT_BASE_SERVICE                        = "base_shell"
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
	baseService := parseBaseService(cmd.Flags())
	runService, runArgs, err := parseRunService(args)
	if err != nil {
		return err
	}

	return pathru.Process(opts, baseService, runService, runArgs)
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

func parseBaseService(f *pflag.FlagSet) string {
	var res string
	f.StringVarP(&res, "base-service", "b", DEFAULT_BASE_SERVICE, "base current service name")
	return res
}

func parseRunService(args []string) (string, []string, error) {
	if len(args) < 1 {
		return "", nil, fmt.Errorf("[Error] not enough argument(s) are given")
	}
	return args[0], args[1:], nil
}
