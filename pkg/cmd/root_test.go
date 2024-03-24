package cmd

import (
	"github.com/docker/compose/v2/cmd/compose"
	"github.com/spf13/cobra"
	"os"
	"testing"
)

type parseOptionsSuccessTestCase struct {
	args     []string
	envVars  map[string]string
	expected *parseOptionsSuccessExpectedValues
}

type parseOptionsSuccessExpectedValues struct {
	prjOpts     *compose.ProjectOptions
	baseService string
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
			&parseOptionsSuccessExpectedValues{
				prjOpts: &compose.ProjectOptions{
					ConfigPaths: []string{"./docker-compose.yml"},
				},
			},
		},
		"[basic; compose options] single config file specified case": {
			[]string{"-f", "./compose.yml"},
			map[string]string{},
			&parseOptionsSuccessExpectedValues{
				prjOpts: &compose.ProjectOptions{
					ConfigPaths: []string{"./compose.yml"},
				},
			},
		},
		"[basic; compose options] multiple config files specified case": {
			[]string{"-f", "./docker-compose.yml", "-f", "./docker-compose.override.yml"},
			map[string]string{},
			&parseOptionsSuccessExpectedValues{
				prjOpts: &compose.ProjectOptions{
					ConfigPaths: []string{"./docker-compose.yml", "./docker-compose.override.yml"},
				},
			},
		},
		"[basic; devcontainer options] single normal case": {
			[]string{"-c", fixturePaths["single_normal"]},
			map[string]string{
				"LOCAL_WORKSPACE_FOLDER": "/workspace",
			},
			&parseOptionsSuccessExpectedValues{
				prjOpts: &compose.ProjectOptions{
					ConfigPaths: []string{"test_data/docker-compose.yml"},
					ProjectDir:  "/workspace",
				},
				baseService: "base_shell",
			},
		},
		"[basic; devcontainer options] multiple normal case": {
			[]string{"-c", fixturePaths["multiple_normal"]},
			map[string]string{
				"PROJECT_DIR": "/project/tmp",
			},
			&parseOptionsSuccessExpectedValues{
				prjOpts: &compose.ProjectOptions{
					ConfigPaths: []string{"src/docker-compose.yml", "test_data/compose.yaml"},
					ProjectDir:  "/project/tmp",
				},
				baseService: "shell",
			},
		},
		"[basic; base service option] base service specified case": {
			[]string{"-b", "base"},
			map[string]string{},
			&parseOptionsSuccessExpectedValues{
				prjOpts: &compose.ProjectOptions{
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
			&parseOptionsSuccessExpectedValues{
				prjOpts: &compose.ProjectOptions{
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
			&parseOptionsSuccessExpectedValues{
				prjOpts: &compose.ProjectOptions{
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
			&parseOptionsSuccessExpectedValues{
				prjOpts: &compose.ProjectOptions{
					ConfigPaths: []string{"src/docker-compose.yml", "test_data/compose.yaml"},
					ProjectDir:  "/workspace",
				},
				baseService: "shell",
			},
		},
		"[complicated] devcontainer config file & localhost base service option specified case": {
			[]string{"-c", fixturePaths["multiple_normal"], "-b", ""},
			map[string]string{
				"PROJECT_DIR": "/project/tmp",
			},
			&parseOptionsSuccessExpectedValues{
				prjOpts: &compose.ProjectOptions{
					ConfigPaths: []string{"src/docker-compose.yml", "test_data/compose.yaml"},
					ProjectDir:  "/project/tmp",
				},
				baseService: "",
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

			tc := func(prjOpts *compose.ProjectOptions, baseService string) {
				assertComposeOptions(t, c.expected.prjOpts, prjOpts)
				assertBaseServiceOptions(t, c.expected.baseService, baseService)
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

func assertComposeOptions(t *testing.T, expected, actual *compose.ProjectOptions) {
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

func assertBaseServiceOptions(t *testing.T, expected, actual string) {
	if expected != actual {
		t.Errorf(
			"parsed base service do not match [expected = \"%s\", actual = \"%s\"]",
			expected,
			actual,
		)
		t.FailNow()
	}
}

func NewRootCommandMock(
	tc func(prjOpts *compose.ProjectOptions, baseService string),
) *cobra.Command {
	opts := createRootOptions()
	cmd := &cobra.Command{
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.parse(); err != nil {
				return err
			}
			prjOpts, baseService := opts.getProjectOptions(), opts.getBaseService()
			tc(prjOpts, baseService)
			return nil
		},
	}
	opts.set(cmd.PersistentFlags())
	return cmd
}
