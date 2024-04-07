package domain

import (
	"github.com/docker/compose/v2/cmd/compose"
	"testing"
)

type unmarshalTestCase struct {
	opts     *compose.ProjectOptions
	expected []string
}

func providerTestUnmarshal(t *testing.T) map[string]unmarshalTestCase {
	return map[string]unmarshalTestCase{
		"[basic] single ConfigPath specified case": {
			opts: &compose.ProjectOptions{
				ConfigPaths: []string{"./docker-compose.yml"},
			},
			expected: []string{"docker", "compose", "--file", "./docker-compose.yml"},
		},
		"[basic] multiple ConfigPaths specified case": {
			opts: &compose.ProjectOptions{
				ConfigPaths: []string{"./docker-compose.yml", "./docker-compose.override.yml"},
			},
			expected: []string{"docker", "compose", "--file", "./docker-compose.yml", "--file", "./docker-compose.override.yml"},
		},
		"[basic] multiple options specified case": {
			opts: &compose.ProjectOptions{
				ConfigPaths: []string{"test_data/docker-compose.yml"},
				ProjectDir:  "/workspace",
			},
			expected: []string{"docker", "compose", "--file", "test_data/docker-compose.yml", "--project-directory", "/workspace"},
		},
	}
}

func TestUnmarshal(t *testing.T) {
	cases := providerTestUnmarshal(t)
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			actual := Unmarshal(c.opts)
			if el, al := len(c.expected), len(actual); el != al {
				t.Errorf(
					"extracted options counts do not match [expected = \"%v\", actual = \"%v\"]",
					c.expected,
					actual,
				)
				t.FailNow()
			}
			for i, eo := range c.expected {
				ao := actual[i]
				if eo != ao {
					t.Errorf(
						"extracted options do not match [expected = \"%v\", actual = \"%v\"]",
						c.expected,
						actual,
					)
					t.FailNow()
				}
			}
		})
	}
}
