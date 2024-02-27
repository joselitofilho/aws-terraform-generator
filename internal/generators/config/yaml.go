package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

var (
	osReadFile    = os.ReadFile
	yamlUnmarshal = yaml.Unmarshal
)

type YAML struct {
	fileName string
}

func NewYAML(fileName string) *YAML {
	return &YAML{fileName: fileName}
}

func (y *YAML) Parse() (*Config, error) {
	yamlFile, err := osReadFile(y.fileName)
	if err != nil {
		return nil, fmt.Errorf("read YAML file error: %w", err)
	}

	var config Config
	if err := yamlUnmarshal(yamlFile, &config); err != nil {
		return nil, fmt.Errorf("unmarshal YAML file error: %w", err)
	}

	return &config, nil
}
