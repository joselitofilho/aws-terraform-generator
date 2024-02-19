package cmd

import (
	"fmt"
	"os"

	"github.com/joselitofilho/aws-terraform-generator/internal/drawio"
	"github.com/joselitofilho/aws-terraform-generator/internal/transformers"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// diagramCmd represents the structure command
var diagramCmd = &cobra.Command{
	Use:   "diagram",
	Short: "Manage Diagram",
	Run: func(cmd *cobra.Command, args []string) {
		input, err := cmd.Flags().GetString("input")
		if err != nil {
			panic(err)
		}

		output, err := cmd.Flags().GetString("output")
		if err != nil {
			panic(err)
		}

		root, err := drawio.Parse(input)
		if err != nil {
			panic(err)
		}

		resources, err := drawio.ParseResources(root)
		if err != nil {
			panic(err)
		}

		yamlConfig, err := transformers.TransformDrawIOToYAML("location", resources)
		if err != nil {
			panic(err)
		}

		data, err := yaml.Marshal(yamlConfig)
		if err != nil {
			panic(err)
		}

		err = os.WriteFile(output, data, os.ModePerm)
		if err != nil {
			panic(err)
		}

		fmt.Printf("YAML '%s' has been generated successfully\n", output)
	},
}

func init() {
	rootCmd.AddCommand(diagramCmd)

	diagramCmd.Flags().StringP("input", "i", "", "Path to the yaml file. For example: diagram.xml")
	diagramCmd.Flags().StringP("output", "o", "", "Path to the output file. For example: ./diagram.yaml")

	diagramCmd.MarkFlagRequired("input")
	diagramCmd.MarkFlagRequired("output")
}