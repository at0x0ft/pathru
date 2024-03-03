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

func (opts *rootCommandOptions) equals(arg *rootCommandOptions) bool {
	if opts.baseService != arg.baseService {
		return false
	}
	return opts.composeOptions.equals(&arg.composeOptions)
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

func (opts *rootCommandOptions) createWithOverWrite(
	co *composeOptions,
	do *devcontainerOptions,
) (*rootCommandOptions, error) {
	newOpts := createNewRootCommandOptions()
	if !do.equals(&defaultDevcontainerOptions) {
		newOpts.ConfigPaths = do.dockerComposeFile
		newOpts.baseService = do.service
	}
	if !opts.equals(&defaultRootCommandOptions) {
		newOpts.baseService = opts.baseService
	}
	if co != nil && !co.equals(&defaultComposeOptions) {
		newOpts.composeOptions = *co
	}

	return newOpts, nil;
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
	co.set(f)
	do.set(f)
	ro.setBaseServiceOption(f)
	return cmd
}
