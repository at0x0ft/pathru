package parser

import (
	"github.com/at0x0ft/pathru/pkg/mount"
)

type MountParser interface {
	Parse() map[string]mount.BindMount
}
