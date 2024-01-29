package resolver

import (
	"fmt"
	"path/filepath"
	"strings"
)

type MountPointResolver struct {
	srcMountPoint string
	dstMountPoint string
}

func (m *MountPointResolver) Resolve(path string) (string, error) {
	srcMountPoint := filepath.Clean(m.srcMountPoint)
	if !m.isAncestor(srcMountPoint, path) {
		return "", fmt.Errorf(
			"path cannot reach srcMountPoint [path: \"%s\", srcMountPoint: \"%s\"]",
			path,
			srcMountPoint,
		)
	}
	relPath, err := filepath.Rel(srcMountPoint, path)
	if err != nil {
		return "", err
	}
	return filepath.Join(m.dstMountPoint, relPath), nil
}

func (m *MountPointResolver) isAncestor(parentPath, childPath string) bool {
	return strings.HasPrefix(childPath, parentPath)
}
