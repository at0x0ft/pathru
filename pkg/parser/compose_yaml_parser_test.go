package parser

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

type ParseComposeYamlSuccessTestCase struct {
	content string
	expected map[string]mount.BindMount
}

func providerTestParseComposeYamlSuccess() map[string]ParseComposeYamlSuccessTestCase {
	return map[string]ParseComposeYamlSuccessTestCase{
		"short syntax normal case": {
			`
services:
  base_shell:
    image: example/base_shell
    volumes:
      - type: volume
        source: db-data
        target: /data
        volume:
          nocopy: true
      # replace imagine filepath for existing target mount path for testing
      # - ./src:/workspace
      - /dev/null:/workspace
volumes:
  db-data:
`,
			map[string]mount.BindMount{
				"base_shell": mount.BindMount{
					"/dev/null",
					"/workspace",
				},
			},
		},
		"long syntax normal case": {
			`
services:
  base_shell:
    image: example/base_shell
    volumes:
      - .:/workspace
  golang:
    image: golang:1.22
    volumes:
      - type: bind
        source: /home/testuser/Programming/test_project/golang
        target: /go/src
`,
			map[string]mount.BindMount{
				"base_shell": mount.BindMount{
					".",
					"/workspace",
				},
				"golang": mount.BindMount{
					"/home/testuser/Programming/test_project/golang",
					"/go/src",
				},
			},
		},
		"short syntax not found bind case": {
			`
services:
  base_shell:
    image: example/base_shell
    volumes:
      - src:/workspace
`,
			map[string]mount.BindMount{},
		},
	}
}

func TestParseComposeYamlSuccess(t * testing.T) {
	cases := providerTestParseComposeYamlSuccess()
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			p := ComposeYamlParser{c.content}
			actual, err := p.Parse()
			if err != nil {
				t.Error(err)
			}

			if el, al := len(c.expected), len(actual); el != al {
				t.Errorf(
					"parsed service counts do not match [expected = \"%s\", actual = \"%s\"]",
					c.expected,
					actual,
				)
			}
			for k, em := range c.expected {
				if am, ok := actual[k]; !ok {
					t.Errorf(
						"parsed mounts do not have expected service [\"%s\"]",
						k,
					)
				} else if em != am {
					t.Errorf(
						"parsed mounts do not match [service = \"%s\", expected = \"%s\", actual = \"%s\"]",
						k,
						em,
						am,
					)
				}
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
