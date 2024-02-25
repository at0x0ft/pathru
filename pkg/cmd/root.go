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

var (
	opts = rootCommandOptions{
		composeOpts: compose.ProjectOptions{},
		baseService: "",
	}
)

type rootCommandOptions struct {
	composeOpts compose.ProjectOptions
	baseService string
}

// ref: https://github.com/docker/compose/blob/d10a179f3e451f8b03fd99271f011c34bc31bedb/cmd/compose/compose.go#L157-L167
func (opts *rootCommandOptions) setComposeOptions(f *pflag.FlagSet) {
	f.StringArrayVar(&opts.composeOpts.Profiles, "profile", []string{}, "Specify a profile to enable")
	f.StringVarP(&opts.composeOpts.ProjectName, "project-name", "p", "", "Project name")
	f.StringArrayVarP(&opts.composeOpts.ConfigPaths, "file", "f", []string{COMPOSE_PROJECT_OPTIONS_DEFAULT_CONFIG_PATH}, "Compose configuration files")
	f.StringArrayVar(&opts.composeOpts.EnvFiles, "env-file", nil, "Specify an alternate environment file")
	f.StringVar(&opts.composeOpts.ProjectDir, "project-directory", "", "Specify an alternate working directory\n(default: the path of the, first specified, Compose file)")
}

func (opts *rootCommandOptions) setBaseServiceOption(f *pflag.FlagSet) {
	f.StringVarP(&opts.baseService, "base-service", "b", DEFAULT_BASE_SERVICE, "base current service name")
}

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pathru",
		Short: "Command pass-through helper with path conversion",
		Long: `pathru is a CLI command for help executing command in external container.
Usage: pathru <runtime service name> <execute command> -- [command arguments & options]`,
		RunE: runBody,
	}
	opts.setComposeOptions(cmd.PersistentFlags())
	opts.setBaseServiceOption(cmd.PersistentFlags())
	return cmd
}

func runBody(cmd *cobra.Command, args []string) error {
	runService, runArgs, err := parseRunService(args)
	if err != nil {
		return err
	}
	return pathru.Process(&opts.composeOpts, opts.baseService, runService, runArgs)
}

func parseRunService(args []string) (string, []string, error) {
	if len(args) < 1 {
		return "", nil, fmt.Errorf("[Error] not enough argument(s) are given")
	}
	return args[0], args[1:], nil
}
