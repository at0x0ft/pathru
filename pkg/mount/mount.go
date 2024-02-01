package mount

import (
	"fmt"
	"path/filepath"
	"strings"
)

type BindMount struct {
	Source string
	Target string
}

func (bm *BindMount) ConvertSourceToTarget(pathInSource string) (string, error) {
	return bm.convert(pathInSource, bm.Source, bm.Target);
}

func (bm *BindMount) ConvertTargetToSource(pathInTarget string) (string, error) {
	return bm.convert(pathInTarget, bm.Target, bm.Source);
}

func (bm *BindMount) convert(path, base, dst string) (string, error) {
	bs := filepath.Clean(base)
	if !bm.isAncestor(bs, path) {
		return "", fmt.Errorf(
			"given path cannot reach to its mount base path [given: \"%s\", base: \"%s\"]",
			path,
			bs,
		)
	}
	rel, err := filepath.Rel(base, path)
	if err != nil {
		return "", err
	}
	return filepath.Join(dst, rel), nil
}

func (bm *BindMount) isAncestor(parentPath, childPath string) bool {
	return strings.HasPrefix(childPath, parentPath)
}
