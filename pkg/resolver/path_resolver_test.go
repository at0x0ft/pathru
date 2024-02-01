package resolver

import (
	"testing"
	"os"
	"github.com/at0x0ft/pathru/pkg/mount"
)

func TestMain(t *testing.M) {
	code := t.Run()
	os.Exit(code)
}

func getMounts() map[string]mount.BindMount {
	return map[string]mount.BindMount {
		"base_shell": mount.BindMount{
			"/home/testuser/Programming/test_project",
			"/workspace",
		},
		"golang": mount.BindMount{
			"/home/testuser/Programming/test_project/golang",
			"/go/src",
		},
	}
}

type ResolveSuccessTestCase struct {
	path string
	baseService string
	dstService string
	expected string
}

func providerTestResolveSuccess() map[string]ResolveSuccessTestCase {
	return map[string]ResolveSuccessTestCase{
		"normal case": {
			"/workspace/golang/pkg/cmd/root.go",
			"base_shell",
			"golang",
			"/go/src/pkg/cmd/root.go",
		},
		// "multiple subpath matched case": {
		// 	"/home/testuser",
		// 	"/workspace/hoge",
		// 	"/workspace/hoge/fuga",
		// 	"/home/testuser/fuga",
		// },
		// "mountPoints are directories case": {
		// 	"/home/testuser/",
		// 	"/workspace/",
		// 	"/workspace/hoge/fuga/piyo",
		// 	"/home/testuser/hoge/fuga/piyo",
		// },
		// "directory path converts to slash trailed path": {
		// 	"/home/testuser",
		// 	"/workspace",
		// 	"/workspace/hoge/",
		// 	"/home/testuser/hoge",
		// },
	}
}

func TestResolveSuccess(t * testing.T) {
	cases := providerTestResolveSuccess()
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			c := c
			t.Parallel()
			pr := PathResolver{getMounts()}
			actual, err := pr.Resolve(c.path, c.baseService, c.dstService)
			if err != nil {
				t.Error(err)
			}

			if c.expected != actual {
				t.Errorf(
					"resolved actual path does not match [expected = \"%s\", actual = \"%s\"]",
					c.expected,
					actual,
				)
			}
		})
	}
}

// type ResolveTargetToSourceFailTestCase struct {
// 	source string
// 	target string
// 	path string
// 	expectedMsg string
// }

// func providerTestResolveTargetToSourceFail() map[string]ResolveTargetToSourceFailTestCase {
// 	return map[string]ResolveTargetToSourceFailTestCase{
// 		"path not containing mount.source": {
// 			"/home/testuser",
// 			"/workspace",
// 			"/home/hoge/fuga",
// 			"given path cannot reach to its mount base path [given: \"/home/hoge/fuga\", base: \"/workspace\"]",
// 		},
// 	}
// }

// func TestResolveTargetToSourceFail(t * testing.T) {
// 	cases := providerTestResolveTargetToSourceFail()
// 	for name, c := range cases {
// 		t.Run(name, func(t *testing.T) {
// 			c := c
// 			t.Parallel()
// 			bm := BindMount{c.source, c.target}
// 			result, err := bm.ResolveTargetToSource(c.path)
// 			if err == nil {
// 				t.Errorf("no error thrown [result = \"%s\"]", result)
// 				return
// 			}

// 			actual := err.Error()
// 			if c.expectedMsg != actual {
// 				t.Errorf(
// 					"error message does not match [expected = \"%s\", actual = \"%s\"]",
// 					c.expectedMsg,
// 					actual,
// 				)
// 			}
// 		})
// 	}
// }