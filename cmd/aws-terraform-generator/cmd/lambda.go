package cmd

import (
	"github.com/spf13/cobra"

	"github.com/joselitofilho/aws-terraform-generator/internal/generators/lambda"
)

const (
	lambdaCMDFlagConfig = "config"
	lambdaCMDFlagOutput = "output"
)

// lambdaCmd represents the lambda command
var lambdaCmd = &cobra.Command{
	Use:   "lambda",
	Short: "Manage Lambda",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := cmd.Flags().GetString(lambdaCMDFlagConfig)
		if err != nil {
			printErrorAndExit(err)
		}

		output, err := cmd.Flags().GetString(lambdaCMDFlagOutput)
		if err != nil {
			printErrorAndExit(err)
		}

		err = lambda.NewLambda(config, output).Build()
		if err != nil {
			printErrorAndExit(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(lambdaCmd)

	lambdaCmd.Flags().StringP(lambdaCMDFlagConfig, "c", "", "Path to the configuration file. For example: ./lambda.config.yaml")
	lambdaCmd.Flags().StringP(lambdaCMDFlagOutput, "o", "", "Path to the output folder. For example: ./output")

	lambdaCmd.MarkFlagRequired(lambdaCMDFlagConfig)
	lambdaCmd.MarkFlagRequired(lambdaCMDFlagOutput)
}
