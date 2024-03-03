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
			oldArgs := os.Args
			t.Cleanup(func() {
				os.Args = oldArgs
			})

			tc := func (opts *rootCommandOptions) {
				if el, al := len(c.expected), len(opts.ConfigPaths); el != al {
					t.Errorf(
						"parsed config path counts do not match [expected = \"%s\", actual = \"%s\"]",
						c.expected,
						opts.ConfigPaths,
					)
					t.FailNow()
				}
				for i, ep := range c.expected {
					ap := opts.ConfigPaths[i]
					if ep != ap {
						t.Errorf(
							"parsed path does not match [expected = \"%s\", actual = \"%s\"]",
							ep,
							ap,
						)
						t.FailNow()
					}
				}
			}

			os.Args = append([]string{"command"}, c.args...)
			if err := NewRootCommandMock(tc).Execute(); err != nil {
				t.Errorf(
					"command execute error: %v",
					err.Error(),
				)
				t.FailNow()
			}
		})
	}
}

type SetDevcontainerOptionsSuccessTestCase struct {
	args     []string
	expected *devcontainerOptions
}

// func providerTestSetDevcontainerOptionsSuccess(t *testing.T) map[string]SetDevcontainerOptionsSuccessTestCase {
// 	fixturePaths := map[string]string{
// 		"single_normal": "./test_data/devcontainer.single.normal.json",
// 		"multiple_normal": "./test_data/devcontainer.multiple.normal.json",
// 	}

// 	return map[string]SetDevcontainerOptionsSuccessTestCase{
// 		"single normal case": {
// 			[]string{"-c", fixturePaths["single_normal"]},
// 			&devcontainerOptions{
// 				path: fixturePaths["single_normal"],
// 				dockerComposeFile: []string{"./docker-compose.yml"},
// 				service: "base_shell",
// 			},
// 		},
// 		"multiple normal case": {
// 			[]string{"-c", fixturePaths["multiple_normal"]},
// 			&devcontainerOptions{
// 				path: fixturePaths["multiple_normal"],
// 				dockerComposeFile: []string{"../src/docker-compose.yml", "./compose.yaml"},
// 				service: "shell",
// 			},
// 		},
// 		"not specified case": {
// 			[]string{},
// 			nil,
// 		},
// 	}
// }

// func TestSetDevcontainerOptionsSuccess(t *testing.T) {
// 	cases := providerTestSetDevcontainerOptionsSuccess(t)
// 	for name, c := range cases {
// 		t.Run(name, func(t *testing.T) {
// 			// no parallelization because using os.Args (global variable)
// 			oldArgs := os.Args
// 			t.Cleanup(func() {
// 				os.Args = oldArgs
// 			})

// 			os.Args = append([]string{"command"}, c.args...)
// 			cmd, opts := NewRootCommandMock()
// 			if err := cmd.Execute(); err != nil {
// 				t.Errorf(
// 					"command execute error: %v",
// 					err.Error(),
// 				)
// 				t.FailNow()
// 			}

// 			if c.expected == nil {
// 				if actual == nil {
// 					return
// 				} else {
// 					t.Errorf(
// 						"nil devcontainerOpts expected but actual is not nil [actual = %v]",
// 						*actual,
// 					)
// 					t.FailNow()
// 				}
// 			}
// 			if actual == nil {
// 				t.Errorf(
// 					"not nil devcontainerOpts expected but actual devcontainerOpts is nil [expected = %v]",
// 					*c.expected,
// 				)
// 				t.FailNow()
// 			}

// 			if c.expected.path != actual.path {
// 				t.Errorf(
// 					"parsed devcontainer path does not match [expected = \"%s\", actual = \"%s\"]",
// 					c.expected.path,
// 					actual.path,
// 				)
// 				t.FailNow()
// 			}
// 			for i, expectedFile := range c.expected.dockerComposeFile {
// 				actualFile := actual.dockerComposeFile[i]
// 				if expectedFile != actualFile {
// 					t.Errorf(
// 						"parsed dockerComposeFile does not match [expected = \"%v\", actual = \"%v\"]",
// 						c.expected.dockerComposeFile,
// 						actual.dockerComposeFile,
// 					)
// 					t.FailNow()
// 				}
// 			}
// 			if c.expected.service != actual.service {
// 				t.Errorf(
// 					"parsed baseService path does not match [expected = \"%s\", actual = \"%s\"]",
// 					c.expected.path,
// 					actual.path,
// 				)
// 				t.FailNow()
// 			}
// 		})
// 	}
// }

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
			t.Cleanup(func() {
				os.Args = oldArgs
			})

			tc := func (opts *rootCommandOptions) {
				if c.expected != opts.baseService {
					t.Errorf(
						"parsed base service do not match [expected = \"%s\", actual = \"%s\"]",
						c.expected,
						opts.baseService,
					)
					t.FailNow()
				}
			}

			os.Args = append([]string{"command"}, c.args...)
			if err := NewRootCommandMock(tc).Execute(); err != nil {
				t.Errorf(
					"command execute error: %v",
					err.Error(),
				)
				t.FailNow()
			}
		})
	}
}

func NewRootCommandMock(tc func (opts *rootCommandOptions)) *cobra.Command {
	ro, co, do := &rootCommandOptions{}, &composeOptions{}, &devcontainerOptions{}
	cmd := &cobra.Command{
		RunE: func (cmd *cobra.Command, args []string) error {
			parsedDevcontainerOptions, err := do.parse()
			if err != nil {
				return err
			}
			opts, err := ro.createWithOverWrite(co, parsedDevcontainerOptions)
			if err != nil {
				return err
			}
			tc(opts)
			return nil
		},
	}
	f := cmd.PersistentFlags()
	ro.setBaseServiceOption(f)
	co.set(f)
	do.set(f)
	return cmd
}
