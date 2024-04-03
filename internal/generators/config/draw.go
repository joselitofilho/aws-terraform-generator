package config

import awsresources "github.com/joselitofilho/aws-terraform-generator/internal/resources"

type Images map[awsresources.ResourceType]string
type ReplaceableTexts map[string]string

type Draw struct {
	Name             string           `yaml:"name,omitempty"`
	Orientation      string           `yaml:"orientation,omitempty"`
	ReplaceableTexts ReplaceableTexts `yaml:"replaceable_texts,omitempty"`
	Images           Images           `yaml:"images,omitempty"`
	Filters          Filters          `yaml:"filters,omitempty"`
}
