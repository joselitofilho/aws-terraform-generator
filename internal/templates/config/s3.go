package config

type S3 struct {
	Name   string `yaml:"name"`
	Key    string `yaml:"key,omitempty"`
	Source string `yaml:"source,omitempty"`
}
