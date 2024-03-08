package graphviz

import (
	"github.com/emicklei/dot"

	"github.com/joselitofilho/aws-terraform-generator/internal/drawio"
)

const OrientationLeftRight = "LR"

type Config struct {
	Orientation string
}

func Build(
	resources *drawio.ResourceCollection, resourceImageMap map[drawio.ResourceType]string, config Config,
) (string, error) {
	g := dot.NewGraph(dot.Directed)

	if config.Orientation != "" {
		g.Attr("rankdir", config.Orientation)
	}

	g.NodeInitializer(func(n dot.Node) {
		n.Attrs("shape", "plaintext", "imagepos", "tc", "labelloc", "b", "height", "0.9")
	})

	g.EdgeInitializer(func(e dot.Edge) {
		e.Attrs("arrowhead", "vee", "arrowtail", "normal")
	})

	nodes := map[string]dot.Node{}

	for _, res := range resources.Resources {
		nodes[res.ID()] = g.Node(res.Value()).
			Attr("image", resourceImageMap[res.ResourceType()])
	}

	for _, rel := range resources.Relationships {
		if rel.Source == nil || rel.Target == nil {
			continue
		}

		if _, ok := nodes[rel.Source.ID()]; !ok {
			continue
		}

		if _, ok := nodes[rel.Target.ID()]; !ok {
			continue
		}

		g.Edge(nodes[rel.Source.ID()], nodes[rel.Target.ID()])
	}

	return g.String(), nil
}
