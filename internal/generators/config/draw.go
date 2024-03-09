package config

import "github.com/joselitofilho/aws-terraform-generator/internal/resources"

type Images map[resources.ResourceType]string

type Draw struct {
	Orientation string  `yaml:"orientation,omitempty"`
	Images      Images  `yaml:"images,omitempty"`
	Filters     Filters `yaml:"filters,omitempty"`
}
