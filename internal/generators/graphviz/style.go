package graphviz

import "github.com/diagram-code-generator/resources/pkg/resources"

type Style struct {
	Nodes  map[resources.Resource]string
	Arrows map[string][]map[string]string
}
