package config

import "github.com/joselitofilho/aws-terraform-generator/internal/drawio"

type FilenameTemplateMap map[string]string

type File struct {
	Name    string   `yaml:"name"`
	Tmpl    string   `yaml:"tmpl,omitempty"`
	Imports []string `yaml:"imports,omitempty"`
}

type Filter struct {
	Match    []string `yaml:"match,omitempty"`
	NotMatch []string `yaml:"not_match,omitempty"`
}

type Filters map[drawio.ResourceType]Filter
