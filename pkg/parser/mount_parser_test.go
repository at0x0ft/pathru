package parser

import (
	"github.com/at0x0ft/pathru/pkg/entity"
	"github.com/docker/compose/v2/cmd/compose"
	"github.com/google/go-cmp/cmp"
	"os"
	"path/filepath"
	"testing"
)

func TestMain(t *testing.M) {
	code := t.Run()
	os.Exit(code)
}

type ParseSuccessTestCase struct {
	configPaths []string
	expected    map[string][]entity.BindMount
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
			t.FailNow()
		}
		absContexts[n] = absPath
	}

	return map[string]ParseSuccessTestCase{
		"short syntax normal case": {
			[]string{fixturePaths["short_normal"]},
			map[string][]entity.BindMount{
				"base_shell": []entity.BindMount{
					entity.BindMount{
						Source: filepath.Join(absContexts["short_normal"], "./src"),
						Target: "/workspace",
					},
				},
			},
		},
		"long syntax normal case": {
			[]string{fixturePaths["long_normal"]},
			map[string][]entity.BindMount{
				"base_shell": []entity.BindMount{
					entity.BindMount{
						Source: filepath.Join(absContexts["long_normal"], "."),
						Target: "/workspace",
					},
				},
				"golang": []entity.BindMount{
					entity.BindMount{
						Source: "/home/testuser/Programming/test_project/golang",
						Target: "/go/src",
					},
				},
			},
		},
		"short syntax not found bind case": {
			[]string{fixturePaths["no_bind"]},
			map[string][]entity.BindMount{},
		},
		"override case": {
			[]string{fixturePaths["long_normal"], fixturePaths["short_normal"], fixturePaths["no_bind"]},
			map[string][]entity.BindMount{
				"golang": []entity.BindMount{
					entity.BindMount{
						Source: "/home/testuser/Programming/test_project/golang",
						Target: "/go/src",
					},
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
			opts := &compose.ProjectOptions{ConfigPaths: c.configPaths}
			actual, err := p.Parse(opts)
			if err != nil {
				t.Error(err)
				t.FailNow()
			}

			if el, al := len(c.expected), len(actual); el != al {
				t.Errorf(
					"parsed service counts do not match [expected = \"%s\", actual = \"%s\"]",
					c.expected,
					actual,
				)
				t.FailNow()
			}
			for s, ems := range c.expected {
				ams, ok := actual[s]
				if !ok {
					t.Errorf(
						"actual does not have expected service [\"%s\"]",
						s,
					)
					t.FailNow()
				}

				if !cmp.Equal(ems, ams) {
					t.Errorf(
						"parsed mounts do not match [service = \"%s\", expected = \"%s\", actual = \"%s\"]",
						s,
						ems,
						ams,
					)
					t.FailNow()
				}
			}
		})
	}
}
