package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type YAML struct {
	fileName string
}

func NewYAML(fileName string) *YAML {
	return &YAML{fileName: fileName}
}

func (y *YAML) Parse() (*Config, error) {
	yamlFile, err := os.ReadFile(y.fileName)
	if err != nil {
		return nil, fmt.Errorf("read YAML file error: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(yamlFile, &config); err != nil {
		return nil, fmt.Errorf("unmarshal YAML file error: %w", err)
	}

	return &config, nil
}
