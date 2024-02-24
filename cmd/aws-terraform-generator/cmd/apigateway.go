package cmd

import (
	"github.com/spf13/cobra"

	"github.com/joselitofilho/aws-terraform-generator/internal/generators/apigateway"
)

const (
	apigatewayCMDFlagConfig = "config"
	apigatewayCMDFlagOutput = "output"
)

// apigatewayCmd represents the apigateway command
var apigatewayCmd = &cobra.Command{
	Use:   "apigateway",
	Short: "Manage APIGateway",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := cmd.Flags().GetString(apigatewayCMDFlagConfig)
		if err != nil {
			printErrorAndExit(err)
		}

		output, err := cmd.Flags().GetString(apigatewayCMDFlagOutput)
		if err != nil {
			printErrorAndExit(err)
		}

		err = apigateway.NewAPIGateway(config, output).Build()
		if err != nil {
			printErrorAndExit(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(apigatewayCmd)

	apigatewayCmd.Flags().StringP(apigatewayCMDFlagConfig, "c", "", "Path to the configuration file. For example: ./apigateway.config.yaml")
	apigatewayCmd.Flags().StringP(apigatewayCMDFlagOutput, "o", "", "Path to the output folder. For example: ./output")

	apigatewayCmd.MarkFlagRequired(apigatewayCMDFlagConfig)
	apigatewayCmd.MarkFlagRequired(apigatewayCMDFlagOutput)
}
