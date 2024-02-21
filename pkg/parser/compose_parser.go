package parser

import (
	"github.com/at0x0ft/pathru/pkg/mount"
	"github.com/compose-spec/compose-go/v2/cli"
	"github.com/compose-spec/compose-go/v2/types"
)

type ComposeParser struct{}

func (mp *ComposeParser) Parse(configPaths []string) (map[string]mount.BindMount, error) {
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

	res := make(map[string]mount.BindMount)
	for _, n := range prj.ServiceNames() {
		s, err := prj.GetService(n)
		if err != nil {
			return nil, err
		}
		for _, v := range s.Volumes {
			if v.Type != types.VolumeTypeBind {
				continue
			}
			res[n] = mount.BindMount{Source: v.Source, Target: v.Target}
		}
	}
	return res, nil
}
