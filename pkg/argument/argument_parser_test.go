package argument

import (
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"testing"
)

func TestMain(t *testing.M) {
	code := t.Run()
	os.Exit(code)
}

type ParseSuccessTestCase struct {
	args     []string
	expected []string
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
		"single file specified case": {
			[]string{"-f", "./compose.yml"},
			[]string{"./compose.yml"},
		},
		"multiple files specified case": {
			[]string{"-f", "./docker-compose.yml", "-f", "./docker-compose.override.yml"},
			[]string{"./docker-compose.yml", "./docker-compose.override.yml"},
		},
		"using default path when no file specified case": {
			[]string{},
			[]string{"./docker-compose.yml"},
		},
	}
}

func TestParseSuccess(t *testing.T) {
	cases := providerTestParseSuccess(t)
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			oldArgs := os.Args

			os.Args = append([]string{"command"}, c.args...)
			p := ArgumentParser{}
			cmd := cobra.Command{}
			opts := p.addComposeProjectFlags(cmd.Flags())
			cmd.Execute()

			if el, al := len(c.expected), len(opts.ConfigPaths); el != al {
				t.Errorf(
					"parsed config path counts do not match [expected = \"%s\", actual = \"%s\"]",
					c.expected,
					opts.ConfigPaths,
				)
			}
			for i, ep := range c.expected {
				ap := opts.ConfigPaths[i]
				if ep != ap {
					t.Errorf(
						"parsed path does not match [expected = \"%s\", actual = \"%s\"]",
						ep,
						ap,
					)
				}
			}

			// finally restore os.Args (global variable)
			t.Cleanup(func() {
				os.Args = oldArgs
			})
		})
	}
}
