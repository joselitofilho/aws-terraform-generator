package config

type Folder struct {
	Name  string `yaml:"name"`
	Files []File `yaml:"files"`
}

type Stack struct {
	Name    string   `yaml:"name"`
	Folders []Folder `yaml:"folders"`
}

type Structure struct {
	Stacks           []Stack         `yaml:"stacks"`
	DefaultTemplates []DefaultConfig `yaml:"default_templates,omitempty"`
}
