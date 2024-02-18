package cmd

import (
	"github.com/spf13/cobra"

	"github.com/joselitofilho/aws-terraform-generator/internal/templates/lambda"
)

// lambdaCmd represents the lambda command
var lambdaCmd = &cobra.Command{
	Use:   "lambda",
	Short: "Manage Lambda",
	Run: func(cmd *cobra.Command, args []string) {
		input, err := cmd.Flags().GetString("input")
		if err != nil {
			panic(err)
		}

		output, err := cmd.Flags().GetString("output")
		if err != nil {
			panic(err)
		}

		err = lambda.NewLambda(input, output).Build()
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(lambdaCmd)

	lambdaCmd.Flags().StringP("input", "i", "", "Path to the yaml file. For example: lambdas.yaml")
	lambdaCmd.Flags().StringP("output", "o", "", "Path to the output folder. For example: ./output")

	lambdaCmd.MarkFlagRequired("input")
	lambdaCmd.MarkFlagRequired("output")
}
