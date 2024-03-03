package devcontainer

import (
	"encoding/json"
	"os"
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
	return &config, nil
}
