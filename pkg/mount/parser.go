package mount

import (
	"github.com/compose-spec/compose-go/v2/cli"
	"github.com/compose-spec/compose-go/v2/types"
)

type MountParser struct{}

func (mp *MountParser) Parse(configPaths []string) (map[string]BindMount, error) {
	opts, err := cli.NewProjectOptions(
		configPaths,
		cli.WithOsEnv,
		cli.WithDotEnv,
		cli.WithConfigFileEnv,
		cli.WithDefaultConfigPath,
	)
	if err != nil {
		return nil, err
	}

	prj, err := cli.ProjectFromOptions(opts)
	if err != nil {
		return nil, err
	}

	res := make(map[string]BindMount)
	for _, n := range prj.ServiceNames() {
		s, err := prj.GetService(n)
		if err != nil {
			return nil, err
		}
		for _, v := range s.Volumes {
			if v.Type != types.VolumeTypeBind {
				continue
			}
			res[n] = BindMount{Source: v.Source, Target: v.Target}
		}
	}
	return res, nil
}
