package argument

import (
	"github.com/docker/compose/v2/cmd/compose"
	"github.com/spf13/pflag"
)

const (
	COMPOSE_PROJECT_OPTIONS_DEFAULT_CONFIG_PATH = "./docker-compose.yml"
)

type ArgumentParser struct{}

func (ap *ArgumentParser) Parse() ([]string, []string, error) {

	// opts, err := cli.NewProjectOptions(
	// 	configPaths,
	// 	cli.WithOsEnv,
	// 	cli.WithDotEnv,
	// 	cli.WithConfigFileEnv,
	// 	cli.WithDefaultConfigPath,
	// )
	// if err != nil {
	// 	return nil, err
	// }

	// prj, err := cli.ProjectFromOptions(opts)
	// if err != nil {
	// 	return nil, err
	// }

	// res := make(map[string]mount.BindMount)
	// for _, n := range prj.ServiceNames() {
	// 	s, err := prj.GetService(n)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	for _, v := range s.Volumes {
	// 		if v.Type != types.VolumeTypeBind {
	// 			continue
	// 		}
	// 		res[n] = mount.BindMount{Source: v.Source, Target: v.Target}
	// 	}
	// }
	// fmt.Println(ap.addComposeProjectFlags((&pflag.FlagSet{})))
	return nil, nil, nil
}

// ref: https://github.com/docker/compose/blob/d10a179f3e451f8b03fd99271f011c34bc31bedb/cmd/compose/compose.go#L157-L167
func (ap *ArgumentParser) addComposeProjectFlags(f *pflag.FlagSet) *compose.ProjectOptions {
	opts := compose.ProjectOptions{}
	f.StringArrayVar(&opts.Profiles, "profile", []string{}, "Specify a profile to enable")
	f.StringVarP(&opts.ProjectName, "project-name", "p", "", "Project name")
	f.StringArrayVarP(&opts.ConfigPaths, "file", "f", []string{COMPOSE_PROJECT_OPTIONS_DEFAULT_CONFIG_PATH}, "Compose configuration files")
	f.StringArrayVar(&opts.EnvFiles, "env-file", nil, "Specify an alternate environment file")
	f.StringVar(&opts.ProjectDir, "project-directory", "", "Specify an alternate working directory\n(default: the path of the, first specified, Compose file)")
	return &opts
}

// func (ap *ArgumentParser) addDevcontainerFlag(f *pflag.FlagSet) {
// 	f.StringVarP(&opts.ProjectName, "project-name", "p", "", "Project name")
// }
