package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/joselitofilho/aws-terraform-generator/internal/drawio"
	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
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
			printErrorAndExit(err)
		}

		configFile, err := cmd.Flags().GetString(diagramCMDFlagConfig)
		if err != nil {
			printErrorAndExit(err)
		}

		output, err := cmd.Flags().GetString(diagramCMDFlagOutput)
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

	mxFile, err := drawio.Parse(diagram)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	resources, err := drawio.ParseResources(mxFile)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	yamlConfigOut, err := transformers.TransformDrawIOToYAML(yamlConfig, resources)
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

	diagramCmd.Flags().StringP(diagramCMDFlagDiagram, "d", "", "Path to the XML file. For example: ./diagram.xml")
	diagramCmd.Flags().StringP(diagramCMDFlagConfig, "c", "", "Path to the YAML config file. For example: ./diagram.config.yaml")
	diagramCmd.Flags().StringP(diagramCMDFlagOutput, "o", "", "Path to the output file. For example: ./diagram.yaml")

	diagramCmd.MarkFlagRequired(diagramCMDFlagDiagram)
	diagramCmd.MarkFlagRequired(diagramCMDFlagConfig)
	diagramCmd.MarkFlagRequired(diagramCMDFlagOutput)
}
