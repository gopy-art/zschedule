package configs

import (
	"os"

	"gopkg.in/yaml.v2"
)

func ReadConfigFile(filepath string) ([]CommandLineConfig, error) {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var configs []CommandLineConfig
	if err := yaml.Unmarshal(content, &configs); err != nil {
		return nil, err
	}
	return configs, nil
}