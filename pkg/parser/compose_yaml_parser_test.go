package parser

import (
	"github.com/at0x0ft/pathru/pkg/mount"
	"os"
	"testing"
)

func TestMain(t *testing.M) {
	code := t.Run()
	os.Exit(code)
}

type ParseComposeYamlSuccessTestCase struct {
	content  string
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
      - ./src:/workspace
volumes:
  db-data:
`,
			map[string]mount.BindMount{
				"base_shell": mount.BindMount{
					"./src",
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
volumes:
  src:
`,
			map[string]mount.BindMount{},
		},
	}
}

func TestParseComposeYamlSuccess(t *testing.T) {
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
