package config

import (
	"github.com/diagram-code-generator/resources/pkg/parser/graphviz/dot"
	awsresources "github.com/joselitofilho/aws-terraform-generator/internal/resources"
)

type Images map[awsresources.ResourceType]string

func (m Images) ToStringMap() map[string]string {
	result := map[string]string{}

	for k, v := range m {
		result[k.String()] = v
	}

	return result
}

type ReplaceableTexts map[string]string

type Draw struct {
	Name             string               `yaml:"name,omitempty"`
	Direction        dot.DiagramDirection `yaml:"direction,omitempty"`
	Splines          dot.DiagramSpline    `yaml:"splines,omitempty"`
	ReplaceableTexts ReplaceableTexts     `yaml:"replaceable_texts,omitempty"`
	Images           Images               `yaml:"images,omitempty"`
	Filters          Filters              `yaml:"filters,omitempty"`
}
