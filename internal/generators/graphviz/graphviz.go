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

	applyStyleForNodes(resc, g, resourceImageMap, nodes, style)

	applyStyleForArrows(resc, egdes, g, nodes, style)

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

func applyStyleForNodes(
	resc *resources.ResourceCollection, g *dot.Graph, resourceImageMap map[resources.ResourceType]string,
	nodes map[string]dot.Node, style Style) {
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
}

func applyStyleForArrows(
	resc *resources.ResourceCollection, edges map[string]struct{}, g *dot.Graph, nodes map[string]dot.Node, style Style,
) {
	for _, rel := range resc.Relationships {
		if rel.Source == nil || rel.Target == nil {
			continue
		}

		edgeKey := rel.Source.Value() + "###" + rel.Target.Value()
		if _, ok := edges[edgeKey]; ok {
			continue
		}

		sourceNode := nodes[rel.Source.Value()]
		targetNode := nodes[rel.Target.Value()]

		if colors, exists := style.Arrows[rel.Source.Value()]; exists {
			if color, ok := colors[rel.Target.Value()]; ok {
				g.Edge(sourceNode, targetNode).Attr("color", color)
			} else {
				g.Edge(sourceNode, targetNode)
			}
		} else {
			g.Edge(sourceNode, targetNode)
		}

		edges[edgeKey] = struct{}{}
	}
}
