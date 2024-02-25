//nolint:dupl // That is a false positive
package cmd

import (
	"github.com/spf13/cobra"

	"github.com/joselitofilho/aws-terraform-generator/internal/generators/structure"
)

// structureCmd represents the structure command
var structureCmd = &cobra.Command{
	Use:   "structure",
	Short: "Manage Structure",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := cmd.Flags().GetString(flagConfig)
		if err != nil {
			printErrorAndExit(err)
		}

		output, err := cmd.Flags().GetString(flagOutput)
		if err != nil {
			printErrorAndExit(err)
		}

		err = structure.NewStructure(config, output).Build()
		if err != nil {
			printErrorAndExit(err)
		}
	},
}

//nolint:gochecknoinits // That is the way the cobra library operates
func init() {
	rootCmd.AddCommand(structureCmd)

	structureCmd.Flags().StringP(flagConfig, "c", "",
		"Path to the configuration file. For example: ./structure.config.yaml")
	structureCmd.Flags().StringP(flagOutput, "o", "",
		"Path to the output folder. For example: ./output")

	_ = structureCmd.MarkFlagRequired(flagConfig)
	_ = structureCmd.MarkFlagRequired(flagOutput)
}
