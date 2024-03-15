package graphviz

import "github.com/joselitofilho/aws-terraform-generator/internal/resources"

type Style struct {
	Nodes  map[resources.Resource]string
	Arrows map[string][]map[string]string
}
