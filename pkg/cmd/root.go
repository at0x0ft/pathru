package cmd

import (
	"fmt"
	"github.com/at0x0ft/pathru/pkg/pathru"
	"github.com/docker/compose/v2/cmd/compose"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"strings"
)

type rootCommandOptions struct {
	composeOptions
	baseService string
}

var defaultRootCommandOptions = rootCommandOptions{
	composeOptions: defaultComposeOptions,
	baseService: "base_shell",
}

func createNewRootCommandOptions() *rootCommandOptions {
	res := defaultRootCommandOptions
	return &res
}

func (opts *rootCommandOptions) setBaseServiceOption(f *pflag.FlagSet) {
	f.StringVarP(
		&opts.baseService,
		"base-service",
		"b",
		defaultRootCommandOptions.baseService,
		"base current service name",
	)
}

func (opts *rootCommandOptions) createWithMerge(
	co *composeOptions,
	do *devcontainerOptions,
) (*rootCommandOptions, error) {
	newOpts := createNewRootCommandOptions()
	newOpts.overwriteWithDevcontainerOptions(do)
	newOpts.overwriteWithBaseServiceOption(opts)
	newOpts.overwriteWithComposeOptions(co)
	return newOpts, nil;
}

func (opts *rootCommandOptions) overwriteWithDevcontainerOptions(do *devcontainerOptions) {
	if !stringArrayEquals(do.dockerComposeFile, defaultDevcontainerOptions.dockerComposeFile) {
		opts.ConfigPaths = do.dockerComposeFile
	}
	if do.service != defaultDevcontainerOptions.service {
		opts.baseService = do.service
	}
}

func (opts *rootCommandOptions) overwriteWithBaseServiceOption(ro *rootCommandOptions) {
	if ro.baseService != defaultRootCommandOptions.baseService {
		opts.baseService = ro.baseService
	}
}

func (opts *rootCommandOptions) overwriteWithComposeOptions(co *composeOptions) {
	if !stringArrayEquals(co.Profiles, defaultComposeOptions.Profiles) {
		opts.Profiles = co.Profiles
	}
	if co.ProjectName != defaultComposeOptions.ProjectName {
		opts.ProjectName = co.ProjectName
	}
	if !stringArrayEquals(co.ConfigPaths, defaultComposeOptions.ConfigPaths) {
		opts.ConfigPaths = co.ConfigPaths
	}
	if !stringArrayEquals(co.EnvFiles, defaultComposeOptions.EnvFiles) {
		opts.EnvFiles = co.EnvFiles
	}
	if co.ProjectDir != defaultComposeOptions.ProjectDir {
		opts.ProjectDir = co.ProjectDir
	}
}

func NewRootCommand() *cobra.Command {
	do := createNewDevcontainerOptions()
	co := createNewComposeOptions()
	ro := createNewRootCommandOptions()
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
			opts, err := ro.createWithMerge(co, parsedDevcontainerOptions)
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
	co.set(f)
	do.set(f)
	ro.setBaseServiceOption(f)
	return cmd
}
