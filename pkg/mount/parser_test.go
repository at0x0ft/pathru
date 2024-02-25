package mount

import (
	"path/filepath"
	"testing"
)

type ParseSuccessTestCase struct {
	configPaths []string
	expected    map[string]BindMount
}

func providerTestParseSuccess(t *testing.T) map[string]ParseSuccessTestCase {
	fixturePaths := map[string]string{
		"short_normal": "./test_data/compose.short.normal.yml",
		"long_normal":  "./test_data/compose.long.normal.yml",
		"no_bind":      "./test_data/compose.no.bind.yml",
	}
	absContexts := make(map[string]string)
	for n, path := range fixturePaths {
		absPath, err := filepath.Abs(filepath.Dir(path))
		if err != nil {
			t.Errorf("%v", err.Error())
			return nil
		}
		absContexts[n] = absPath
	}

	return map[string]ParseSuccessTestCase{
		"short syntax normal case": {
			[]string{fixturePaths["short_normal"]},
			map[string]BindMount{
				"base_shell": BindMount{
					Source: filepath.Join(absContexts["short_normal"], "./src"),
					Target: "/workspace",
				},
			},
		},
		"long syntax normal case": {
			[]string{fixturePaths["long_normal"]},
			map[string]BindMount{
				"base_shell": BindMount{
					Source: filepath.Join(absContexts["long_normal"], "."),
					Target: "/workspace",
				},
				"golang": BindMount{
					Source: "/home/testuser/Programming/test_project/golang",
					Target: "/go/src",
				},
			},
		},
		"short syntax not found bind case": {
			[]string{fixturePaths["no_bind"]},
			map[string]BindMount{},
		},
		"override case": {
			[]string{fixturePaths["long_normal"], fixturePaths["short_normal"], fixturePaths["no_bind"]},
			map[string]BindMount{
				"golang": BindMount{
					Source: "/home/testuser/Programming/test_project/golang",
					Target: "/go/src",
				},
			},
		},
	}
}

func TestParseSuccess(t *testing.T) {
	cases := providerTestParseSuccess(t)
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			p := MountParser{}
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
