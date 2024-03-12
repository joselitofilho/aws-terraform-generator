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
	return BuildWithStyle(resc, resourceImageMap, config, Style{})
}

func BuildWithStyle(
	resc *resources.ResourceCollection, resourceImageMap map[resources.ResourceType]string, config Config, style Style,
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
	egdes := map[string]struct{}{}

	for i := range resc.Resources {
		res := resc.Resources[i]

		node := g.Node(res.Value()).Attr("image", resourceImageMap[res.ResourceType()])

		if color, ok := style.Nodes[res]; ok {
			node = node.Attr("fontcolor", color)
		}

		nodes[res.Value()] = node
	}

	for k, v := range style.Nodes {
		nodes[k.Value()] = g.Node(k.Value()).Attr("fontcolor", v).Attr("image", resourceImageMap[k.ResourceType()])
	}

	for i := range resc.Relationships {
		rel := resc.Relationships[i]

		if rel.Source == nil || rel.Target == nil {
			continue
		}

		edgeKey := rel.Source.Value() + "###" + rel.Target.Value()

		if style.Arrows[rel.Source.Value()] != nil {
			if color, ok := style.Arrows[rel.Source.Value()][rel.Target.Value()]; ok {
				if _, ok := egdes[edgeKey]; !ok {
					g.Edge(nodes[rel.Source.Value()], nodes[rel.Target.Value()]).Attr("color", color)
					egdes[edgeKey] = struct{}{}
				}
			} else {
				if _, ok := egdes[edgeKey]; !ok {
					g.Edge(nodes[rel.Source.Value()], nodes[rel.Target.Value()])
					egdes[edgeKey] = struct{}{}
				}
			}
		} else {
			if _, ok := egdes[edgeKey]; !ok {
				g.Edge(nodes[rel.Source.Value()], nodes[rel.Target.Value()])
				egdes[edgeKey] = struct{}{}
			}
		}
	}

	for source, v := range style.Arrows {
		for target, color := range v {
			edgeKey := source + "###" + target

			if _, ok := egdes[edgeKey]; !ok {
				g.Edge(nodes[source], nodes[target]).Attr("color", color)
				egdes[edgeKey] = struct{}{}
			}
		}
	}

	return g.String()
}
