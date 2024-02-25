package mount

import (
	"github.com/compose-spec/compose-go/v2/cli"
	"github.com/compose-spec/compose-go/v2/types"
	"github.com/docker/compose/v2/cmd/compose"
)

type MountParser struct{}

func (mp *MountParser) Parse(opts *compose.ProjectOptions) (map[string][]BindMount, error) {
	convertedOpts, err := mp.convertComposeProjectOptions(opts)
	if err != nil {
		return nil, err
	}

	prj, err := cli.ProjectFromOptions(convertedOpts)
	if err != nil {
		return nil, err
	}

	res := make(map[string][]BindMount)
	for _, n := range prj.ServiceNames() {
		s, err := prj.GetService(n)
		if err != nil {
			return nil, err
		}
		for _, v := range s.Volumes {
			if v.Type != types.VolumeTypeBind {
				continue
			}
			res[n] = append(res[n], BindMount{Source: v.Source, Target: v.Target})
		}
	}
	return res, nil
}

// ref: https://github.com/docker/compose/blob/a7224411b4fb179ca47c2d4d86fb3a50a185c5ac/cmd/compose/compose.go#L288-L299
func (mp *MountParser) convertComposeProjectOptions(opts *compose.ProjectOptions) (*cli.ProjectOptions, error) {
	return cli.NewProjectOptions(
		opts.ConfigPaths,
		cli.WithWorkingDirectory(opts.ProjectDir),
		cli.WithOsEnv,
		cli.WithConfigFileEnv,
		cli.WithDefaultConfigPath,
		cli.WithEnvFiles(opts.EnvFiles...),
		cli.WithDotEnv,
		cli.WithDefaultProfiles(opts.Profiles...),
		cli.WithName(opts.ProjectName),
	)
}
