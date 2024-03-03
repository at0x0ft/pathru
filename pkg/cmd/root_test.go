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

type parseOptionsSuccessTestCase struct {
	args     []string
	expected *rootCommandOptions
}

func providerTestParseOptionsSuccess(t *testing.T) map[string]parseOptionsSuccessTestCase {
	fixturePaths := map[string]string{
		"single_normal": "./test_data/devcontainer.single.normal.json",
		"multiple_normal": "./test_data/devcontainer.multiple.normal.json",
	}

	return map[string]parseOptionsSuccessTestCase{
		"[basic] no specified case": {
			[]string{},
			&rootCommandOptions{
				composeOptions: composeOptions{
					ConfigPaths: []string{"./docker-compose.yml"},
				},
				baseService: "base_shell",
			},
		},
		"[basic; compose options] single config file specified case": {
			[]string{"-f", "./compose.yml"},
			&rootCommandOptions{
				composeOptions: composeOptions{
					ConfigPaths: []string{"./compose.yml"},
				},
				baseService: "base_shell",
			},
		},
		"[basic; compose options] multiple config files specified case": {
			[]string{"-f", "./docker-compose.yml", "-f", "./docker-compose.override.yml"},
			&rootCommandOptions{
				composeOptions: composeOptions{
					ConfigPaths: []string{"./docker-compose.yml", "./docker-compose.override.yml"},
				},
				baseService: "base_shell",
			},
		},
		"[basic; devcontainer options] single normal case": {
			[]string{"-c", fixturePaths["single_normal"]},
			&rootCommandOptions{
				composeOptions: composeOptions{
					ConfigPaths: []string{"./docker-compose.yml"},
				},
				baseService: "base_shell2",
			},
		},
		"[basic; devcontainer options] multiple normal case": {
			[]string{"-c", fixturePaths["multiple_normal"]},
			&rootCommandOptions{
				composeOptions: composeOptions{
					ConfigPaths: []string{"../src/docker-compose.yml", "./compose.yaml"},
				},
				baseService: "shell",
			},
		},
		"[basic; base service option] base service specified case": {
			[]string{"-b", "base"},
			&rootCommandOptions{
				composeOptions: composeOptions{
					ConfigPaths: []string{"./docker-compose.yml"},
				},
				baseService: "base",
			},
		},
		"[complicated] devcontainer config file & base service option specified case": {
			[]string{"-c", fixturePaths["multiple_normal"], "-b", "overwritten_base"},
			&rootCommandOptions{
				composeOptions: composeOptions{
					ConfigPaths: []string{"../src/docker-compose.yml", "./compose.yaml"},
				},
				baseService: "overwritten_base",
			},
		},
		"[complicated] devcontainer config file & compose config files option specified case": {
			[]string{"-c", fixturePaths["multiple_normal"], "-f", "./compose.yml", "-f", "./compose.override.yml"},
			&rootCommandOptions{
				composeOptions: composeOptions{
					ConfigPaths: []string{"./compose.yml", "./compose.override.yml"},
				},
				baseService: "shell",
			},
		},
	}
}

func TestParseOptionsSuccess(t *testing.T) {
	cases := providerTestParseOptionsSuccess(t)
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			// no parallelization because using os.Args (global variable)
			oldArgs := os.Args
			t.Cleanup(func() {
				os.Args = oldArgs
			})

			tc := func (opts *rootCommandOptions) {
				assertComposeOptions(t, c.expected, opts)
				assertBaseServiceOptions(t, c.expected, opts)
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

func assertComposeOptions(
	t *testing.T,
	expected *rootCommandOptions,
	actual *rootCommandOptions,
) {
	if el, al := len(expected.ConfigPaths), len(actual.ConfigPaths); el != al {
		t.Errorf(
			"parsed config path counts do not match [expected = \"%s\", actual = \"%s\"]",
			expected.ConfigPaths,
			actual.ConfigPaths,
		)
		t.FailNow()
	}
	for i, ep := range expected.ConfigPaths {
		ap := actual.ConfigPaths[i]
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

func assertBaseServiceOptions(
	t *testing.T,
	expected *rootCommandOptions,
	actual *rootCommandOptions,
) {
	if expected.baseService != actual.baseService {
		t.Errorf(
			"parsed base service do not match [expected = \"%s\", actual = \"%s\"]",
			expected.baseService,
			actual.baseService,
		)
		t.FailNow()
	}
}

func NewRootCommandMock(tc func (opts *rootCommandOptions)) *cobra.Command {
	do := createNewDevcontainerOptions()
	co := createNewComposeOptions()
	ro := createNewRootCommandOptions()
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
	co.set(f)
	do.set(f)
	ro.setBaseServiceOption(f)
	return cmd
}
