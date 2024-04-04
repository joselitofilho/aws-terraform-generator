package config

type Folder struct {
	Name    string   `yaml:"name"`
	Files   []File   `yaml:"files"`
	Folders []Folder `yaml:"folders"`
}

type Stack struct {
	Name    string   `yaml:"name"`
	Files   []File   `yaml:"files"`
	Folders []Folder `yaml:"folders"`
}

type Structure struct {
	Stacks           []Stack               `yaml:"stacks"`
	DefaultTemplates []FilenameTemplateMap `yaml:"default_templates,omitempty"`
}
