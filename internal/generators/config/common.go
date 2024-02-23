package config

// type Code struct {
// 	Key     string   `yaml:"key"`
// 	Tmpl    string   `yaml:"tmpl,omitempty"`
// 	Imports []string `yaml:"imports,omitempty"`
// }

type DefaultConfig map[string]string

type File struct {
	Name    string   `yaml:"name"`
	Tmpl    string   `yaml:"tmpl,omitempty"`
	Imports []string `yaml:"imports,omitempty"`
}
