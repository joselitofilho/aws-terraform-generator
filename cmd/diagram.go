package cmd

import (
	"github.com/spf13/cobra"

	"github.com/joselitofilho/aws-terraform-generator/internal/fmtcolor"
	"github.com/joselitofilho/aws-terraform-generator/internal/generators/diagram"
)

// diagramCmd represents the structure command.
var diagramCmd = &cobra.Command{
	Use:   "diagram",
	Short: "Manage Diagram",
	Run: func(cmd *cobra.Command, _ []string) {
		diagramFilename, _ := cmd.Flags().GetString(flagDiagram)
		configFile, _ := cmd.Flags().GetString(flagConfig)
		output, _ := cmd.Flags().GetString(flagOutput)

		if err := diagram.NewDiagram(diagramFilename, configFile, output).Build(); err != nil {
			printErrorAndExit(err)
		}

		fmtcolor.White.Printf("Configuration file '%s' has been generated successfully\n", output)
	},
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
