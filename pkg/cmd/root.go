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

type composeOptions compose.ProjectOptions

var defaultComposeOptions = composeOptions{
	Profiles: []string{},
	ProjectName: "",
	ConfigPaths: []string{"./docker-compose.yml"},
	EnvFiles: nil,
	ProjectDir: "",
}

func createNewComposeOptions() *composeOptions {
	res := defaultComposeOptions
	return &res
}

func (opts *composeOptions) equals(arg *composeOptions) bool {
	if arg == nil {
		return false
	}
	if opts.ProjectName != arg.ProjectName || opts.ProjectDir != arg.ProjectDir {
		return false
	}
	return stringArrayEquals(opts.Profiles, arg.Profiles) &&
		stringArrayEquals(opts.ConfigPaths, arg.ConfigPaths) &&
		stringArrayEquals(opts.EnvFiles, arg.EnvFiles)
}

// ref: https://github.com/docker/compose/blob/d10a179f3e451f8b03fd99271f011c34bc31bedb/cmd/compose/compose.go#L157-L167
func (opts *composeOptions) set(f *pflag.FlagSet) {
	f.StringArrayVar(
		&opts.Profiles,
		"profile",
		defaultComposeOptions.Profiles,
		"Specify a profile to enable",
	)
	f.StringVarP(
		&opts.ProjectName,
		"project-name",
		"p",
		defaultComposeOptions.ProjectName,
		"Project name",
	)
	f.StringArrayVarP(
		&opts.ConfigPaths,
		"file",
		"f",
		defaultComposeOptions.ConfigPaths,
		"Compose configuration files",
	)
	f.StringArrayVar(
		&opts.EnvFiles,
		"env-file",
		defaultComposeOptions.EnvFiles,
		"Specify an alternate environment file",
	)
	f.StringVar(
		&opts.ProjectDir,
		"project-directory",
		defaultComposeOptions.ProjectDir,
		"Specify an alternate working directory\n(default: the path of the, first specified, Compose file)",
	)
}

type devcontainerOptions struct {
	path string
	dockerComposeFile []string
	service string
}

var defaultDevcontainerOptions = devcontainerOptions{}

func createNewDevcontainerOptions() *devcontainerOptions {
	res := defaultDevcontainerOptions
	return &res
}

func (opts *devcontainerOptions) equals(arg *devcontainerOptions) bool {
	if arg == nil {
		return false
	}
	if opts.path != arg.path || opts.service != arg.service {
		return false
	}
	return stringArrayEquals(opts.dockerComposeFile, arg.dockerComposeFile)
}

func (opts *devcontainerOptions) set(f *pflag.FlagSet) {
	f.StringVarP(&opts.path, "config-path", "c", "", "path to devcontainer.json")
}

func (opts *devcontainerOptions) parse() (*devcontainerOptions, error) {
	if opts.path == "" {
		return createNewDevcontainerOptions(), nil
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

func stringArrayEquals(ls1 []string, ls2 []string) bool {
	if ls1 == nil && ls2 == nil {
		return true
	} else if ls1 == nil || ls2 == nil {
		return false
	}
	if len(ls1) != len(ls2) {
		return false
	}
	for i, e1 := range ls1 {
		if e1 != ls2[i] {
			return false
		}
	}
	return true
}
