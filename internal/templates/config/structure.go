package config

type File struct {
	Name string `yaml:"name"`
	Tmpl string `yaml:"tmpl,omitempty"`
}

type Folder struct {
	Name  string `yaml:"name"`
	Files []File `yaml:"files"`
}

type Stack struct {
	StackName string   `yaml:"stack_name"`
	Folders   []Folder `yaml:"folders"`
}

type Structure struct {
	Stacks           []Stack             `yaml:"stacks"`
	DefaultTemplates []map[string]string `yaml:"default_templates"`
}
