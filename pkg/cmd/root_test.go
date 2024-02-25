package cmd

import (
	"github.com/spf13/cobra"
	"os"
	"testing"
)

func TestMain(t *testing.M) {
	code := t.Run()
	os.Exit(code)
}

type SetComposeOptionsSuccessTestCase struct {
	args     []string
	expected []string
}

func providerTestSetComposeOptionsSuccess(t *testing.T) map[string]SetComposeOptionsSuccessTestCase {
	return map[string]SetComposeOptionsSuccessTestCase{
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

func TestSetComposeOptionsSuccess(t *testing.T) {
	cases := providerTestSetComposeOptionsSuccess(t)
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			// no parallelization because using os.Args (global variable)
			// t.Parallel()
			oldArgs := os.Args

			os.Args = append([]string{"command"}, c.args...)
			opts := &rootCommandOptions{
				composeOpts: composeOptions{},
				baseService: "",
			}
			cmd := cobra.Command{}
			opts.composeOpts.set(cmd.PersistentFlags())
			if err := cmd.Execute(); err != nil {
				t.Errorf(
					"command execute error: %v",
					err.Error(),
				)
			}

			if el, al := len(c.expected), len(opts.composeOpts.ConfigPaths); el != al {
				t.Errorf(
					"parsed config path counts do not match [expected = \"%s\", actual = \"%s\"]",
					c.expected,
					opts.composeOpts.ConfigPaths,
				)
			}
			for i, ep := range c.expected {
				ap := opts.composeOpts.ConfigPaths[i]
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

type BaseServiceOptionsSuccessTestCase struct {
	args     []string
	expected string
}

func providerTestBaseServiceOptionsSuccess(t *testing.T) map[string]BaseServiceOptionsSuccessTestCase {
	return map[string]BaseServiceOptionsSuccessTestCase{
		"base service specified case": {
			[]string{"-b", "base"},
			"base",
		},
		"no specified case": {
			[]string{},
			"base_shell",
		},
	}
}

func TestBaseServiceOptionsSuccess(t *testing.T) {
	cases := providerTestBaseServiceOptionsSuccess(t)
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			// no parallelization because using os.Args (global variable)
			// t.Parallel()
			oldArgs := os.Args

			os.Args = append([]string{"command"}, c.args...)
			opts := &rootCommandOptions{
				composeOpts: composeOptions{},
				baseService: "",
			}
			cmd := cobra.Command{}
			opts.setBaseServiceOption(cmd.PersistentFlags())
			if err := cmd.Execute(); err != nil {
				t.Errorf(
					"command execute error: %v",
					err.Error(),
				)
			}

			if c.expected != opts.baseService {
				t.Errorf(
					"parsed base service do not match [expected = \"%s\", actual = \"%s\"]",
					c.expected,
					opts.baseService,
				)
			}

			// finally restore os.Args (global variable)
			t.Cleanup(func() {
				os.Args = oldArgs
			})
		})
	}
}
