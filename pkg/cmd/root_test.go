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
	envVars  map[string]string
	expected *rootCommandOptions
}

func providerTestParseOptionsSuccess(t *testing.T) map[string]parseOptionsSuccessTestCase {
	fixturePaths := map[string]string{
		"single_normal":   "./test_data/devcontainer.single.normal.json",
		"multiple_normal": "./test_data/devcontainer.multiple.normal.json",
	}

	return map[string]parseOptionsSuccessTestCase{
		"[basic] no specified case": {
			[]string{},
			map[string]string{},
			&rootCommandOptions{
				composeOptions: composeOptions{
					ConfigPaths: []string{"./docker-compose.yml"},
				},
				baseService: "base_shell",
			},
		},
		"[basic; compose options] single config file specified case": {
			[]string{"-f", "./compose.yml"},
			map[string]string{},
			&rootCommandOptions{
				composeOptions: composeOptions{
					ConfigPaths: []string{"./compose.yml"},
				},
				baseService: "base_shell",
			},
		},
		"[basic; compose options] multiple config files specified case": {
			[]string{"-f", "./docker-compose.yml", "-f", "./docker-compose.override.yml"},
			map[string]string{},
			&rootCommandOptions{
				composeOptions: composeOptions{
					ConfigPaths: []string{"./docker-compose.yml", "./docker-compose.override.yml"},
				},
				baseService: "base_shell",
			},
		},
		"[basic; devcontainer options] single normal case": {
			[]string{"-c", fixturePaths["single_normal"]},
			map[string]string{
				"LOCAL_WORKSPACE_FOLDER": "/workspace",
			},
			&rootCommandOptions{
				composeOptions: composeOptions{
					ConfigPaths: []string{"test_data/docker-compose.yml"},
					ProjectDir:  "/workspace",
				},
				baseService: "base_shell2",
			},
		},
		"[basic; devcontainer options] multiple normal case": {
			[]string{"-c", fixturePaths["multiple_normal"]},
			map[string]string{
				"PROJECT_DIR": "/project/tmp",
			},
			&rootCommandOptions{
				composeOptions: composeOptions{
					ConfigPaths: []string{"src/docker-compose.yml", "test_data/compose.yaml"},
					ProjectDir:  "/project/tmp",
				},
				baseService: "shell",
			},
		},
		"[basic; base service option] base service specified case": {
			[]string{"-b", "base"},
			map[string]string{},
			&rootCommandOptions{
				composeOptions: composeOptions{
					ConfigPaths: []string{"./docker-compose.yml"},
				},
				baseService: "base",
			},
		},
		"[complicated] devcontainer config file & base service option specified case": {
			[]string{"-c", fixturePaths["multiple_normal"], "-b", "overwritten_base"},
			map[string]string{
				"PROJECT_DIR": "/project/tmp",
			},
			&rootCommandOptions{
				composeOptions: composeOptions{
					ConfigPaths: []string{"src/docker-compose.yml", "test_data/compose.yaml"},
					ProjectDir:  "/project/tmp",
				},
				baseService: "overwritten_base",
			},
		},
		"[complicated] devcontainer config file & compose config files option specified case": {
			[]string{"-c", fixturePaths["multiple_normal"], "-f", "./compose.yml", "-f", "./compose.override.yml"},
			map[string]string{
				"PROJECT_DIR": "/project/tmp",
			},
			&rootCommandOptions{
				composeOptions: composeOptions{
					ConfigPaths: []string{"./compose.yml", "./compose.override.yml"},
					ProjectDir:  "/project/tmp",
				},
				baseService: "shell",
			},
		},
		"[complicated] devcontainer config file & compose project directory option specified case": {
			[]string{"-c", fixturePaths["multiple_normal"], "--project-directory", "/workspace"},
			map[string]string{
				"PROJECT_DIR": "/project/tmp",
			},
			&rootCommandOptions{
				composeOptions: composeOptions{
					ConfigPaths: []string{"src/docker-compose.yml", "test_data/compose.yaml"},
					ProjectDir:  "/workspace",
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
			for envName, value := range c.envVars {
				os.Setenv(envName, value)
			}
			t.Cleanup(func() {
				for envName := range c.envVars {
					os.Unsetenv(envName)
				}
				os.Args = oldArgs
			})

			tc := func(opts *rootCommandOptions) {
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
			"parsed config path counts do not match [expected = \"%v\", actual = \"%v\"]",
			expected.ConfigPaths,
			actual.ConfigPaths,
		)
		t.FailNow()
	}
	for i, ep := range expected.ConfigPaths {
		ap := actual.ConfigPaths[i]
		if ep != ap {
			t.Errorf(
				"parsed config path does not match [expected = \"%v\", actual = \"%v\"]",
				expected.ConfigPaths,
				actual.ConfigPaths,
			)
			t.FailNow()
		}
	}
	if expected.ProjectDir != actual.ProjectDir {
		t.Errorf(
			"parsed project directory does not match [expected = \"%v\", actual = \"%v\"]",
			expected.ProjectDir,
			actual.ProjectDir,
		)
		t.FailNow()
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

func NewRootCommandMock(tc func(opts *rootCommandOptions)) *cobra.Command {
	do := createNewDevcontainerOptions()
	co := createNewComposeOptions()
	ro := createNewRootCommandOptions()
	cmd := &cobra.Command{
		RunE: func(cmd *cobra.Command, args []string) error {
			parsedDevcontainerOptions, err := do.parse()
			if err != nil {
				return err
			}
			opts, err := ro.createWithMerge(co, parsedDevcontainerOptions)
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
