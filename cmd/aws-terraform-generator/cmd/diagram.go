package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/joselitofilho/aws-terraform-generator/internal/drawio"
	"github.com/joselitofilho/aws-terraform-generator/internal/templates/config"
	"github.com/joselitofilho/aws-terraform-generator/internal/transformers"
)

const (
	diagramCMDFlagConfig  = "config"
	diagramCMDFlagDiagram = "diagram"
	diagramCMDFlagOutput  = "output"
)

// diagramCmd represents the structure command
var diagramCmd = &cobra.Command{
	Use:   "diagram",
	Short: "Manage Diagram",
	Run: func(cmd *cobra.Command, args []string) {
		diagram, err := cmd.Flags().GetString(diagramCMDFlagDiagram)
		if err != nil {
			panic(err)
		}

		configFile, err := cmd.Flags().GetString(diagramCMDFlagConfig)
		if err != nil {
			panic(err)
		}

		output, err := cmd.Flags().GetString(diagramCMDFlagOutput)
		if err != nil {
			panic(err)
		}

		root, err := drawio.Parse(diagram)
		if err != nil {
			panic(err)
		}

		resources, err := drawio.ParseResources(root)
		if err != nil {
			panic(err)
		}

		yamlParser := config.NewYAML(configFile)

		yamlConfig, err := yamlParser.Parse()
		if err != nil {
			panic(err)
		}

		yamlConfigOut, err := transformers.TransformDrawIOToYAML(yamlConfig, resources)
		if err != nil {
			panic(err)
		}

		data, err := yaml.Marshal(yamlConfigOut)
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

	diagramCmd.Flags().StringP(diagramCMDFlagDiagram, "d", "", "Path to the xml file. For example: ./diagram.xml")
	diagramCmd.Flags().StringP(diagramCMDFlagConfig, "c", "", "Path to the YAML config file. For example: ./diagram.config.yaml")
	diagramCmd.Flags().StringP(diagramCMDFlagOutput, "o", "", "Path to the output file. For example: ./diagram.yaml")

	diagramCmd.MarkFlagRequired(diagramCMDFlagDiagram)
	diagramCmd.MarkFlagRequired(diagramCMDFlagConfig)
	diagramCmd.MarkFlagRequired(diagramCMDFlagOutput)
}
