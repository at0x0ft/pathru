package cmd

import (
	"fmt"
	"github.com/at0x0ft/pathru/pkg/devcontainer"
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

type rootCommandOptions struct {
	composeOptions
	baseService string
}

type composeOptions compose.ProjectOptions

// ref: https://github.com/docker/compose/blob/d10a179f3e451f8b03fd99271f011c34bc31bedb/cmd/compose/compose.go#L157-L167
func (opts *composeOptions) set(f *pflag.FlagSet) {
	f.StringArrayVar(&opts.Profiles, "profile", []string{}, "Specify a profile to enable")
	f.StringVarP(&opts.ProjectName, "project-name", "p", "", "Project name")
	f.StringArrayVarP(&opts.ConfigPaths, "file", "f", []string{COMPOSE_PROJECT_OPTIONS_DEFAULT_CONFIG_PATH}, "Compose configuration files")
	f.StringArrayVar(&opts.EnvFiles, "env-file", nil, "Specify an alternate environment file")
	f.StringVar(&opts.ProjectDir, "project-directory", COMPOSE_PROJECT_OPTIONS_DEFAULT_PROJECT_DIR, "Specify an alternate working directory\n(default: the path of the, first specified, Compose file)")
}

type devcontainerOptions struct {
	path string
	dockerComposeFile []string
	service string
}

func (opts *devcontainerOptions) set(f *pflag.FlagSet) {
	f.StringVarP(&opts.path, "config-path", "c", "", "path to devcontainer.json")
}

func (opts *devcontainerOptions) parse() (*devcontainerOptions, error) {
	if opts.path == "" {
		return nil, nil
	}

	config, err := devcontainer.Parse(opts.path)
	if err != nil {
		return nil, err
	}

	newOpts := *opts
	newOpts.dockerComposeFile = config.DockerComposeFile
	newOpts.service = config.Service
	return &newOpts, nil
}

func (opts *rootCommandOptions) setBaseServiceOption(f *pflag.FlagSet) {
	f.StringVarP(&opts.baseService, "base-service", "b", DEFAULT_BASE_SERVICE, "base current service name")
}

func (opts *rootCommandOptions) createWithOverWrite(
	co *composeOptions,
	do *devcontainerOptions,
) (*rootCommandOptions, error) {
	newOpts := &rootCommandOptions{baseService: DEFAULT_BASE_SERVICE}
	newOpts.ConfigPaths = []string{COMPOSE_PROJECT_OPTIONS_DEFAULT_CONFIG_PATH}
	newOpts.ProjectDir = COMPOSE_PROJECT_OPTIONS_DEFAULT_PROJECT_DIR

	if do != nil {
		newOpts.ConfigPaths = do.dockerComposeFile
		newOpts.baseService = do.service
	}
	newOpts.baseService = opts.baseService
	if co != nil {
		newOpts.composeOptions = *co
	}

	return newOpts, nil;
}

func NewRootCommand() *cobra.Command {
	ro, co, do := &rootCommandOptions{}, &composeOptions{}, &devcontainerOptions{}
	cmd := &cobra.Command{
		Use:   "pathru",
		Short: "Command pass-through helper with path conversion",
		Long: `pathru is a CLI command for help executing command in external container.
Usage: pathru <runtime service name> <execute command> -- [command arguments & options]`,
		Args: func (cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf(
					"arguments must be given more than 1 [actual = \"%v\"]",
					args,
				)
			}
			return nil
		},
		RunE: func (cmd *cobra.Command, args []string) error {
			parsedDevcontainerOptions, err := do.parse()
			if err != nil {
				return err
			}
			opts, err := ro.createWithOverWrite(co, parsedDevcontainerOptions)
			if err != nil {
				return err
			}

			runService, runArgs := args[0], args[1:]
			prjOpts := compose.ProjectOptions(opts.composeOptions)
			convertedArgs, err := pathru.Process(
				&prjOpts,
				opts.baseService,
				runService,
				runArgs,
			)
			if err != nil {
				return err
			}

			fmt.Printf("%v\n", strings.Join(convertedArgs, " "))
			return nil
		},
	}
	f := cmd.PersistentFlags()
	ro.setBaseServiceOption(f)
	co.set(f)
	do.set(f)
	return cmd
}
