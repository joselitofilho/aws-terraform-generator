package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/joselitofilho/aws-terraform-generator/internal/drawio"
	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
	"github.com/joselitofilho/aws-terraform-generator/internal/transformers/drawiotoresources"
	"github.com/joselitofilho/aws-terraform-generator/internal/transformers/resourcestoyaml"
)

// diagramCmd represents the structure command.
var diagramCmd = &cobra.Command{
	Use:   "diagram",
	Short: "Manage Diagram",
	Run: func(cmd *cobra.Command, _ []string) {
		diagram, err := cmd.Flags().GetString(flagDiagram)
		if err != nil {
			printErrorAndExit(err)
		}

		configFile, err := cmd.Flags().GetString(flagConfig)
		if err != nil {
			printErrorAndExit(err)
		}

		output, err := cmd.Flags().GetString(flagOutput)
		if err != nil {
			printErrorAndExit(err)
		}

		if err := build(diagram, configFile, output); err != nil {
			printErrorAndExit(err)
		}

		fmt.Printf("Configuration file '%s' has been generated successfully\n", output)
	},
}

func build(diagram, configFile, output string) error {
	yamlConfig, err := config.NewYAML(configFile).Parse()
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	mxFile, err := drawio.ParseXML(diagram)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	resources, err := drawiotoresources.Transform(mxFile)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	yamlConfigOut, err := resourcestoyaml.NewTransformer(yamlConfig, resources).Transform()
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	data, err := yaml.Marshal(yamlConfigOut)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	err = os.WriteFile(output, data, os.ModePerm)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(diagramCmd)

	diagramCmd.Flags().StringP(flagDiagram, "d", "", "Path to the XML file. For example: ./diagram.xml")
	diagramCmd.Flags().StringP(flagConfig, "c", "",
		"Path to the YAML config file. For example: ./diagram.config.yaml")
	diagramCmd.Flags().StringP(flagOutput, "o", "", "Path to the output file. For example: ./diagram.yaml")

	_ = diagramCmd.MarkFlagRequired(flagDiagram)
	_ = diagramCmd.MarkFlagRequired(flagConfig)
	_ = diagramCmd.MarkFlagRequired(flagOutput)
}
