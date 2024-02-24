package cmd

import (
	"github.com/spf13/cobra"

	"github.com/joselitofilho/aws-terraform-generator/internal/generators/structure"
)

const (
	structureCMDFlagConfig = "config"
	structureCMDFlagOutput = "output"
)

// structureCmd represents the structure command
var structureCmd = &cobra.Command{
	Use:   "structure",
	Short: "Manage Structure",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := cmd.Flags().GetString(structureCMDFlagConfig)
		if err != nil {
			printErrorAndExit(err)
		}

		output, err := cmd.Flags().GetString(structureCMDFlagOutput)
		if err != nil {
			printErrorAndExit(err)
		}

		err = structure.NewStructure(config, output).Build()
		if err != nil {
			printErrorAndExit(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(structureCmd)

	structureCmd.Flags().StringP(structureCMDFlagConfig, "c", "", "Path to the configuration file. For example: ./structure.config.yaml")
	structureCmd.Flags().StringP(structureCMDFlagOutput, "o", "", "Path to the output folder. For example: ./output")

	structureCmd.MarkFlagRequired(structureCMDFlagConfig)
	structureCmd.MarkFlagRequired(structureCMDFlagOutput)
}
