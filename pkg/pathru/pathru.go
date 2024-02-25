package pathru

import (
	"fmt"
	"github.com/at0x0ft/pathru/pkg/mount"
	"github.com/at0x0ft/pathru/pkg/resolver"
	"github.com/docker/compose/v2/cmd/compose"
	"os"
	"path/filepath"
)

func Process(
	opts *compose.ProjectOptions,
	baseService string,
	baseServiceWorkDirMount *mount.BindMount,
	runService string,
	args []string,
) ([]string, error) {
	mounts, err := (&mount.MountParser{}).Parse(opts)
	if err != nil {
		return nil, err
	}
	if err := validateWorkDirMount(baseService, baseServiceWorkDirMount, mounts); err != nil {
		return nil, err
	}

	return resolveArgs(args, baseService, runService, mounts)
}

func validateWorkDirMount(
	baseService string,
	baseServiceWorkDirMount *mount.BindMount,
	mounts map[string][]mount.BindMount,
) error {
	baseServiceMounts, ok := mounts[baseService]
	if !ok {
		return fmt.Errorf(
			"base service not found in mounts [service = \"%v\", mounts = \"%v\"]",
			baseService,
			mounts,
		)
	}

	for _, mount := range baseServiceMounts {
		if *baseServiceWorkDirMount == mount {
			return nil
		}
	}

	return fmt.Errorf(
		"base service working_dir mount does not exists in mounts [working_dir mount = \"%v\", mounts = \"%v\"]",
		*baseServiceWorkDirMount,
		baseServiceMounts,
	)
}

func pathExists(path string) bool {
	_, err := os.Stat(filepath.Clean(path))
	return err == nil
}

func resolveArgs(args []string, baseService, runtimeService string, mounts map[string][]mount.BindMount) ([]string, error) {
	r := resolver.PathResolver{Mounts: mounts}
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
		p, err := r.Resolve(absPath, baseService, runtimeService)
		if err != nil {
			return nil, err
		}
		res[i] = p
	}
	return res, nil
}
