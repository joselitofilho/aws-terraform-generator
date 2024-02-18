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
		input, err := cmd.Flags().GetString("input")
		if err != nil {
			panic(err)
		}

		output, err := cmd.Flags().GetString("output")
		if err != nil {
			panic(err)
		}

		err = structure.NewStructure(input, output).Build()
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(structureCmd)

	structureCmd.Flags().StringP("input", "i", "", "Path to the yaml file. For example: lambdas.yaml")
	structureCmd.Flags().StringP("output", "o", "", "Path to the output folder. For example: ./output")

	structureCmd.MarkFlagRequired("input")
	structureCmd.MarkFlagRequired("output")
}
