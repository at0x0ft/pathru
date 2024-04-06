package resolver

import (
	"github.com/at0x0ft/pathru/pkg/entity"
	"os"
	"testing"
)

func TestMain(t *testing.M) {
	code := t.Run()
	os.Exit(code)
}

func getMounts() map[string][]entity.BindMount {
	return map[string][]entity.BindMount{
		"base_shell": []entity.BindMount{
			entity.BindMount{
				Source: "/etc/passwd",
				Target: "/etc/passwd",
			},
			entity.BindMount{
				Source: "/etc/group",
				Target: "/etc/group",
			},
			entity.BindMount{
				Source: "/home/testuser/Programming/test_project",
				Target: "/workspace",
			},
		},
		"golang": []entity.BindMount{
			entity.BindMount{
				Source: "/etc/passwd",
				Target: "/etc/passwd",
			},
			entity.BindMount{
				Source: "/etc/group",
				Target: "/etc/group",
			},
			entity.BindMount{
				Source: "/home/testuser/Programming/test_project/golang",
				Target: "/go/src",
			},
		},
	}
}

type ResolveSuccessTestCase struct {
	path        string
	baseService string
	dstService  string
	expected    string
}

func providerTestResolveSuccess() map[string]ResolveSuccessTestCase {
	return map[string]ResolveSuccessTestCase{
		"normal case": {
			"/workspace/golang/pkg/cmd/root.go",
			"base_shell",
			"golang",
			"/go/src/pkg/cmd/root.go",
		},
		"services are same case": {
			"/workspace/pkg/cmd/root.go",
			"base_shell",
			"base_shell",
			"/workspace/pkg/cmd/root.go",
		},
	}
}

func TestResolveSuccess(t *testing.T) {
	cases := providerTestResolveSuccess()
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			pr := PathResolver{getMounts()}
			actual, err := pr.Resolve(c.path, c.baseService, c.dstService)
			if err != nil {
				t.Errorf("%v", err)
				t.FailNow()
			}

			if c.expected != actual {
				t.Errorf(
					"resolved actual path does not match [expected = \"%s\", actual = \"%s\"]",
					c.expected,
					actual,
				)
				t.FailNow()
			}
		})
	}
}
