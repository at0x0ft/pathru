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
	baseService *OptionData[string]
}

func (opts *rootCommandOptions) set(f *pflag.FlagSet) {
	opts.baseService = CreateStringPersistentOptionData(
		f,
		pathru.HOST_BASE_SERVICE,
		"base-service",
		"b",
		"base current service name",
	)
}

func NewRootCommand() *cobra.Command {
	ro, co, do := &rootCommandOptions{}, &composeOptions{}, &devcontainerOptions{}
	cmd := &cobra.Command{
		Use:   "pathru",
		Short: "Command pass-through helper with path conversion",
		Long: `pathru is a CLI command for help executing command in external container.
Usage: pathru <runtime service name> <execute command> -- [command arguments & options]`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf(
					"arguments must be given more than 1 [actual = \"%v\"]",
					args,
				)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			configData, err := do.parse()
			if err != nil {
				return err
			}

			prjOpts, baseService := mergeOptions(ro, co, configData)
			runService, runArgs := args[0], args[1:]
			convertedArgs, err := pathru.Process(
				prjOpts,
				baseService,
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
	ro.set(f)
	co.set(f)
	do.set(f)
	return cmd
}

func mergeOptions(
	r *rootCommandOptions,
	c *composeOptions,
	d *devcontainerData,
) (*compose.ProjectOptions, string) {
	return getProjectOptions(r, c, d), getBaseService(r, d)
}

func getProjectOptions(
	r *rootCommandOptions,
	c *composeOptions,
	d *devcontainerData,
) *compose.ProjectOptions {
	res := &compose.ProjectOptions{
		Profiles:    c.profiles.Value(),
		ProjectName: c.projectName.Value(),
		ConfigPaths: c.configPaths.Value(),
		EnvFiles:    c.envFiles.Value(),
		ProjectDir:  c.projectDir.Value(),
	}
	if !c.configPaths.IsSet() && d != nil {
		res.ConfigPaths = d.dockerComposeFile
	}
	if !c.projectDir.IsSet() && d != nil {
		res.ProjectDir = d.localWorkspaceFolder
	}
	return res
}

func getBaseService(r *rootCommandOptions, d *devcontainerData) string {
	if !r.baseService.IsSet() && d != nil {
		return d.service
	}
	return r.baseService.Value()
}
