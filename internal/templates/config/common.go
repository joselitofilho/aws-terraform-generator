package config

type Code struct {
	Key     string   `yaml:"key"`
	Tmpl    string   `yaml:"tmpl,omitempty"`
	Imports []string `yaml:"imports,omitempty"`
}
