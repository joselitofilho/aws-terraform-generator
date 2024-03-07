package graphviz

import (
	"github.com/emicklei/dot"

	"github.com/joselitofilho/aws-terraform-generator/internal/drawio"
)

func Build(resources *drawio.ResourceCollection, resourceImageMap map[drawio.ResourceType]string) (string, error) {
	g := dot.NewGraph(dot.Directed)

	g.Attr("rankdir", "LR")

	g.NodeInitializer(func(n dot.Node) {
		n.Attrs("shape", "plaintext", "imagepos", "tc", "labelloc", "b", "height", "0.9")
	})

	nodes := map[string]dot.Node{}

	for _, res := range resources.Resources {
		nodes[res.ID()] = g.Node(res.Value()).
			Attr("image", resourceImageMap[res.ResourceType()])
	}

	for _, rel := range resources.Relationships {
		g.Edge(nodes[rel.Source.ID()], nodes[rel.Target.ID()])
	}

	return g.String(), nil
}
