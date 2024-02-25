package resolver

import (
	"fmt"
	"github.com/at0x0ft/pathru/pkg/mount"
	"strings"
)

type PathResolver struct {
	Mounts map[string][]mount.BindMount
}

type pathGetter func(mount.BindMount) string

func (pr *PathResolver) Resolve(path, baseService, runService string) (string, error) {
	baseMount, err := pr.findMountFromTarget(path, baseService)
	if err != nil {
		return "", err
	}
	basePath, err := baseMount.ConvertTargetToSource(path)
	if err != nil {
		return "", err
	}

	dstMount, err := pr.findMountFromSource(basePath, runService)
	if err != nil {
		return "", err
	}
	return dstMount.ConvertSourceToTarget(basePath)
}

func (pr *PathResolver) findMountFromSource(path, service string) (*mount.BindMount, error) {
	getter := func(m mount.BindMount) string {
		return m.Source
	}
	return pr.findMountBase(path, service, getter)
}

func (pr *PathResolver) findMountFromTarget(path, service string) (*mount.BindMount, error) {
	getter := func(m mount.BindMount) string {
		return m.Target
	}
	return pr.findMountBase(path, service, getter)
}

func (pr *PathResolver) isAncestor(ancestorPath, childPath string) bool {
	return strings.HasPrefix(childPath, ancestorPath)
}

func (pr *PathResolver) findMountBase(path, service string, getter pathGetter) (*mount.BindMount, error) {
	mounts, ok := pr.Mounts[service]
	if !ok {
		return nil, fmt.Errorf(
			"service not found [service = \"%v\"]",
			service,
		)
	}

	for _, mount := range mounts {
		if pr.isAncestor(getter(mount), path) {
			return &mount, nil
		}
	}
	return nil, fmt.Errorf(
		"not found mount [service = \"%v\", path = \"%v\"]",
		service,
		path,
	)
}
