package devcontainer

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func Parse(path string) (*DevcontainerConfig, error) {
	rawContent, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config DevcontainerConfig
	if err := json.Unmarshal(rawContent, &config); err != nil {
		return nil, err
	}
	config.DockerComposeFile = resolveComposeFilePaths(path, config.DockerComposeFile)
	return &config, nil
}

func resolveComposeFilePaths(devcontainerPath string, composeFiles []string) []string {
	dirname := filepath.Dir(devcontainerPath)
	res := make([]string, len(composeFiles))
	for i, composeFile := range composeFiles {
		if filepath.IsAbs(composeFile) {
			res[i] = composeFile
			continue
		}
		res[i] = filepath.Join(dirname, composeFile)
	}
	return res
}
