package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Lambda struct {
	Name        string              `yaml:"name"`
	Description string              `yaml:"description"`
	Envars      []map[string]string `yaml:"envars"`
	SQSTriggers []SQSTrigger        `yaml:"sqs-triggers"`
	Cron        []Cron              `yaml:"crons"`
	Code        []Code              `yaml:"code"`
}

type APILambda struct {
	Name        string              `yaml:"name"`
	Description string              `yaml:"description"`
	Envars      []map[string]string `yaml:"envars"`
	Code        []Code              `yaml:"code"`
}

type SQSTrigger struct {
	SourceARN string `yaml:"source_arn"`
}

type Cron struct {
	ScheduleExpression string `yaml:"schedule_expression"`
	IsEnabled          string `yaml:"is_enabled"`
}

type Code struct {
	Key     string   `yaml:"key"`
	Tmpl    string   `yaml:"tmpl"`
	Imports []string `yaml:"imports"`
}

type Config struct {
	Lambdas []Lambda `yaml:"lambdas"`
}

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
