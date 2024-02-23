package cmd

import (
	"github.com/spf13/cobra"

	"github.com/joselitofilho/aws-terraform-generator/internal/generators/lambda"
)

// lambdaCmd represents the lambda command
var lambdaCmd = &cobra.Command{
	Use:   "lambda",
	Short: "Manage Lambda",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := cmd.Flags().GetString("config")
		if err != nil {
			panic(err)
		}

		output, err := cmd.Flags().GetString("output")
		if err != nil {
			panic(err)
		}

		err = lambda.NewLambda(config, output).Build()
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(lambdaCmd)

	lambdaCmd.Flags().StringP("config", "c", "", "Path to the configuration file. For example: ./lambda.config.yaml")
	lambdaCmd.Flags().StringP("output", "o", "", "Path to the output folder. For example: ./output")

	lambdaCmd.MarkFlagRequired("config")
	lambdaCmd.MarkFlagRequired("output")
}
