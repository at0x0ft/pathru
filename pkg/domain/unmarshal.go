package domain

import (
	"github.com/docker/compose/v2/cmd/compose"
)

func Unmarshal(opts *compose.ProjectOptions) []string {
	extractedOpts := extractComposeProjectOptions(opts)
	return append([]string{"docker", "compose"}, extractedOpts...)
}

// ref: https://github.com/docker/compose/blob/d10a179f3e451f8b03fd99271f011c34bc31bedb/cmd/compose/compose.go#L157-L167
func extractComposeProjectOptions(opts *compose.ProjectOptions) []string {
	var result []string

	for _, profile := range opts.Profiles {
		result = append(result, "--profile", profile)
	}

	if opts.ProjectName != "" {
		result = append(result, "--project-name", opts.ProjectName)
	}

	for _, configPath := range opts.ConfigPaths {
		result = append(result, "--file", configPath)
	}

	if opts.EnvFiles != nil {
		for _, envFile := range opts.EnvFiles {
			result = append(result, "--env-file", envFile)
		}
	}

	if opts.ProjectDir != "" {
		result = append(result, "--project-directory", opts.ProjectDir)
	}

	return result
}
