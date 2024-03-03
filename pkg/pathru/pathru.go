package pathru

import (
	"github.com/at0x0ft/pathru/pkg/mount"
	"github.com/at0x0ft/pathru/pkg/resolver"
	"github.com/docker/compose/v2/cmd/compose"
	"os"
	"path/filepath"
)

const HOST_BASE_SERVICE = ""

func Process(
	opts *compose.ProjectOptions,
	baseService string,
	runService string,
	args []string,
) ([]string, error) {
	mounts, err := (&mount.MountParser{}).Parse(opts)
	if err != nil {
		return nil, err
	}
	return resolveArgs(args, baseService, runService, mounts)
}

func pathExists(path string) bool {
	_, err := os.Stat(filepath.Clean(path))
	return err == nil
}

func resolveArgs(args []string, baseService, runtimeService string, mounts map[string][]mount.BindMount) ([]string, error) {
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
