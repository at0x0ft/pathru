package cmd

import (
	"github.com/at0x0ft/pathru/pkg/mount"
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
				composeOpts:     composeOptions{},
				baseServiceOpts: rootCommandBaseServiceOptions{},
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

type BaseServiceOptionsSuccessTestCaseExpectedValues struct {
	name         string
	workDirMount mount.BindMount
}
type BaseServiceOptionsSuccessTestCase struct {
	args     []string
	expected BaseServiceOptionsSuccessTestCaseExpectedValues
}

func providerTestBaseServiceOptionsSuccess(t *testing.T) map[string]BaseServiceOptionsSuccessTestCase {
	return map[string]BaseServiceOptionsSuccessTestCase{
		"simple case": {
			[]string{"-w", "/home/testuser/Programming:/workspace"},
			BaseServiceOptionsSuccessTestCaseExpectedValues{
				name: "base_shell",
				workDirMount: mount.BindMount{
					Source: "/home/testuser/Programming",
					Target: "/workspace",
				},
			},
		},
		"base service specified case": {
			[]string{"-w", "/tmp:/workspace/src", "-b", "base"},
			BaseServiceOptionsSuccessTestCaseExpectedValues{
				name: "base",
				workDirMount: mount.BindMount{
					Source: "/tmp",
					Target: "/workspace/src",
				},
			},
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
				composeOpts:     composeOptions{},
				baseServiceOpts: rootCommandBaseServiceOptions{},
			}
			cmd := cobra.Command{}
			opts.baseServiceOpts.set(cmd.PersistentFlags())
			if err := cmd.Execute(); err != nil {
				t.Errorf(
					"command execute error: %v",
					err.Error(),
				)
			}
			actual := &(opts.baseServiceOpts)
			if err := actual.parseOptions(); err != nil {
				t.Errorf("%v", err)
			}

			if c.expected.name != actual.name {
				t.Errorf(
					"parsed base service names do not match [expected = \"%s\", actual = \"%s\"]",
					c.expected.name,
					actual.name,
				)
			}
			if c.expected.workDirMount != actual.workDirMount {
				t.Errorf(
					"parsed base service working directory mounts do not match [expected = \"%s\", actual = \"%s\"]",
					c.expected.workDirMount,
					actual.workDirMount,
				)
			}

			// finally restore os.Args (global variable)
			t.Cleanup(func() {
				os.Args = oldArgs
			})
		})
	}
}
