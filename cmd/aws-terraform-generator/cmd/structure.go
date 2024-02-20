package cmd

import (
	"github.com/spf13/cobra"

	"github.com/joselitofilho/aws-terraform-generator/internal/templates/structure"
)

// structureCmd represents the structure command
var structureCmd = &cobra.Command{
	Use:   "structure",
	Short: "Manage Structure",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := cmd.Flags().GetString("config")
		if err != nil {
			panic(err)
		}

		output, err := cmd.Flags().GetString("output")
		if err != nil {
			panic(err)
		}

		err = structure.NewStructure(config, output).Build()
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(structureCmd)

	structureCmd.Flags().StringP("config", "c", "", "Path to the configuration file. For example: ./structure.config.yaml")
	structureCmd.Flags().StringP("output", "o", "", "Path to the output folder. For example: ./output")

	structureCmd.MarkFlagRequired("config")
	structureCmd.MarkFlagRequired("output")
}
