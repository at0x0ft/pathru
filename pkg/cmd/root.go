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
	COMPOSE_PROJECT_OPTIONS_DEFAULT_PROJECT_DIR = ""
	DEFAULT_BASE_SERVICE                        = "base_shell"
)

type composeOptions compose.ProjectOptions

type rootCommandBaseServiceOptions struct {
	name         string
	rawWorkDir   string
	workDirMount mount.BindMount
}

type rootCommandOptions struct {
	composeOpts     composeOptions
	baseServiceOpts rootCommandBaseServiceOptions
}

// ref: https://github.com/docker/compose/blob/d10a179f3e451f8b03fd99271f011c34bc31bedb/cmd/compose/compose.go#L157-L167
func (opts *composeOptions) set(f *pflag.FlagSet) {
	f.StringArrayVar(&opts.Profiles, "profile", []string{}, "Specify a profile to enable")
	f.StringVarP(&opts.ProjectName, "project-name", "p", "", "Project name")
	f.StringArrayVarP(&opts.ConfigPaths, "file", "f", []string{COMPOSE_PROJECT_OPTIONS_DEFAULT_CONFIG_PATH}, "Compose configuration files")
	f.StringArrayVar(&opts.EnvFiles, "env-file", nil, "Specify an alternate environment file")
	f.StringVar(&opts.ProjectDir, "project-directory", COMPOSE_PROJECT_OPTIONS_DEFAULT_PROJECT_DIR, "Specify an alternate working directory\n(default: the path of the, first specified, Compose file)")
}

func (opts *rootCommandBaseServiceOptions) set(f *pflag.FlagSet) {
	f.StringVarP(&(opts.name), "base-service", "b", DEFAULT_BASE_SERVICE, "base current service name")
	f.StringVarP(&(opts.rawWorkDir), "working-dir", "w", "", "working directory mount setting for base current service")
}

func (opts *rootCommandBaseServiceOptions) parseOptions() error {
	paths := strings.Split(opts.rawWorkDir, ":")
	if actual := len(paths); actual != 2 {
		return fmt.Errorf(
			"just 2 paths must be specified for working-dir option [actual count = \"%v\"]",
			actual,
		)
	}

	opts.workDirMount = mount.BindMount{Source: paths[0], Target: paths[1]}
	return nil
}

func NewRootCommand() *cobra.Command {
	opts := rootCommandOptions{
		composeOpts:     composeOptions{},
		baseServiceOpts: rootCommandBaseServiceOptions{},
	}
	cmd := &cobra.Command{
		Use:   "pathru",
		Short: "Command pass-through helper with path conversion",
		Long: `pathru is a CLI command for help executing command in external container.
Usage: pathru <runtime service name> <execute command> -- [command arguments & options]`,
		RunE: opts.runBody,
	}
	opts.composeOpts.set(cmd.PersistentFlags())
	opts.baseServiceOpts.set(cmd.PersistentFlags())
	return cmd
}

func (opts *rootCommandOptions) runBody(cmd *cobra.Command, args []string) error {
	if err := opts.baseServiceOpts.parseOptions(); err != nil {
		return err
	}
	runService, runArgs, err := opts.parseRunService(args)
	if err != nil {
		return err
	}
	newOpts := opts.updateProjectDirOption()
	prjOpts := compose.ProjectOptions(newOpts.composeOpts)

	convertedArgs, err := pathru.Process(
		&prjOpts,
		newOpts.baseServiceOpts.name,
		&(newOpts.baseServiceOpts.workDirMount),
		runService,
		runArgs,
	)
	if err != nil {
		return err
	}

	fmt.Printf("%v\n", strings.Join(convertedArgs, " "))
	return nil
}

func (opts *rootCommandOptions) parseRunService(args []string) (string, []string, error) {
	if len(args) < 1 {
		return "", nil, fmt.Errorf("not enough argument(s) are given")
	}
	return args[0], args[1:], nil
}

func (opts *rootCommandOptions) updateProjectDirOption() *rootCommandOptions {
	newOpts := *opts
	if newOpts.composeOpts.ProjectDir == COMPOSE_PROJECT_OPTIONS_DEFAULT_PROJECT_DIR {
		newOpts.composeOpts.ProjectDir = opts.baseServiceOpts.workDirMount.Source
	}
	return &newOpts
}
