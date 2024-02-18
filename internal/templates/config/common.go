package config

type Code struct {
	Key     string   `yaml:"key"`
	Tmpl    string   `yaml:"tmpl"`
	Imports []string `yaml:"imports"`
}
