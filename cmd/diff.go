package cmd

import (
	"os"
	"path"

	"github.com/spf13/cobra"

	"github.com/diagram-code-generator/resources/pkg/parser/graphviz"
	"github.com/diagram-code-generator/resources/pkg/resources"

	"github.com/joselitofilho/aws-terraform-generator/internal/fmtcolor"
	"github.com/joselitofilho/aws-terraform-generator/internal/generators/draw"
	awsresources "github.com/joselitofilho/aws-terraform-generator/internal/resources"
	"github.com/joselitofilho/aws-terraform-generator/internal/transformers/yamltoresources"
)

// diffCmd represents the sqs command.
var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Manage Diff",
	Run: func(cmd *cobra.Command, _ []string) {
		left, err := cmd.Flags().GetString(flagLeft)
		if err != nil {
			printErrorAndExit(err)
		}

		right, err := cmd.Flags().GetString(flagRight)
		if err != nil {
			printErrorAndExit(err)
		}

		leftRc, err := yamltoresources.Parse(left)
		if err != nil {
			printErrorAndExit(err)
		}

		rightRc, err := yamltoresources.Parse(right)
		if err != nil {
			printErrorAndExit(err)
		}

		resources.PrintDiff(leftRc, rightRc, awsresources.AvailableTypes)

		addedResourcesByType, removedResourcesByType, addedRelationships, removedRelationships :=
			resources.FindDifferences(leftRc, rightRc)

		dotConfig := graphviz.Config{}
		style := graphviz.Style{Nodes: map[resources.Resource]string{}, Arrows: map[string][]map[string]string{}}

		for _, rscs := range addedResourcesByType {
			for i := range rscs {
				style.Nodes[rscs[i]] = "green"
			}
		}

		for _, rscs := range removedResourcesByType {
			for i := range rscs {
				style.Nodes[rscs[i]] = "red"
			}
		}

		for i := range addedRelationships {
			arrowTarget := style.Arrows[addedRelationships[i].Source.Value()]
			arrowTarget = append(arrowTarget, map[string]string{addedRelationships[i].Target.Value(): "green"})

			style.Arrows[addedRelationships[i].Source.Value()] = arrowTarget
		}

		for i := range removedRelationships {
			arrowTarget := style.Arrows[removedRelationships[i].Source.Value()]
			arrowTarget = append(arrowTarget, map[string]string{removedRelationships[i].Target.Value(): "red"})

			style.Arrows[removedRelationships[i].Source.Value()] = arrowTarget
		}

		dotContent := graphviz.BuildWithStyle(leftRc, draw.DefaultResourceImageMap.ToStringMap(), dotConfig, style)

		dotFilename := "diff.dot"

		dotfile, err := os.Create(path.Join(".", dotFilename))
		if err != nil {
			printErrorAndExit(err)
		}
		defer dotfile.Close()

		if _, err := dotfile.WriteString(dotContent); err != nil {
			printErrorAndExit(err)
		}

		fmtcolor.White.Println("The graphviz dot file has been generated successfully.")
	},
}

func init() {
	rootCmd.AddCommand(diffCmd)

	diffCmd.Flags().StringP(flagLeft, "l", "", "Path to the left YAML config file. For example: ./left.yaml")
	diffCmd.Flags().StringP(flagRight, "r", "", "Path to the right YAML config file. For example: ./right.yaml")

	_ = diffCmd.MarkFlagRequired(flagLeft)
	_ = diffCmd.MarkFlagRequired(flagRight)
}
