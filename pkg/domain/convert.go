package domain

import (
	"github.com/at0x0ft/pathru/pkg/entity"
	"github.com/at0x0ft/pathru/pkg/parser"
	"github.com/at0x0ft/pathru/pkg/resolver"
	"github.com/docker/compose/v2/cmd/compose"
	"os"
	"path/filepath"
)

const HOST_BASE_SERVICE = ""

func Convert(
	opts *compose.ProjectOptions,
	baseService string,
	args []string,
) ([]string, error) {
	runService := args[0]
	args = args[1:]
	mounts, err := (&parser.MountParser{}).Parse(opts)
	if err != nil {
		return nil, err
	}
	resolvedArgs, err := resolveArgs(args, baseService, runService, mounts)
	if err != nil {
		return nil, err
	}
	return append([]string{runService}, resolvedArgs...), nil
}

func pathExists(path string) bool {
	_, err := os.Stat(filepath.Clean(path))
	return err == nil
}

func resolveArgs(args []string, baseService, runtimeService string, mounts map[string][]entity.BindMount) ([]string, error) {
	r := &resolver.PathResolver{Mounts: mounts}
	res := make([]string, len(args))
	for i, arg := range args {
		if !pathExists(arg) {
			res[i] = arg
			continue
		}

		absPath, err := filepath.Abs(arg)
		if err != nil {
			return nil, err
		}
		p, err := resolve(r, absPath, baseService, runtimeService)
		if err != nil {
			return nil, err
		}
		res[i] = p
	}
	return res, nil
}

func resolve(
	r *resolver.PathResolver,
	absPath,
	baseService,
	runtimeService string,
) (string, error) {
	if baseService == HOST_BASE_SERVICE {
		return r.ResolveFromHost(absPath, runtimeService)
	}
	return r.Resolve(absPath, baseService, runtimeService)
}
