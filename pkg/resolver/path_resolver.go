package resolver

import (
	"fmt"
	"github.com/at0x0ft/pathru/pkg/mount"
)

type PathResolver struct {
	Mounts map[string]mount.BindMount
}

func (pr *PathResolver) Resolve(path, baseService, dstService string) (string, error) {
	var baseMount, dstMount mount.BindMount
	var ok bool
	var basePath string
	var err error

	if baseMount, ok = pr.Mounts[baseService]; !ok {
		return "", fmt.Errorf(
			"",
		)
	}
	if basePath, err = baseMount.ConvertTargetToSource(path); err != nil {
		return "", err
	}

	if dstMount, ok = pr.Mounts[dstService]; !ok {
		return "", fmt.Errorf(
			"",
		)
	}
	return dstMount.ConvertSourceToTarget(basePath)
}
