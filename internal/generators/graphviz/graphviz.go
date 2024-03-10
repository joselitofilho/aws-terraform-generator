package graphviz

import (
	"github.com/emicklei/dot"

	"github.com/joselitofilho/aws-terraform-generator/internal/resources"
)

const OrientationLeftRight = "LR"

type Config struct {
	Orientation string
}

func Build(
	resc *resources.ResourceCollection, resourceImageMap map[resources.ResourceType]string, config Config,
) string {
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

	for _, res := range resc.Resources {
		nodes[res.ID()] = g.Node(res.Value()).
			Attr("image", resourceImageMap[res.ResourceType()])
	}

	for _, rel := range resc.Relationships {
		if rel.Source == nil || rel.Target == nil {
			continue
		}

		g.Edge(nodes[rel.Source.ID()], nodes[rel.Target.ID()])
	}

	return g.String()
}
