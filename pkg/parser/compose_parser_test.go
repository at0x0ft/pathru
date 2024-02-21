package parser

import (
	"github.com/at0x0ft/pathru/pkg/mount"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestMain(t *testing.M) {
	code := t.Run()
	os.Exit(code)
}

type ParseComposeYamlSuccessTestCase struct {
	configPaths []string
	expected    map[string]mount.BindMount
}

func providerTestParseComposeYamlSuccess() map[string]ParseComposeYamlSuccessTestCase {
	fixturePaths := map[string]string{
		"short_normal": "./test_data/compose.short.normal.yml",
		"long_normal":  "./test_data/compose.long.normal.yml",
		"no_bind":      "./test_data/compose.no.bind.yml",
	}
	absContexts := make(map[string]string)
	for n, path := range fixturePaths {
		absPath, err := filepath.Abs(filepath.Dir(path))
		if err != nil {
			exitWithError(err)
			return nil
		}
		absContexts[n] = absPath
	}

	return map[string]ParseComposeYamlSuccessTestCase{
		"short syntax normal case": {
			[]string{fixturePaths["short_normal"]},
			map[string]mount.BindMount{
				"base_shell": mount.BindMount{
					Source: filepath.Join(absContexts["short_normal"], "./src"),
					Target: "/workspace",
				},
			},
		},
		"long syntax normal case": {
			[]string{fixturePaths["long_normal"]},
			map[string]mount.BindMount{
				"base_shell": mount.BindMount{
					Source: filepath.Join(absContexts["long_normal"], "."),
					Target: "/workspace",
				},
				"golang": mount.BindMount{
					Source: "/home/testuser/Programming/test_project/golang",
					Target: "/go/src",
				},
			},
		},
		"short syntax not found bind case": {
			[]string{fixturePaths["no_bind"]},
			map[string]mount.BindMount{},
		},
		"override case": {
			[]string{fixturePaths["long_normal"], fixturePaths["short_normal"], fixturePaths["no_bind"]},
			map[string]mount.BindMount{
				"golang": mount.BindMount{
					Source: "/home/testuser/Programming/test_project/golang",
					Target: "/go/src",
				},
			},
		},
	}
}

func TestParseComposeYamlSuccess(t *testing.T) {
	cases := providerTestParseComposeYamlSuccess()
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			p := ComposeParser{}
			actual, err := p.Parse(c.configPaths)
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

func exitWithError(err error) {
	log.Println(err)
	os.Exit(1)
}
