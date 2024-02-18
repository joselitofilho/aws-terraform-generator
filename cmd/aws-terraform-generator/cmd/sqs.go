package cmd

import (
	"github.com/joselitofilho/aws-terraform-generator/internal/templates/sqs"
	"github.com/spf13/cobra"
)

// sqsCmd represents the sqs command
var sqsCmd = &cobra.Command{
	Use:   "sqs",
	Short: "Manage SQS",
	Run: func(cmd *cobra.Command, args []string) {
		input, err := cmd.Flags().GetString("input")
		if err != nil {
			panic(err)
		}

		output, err := cmd.Flags().GetString("output")
		if err != nil {
			panic(err)
		}

		sqsTmpl := sqs.NewSQS(input, output)

		err = sqsTmpl.Build()
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(sqsCmd)

	sqsCmd.Flags().StringP("input", "i", "", "Path to the yaml file. For example: lambdas.yaml")
	sqsCmd.Flags().StringP("output", "o", "", "Path to the output file. For example: ./output/sqs.tf")

	sqsCmd.MarkFlagRequired("input")
	sqsCmd.MarkFlagRequired("output")
}
