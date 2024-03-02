package cmd

import (
	"github.com/spf13/cobra"

	"github.com/joselitofilho/aws-terraform-generator/internal/generators/lambda"
)

// lambdaCmd represents the lambda command.
var lambdaCmd = &cobra.Command{
	Use:   "lambda",
	Short: "Manage Lambda",
	Run: func(cmd *cobra.Command, _ []string) {
		config, err := cmd.Flags().GetString(flagConfig)
		if err != nil {
			printErrorAndExit(err)
		}

		output, err := cmd.Flags().GetString(flagOutput)
		if err != nil {
			printErrorAndExit(err)
		}

		err = lambda.NewLambda(config, output).Build()
		if err != nil {
			printErrorAndExit(err)
		}
	},
}

//nolint:gochecknoinits // That is the way the cobra library operates
func init() {
	rootCmd.AddCommand(lambdaCmd)

	lambdaCmd.Flags().StringP(flagConfig, "c", "",
		"Path to the configuration file. For example: ./lambda.config.yaml")
	lambdaCmd.Flags().StringP(flagOutput, "o", "",
		"Path to the output folder. For example: ./output")

	_ = lambdaCmd.MarkFlagRequired(flagConfig)
	_ = lambdaCmd.MarkFlagRequired(flagOutput)
}
