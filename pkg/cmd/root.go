package cmd

import (
	"fmt"
	"github.com/at0x0ft/pathru/pkg/mount"
	"github.com/at0x0ft/pathru/pkg/pathru"
	"github.com/docker/compose/v2/cmd/compose"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"strings"
)

const (
	COMPOSE_PROJECT_OPTIONS_DEFAULT_CONFIG_PATH = "./docker-compose.yml"
	DEFAULT_BASE_SERVICE                        = "base_shell"
)

type rootCommandBaseServiceOptions struct {
	name         string
	rawWorkDir   string
	workDirMount mount.BindMount
}

type rootCommandOptions struct {
	composeOpts     compose.ProjectOptions
	baseServiceOpts rootCommandBaseServiceOptions
}

// ref: https://github.com/docker/compose/blob/d10a179f3e451f8b03fd99271f011c34bc31bedb/cmd/compose/compose.go#L157-L167
func (opts *rootCommandOptions) setComposeOptions(f *pflag.FlagSet) {
	f.StringArrayVar(&opts.composeOpts.Profiles, "profile", []string{}, "Specify a profile to enable")
	f.StringVarP(&opts.composeOpts.ProjectName, "project-name", "p", "", "Project name")
	f.StringArrayVarP(&opts.composeOpts.ConfigPaths, "file", "f", []string{COMPOSE_PROJECT_OPTIONS_DEFAULT_CONFIG_PATH}, "Compose configuration files")
	f.StringArrayVar(&opts.composeOpts.EnvFiles, "env-file", nil, "Specify an alternate environment file")
	f.StringVar(&opts.composeOpts.ProjectDir, "project-directory", "", "Specify an alternate working directory\n(default: the path of the, first specified, Compose file)")
}

func (opts *rootCommandOptions) setBaseServiceOptions(f *pflag.FlagSet) {
	f.StringVarP(&(opts.baseServiceOpts.name), "base-service", "b", DEFAULT_BASE_SERVICE, "base current service name")
	f.StringVarP(&(opts.baseServiceOpts.rawWorkDir), "working-dir", "w", "", "working directory mount setting for base current service")
}

func NewRootCommand() *cobra.Command {
	opts := rootCommandOptions{
		composeOpts:     compose.ProjectOptions{},
		baseServiceOpts: rootCommandBaseServiceOptions{},
	}
	cmd := &cobra.Command{
		Use:   "pathru",
		Short: "Command pass-through helper with path conversion",
		Long: `pathru is a CLI command for help executing command in external container.
Usage: pathru <runtime service name> <execute command> -- [command arguments & options]`,
		RunE: opts.runBody,
	}
	opts.setComposeOptions(cmd.PersistentFlags())
	opts.setBaseServiceOptions(cmd.PersistentFlags())
	return cmd
}

func (opts *rootCommandOptions) runBody(cmd *cobra.Command, args []string) error {
	if err := opts.baseServiceOpts.parseOptions(); err != nil {
		return err
	}

	runService, runArgs, err := parseRunService(args)
	if err != nil {
		return err
	}
	convertedArgs, err := pathru.Process(&(opts.composeOpts), opts.baseServiceOpts.name, runService, runArgs)
	if err != nil {
		return err
	}

	fmt.Printf("%v\n", strings.Join(convertedArgs, " "))
	return nil
}

func parseRunService(args []string) (string, []string, error) {
	if len(args) < 1 {
		return "", nil, fmt.Errorf("not enough argument(s) are given")
	}
	return args[0], args[1:], nil
}

func (bopts *rootCommandBaseServiceOptions) parseOptions() error {
	paths := strings.Split(bopts.rawWorkDir, ":")
	if actual := len(paths); actual != 2 {
		return fmt.Errorf(
			"just 2 paths must be specified for working-dir option [actual count = \"%v\"]",
			actual,
		)
	}

	bopts.workDirMount = mount.BindMount{Source: paths[0], Target: paths[1]}
	return nil
}
