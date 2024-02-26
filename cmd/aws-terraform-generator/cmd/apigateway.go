//nolint:dupl // That is a false positive
package cmd

import (
	"github.com/spf13/cobra"

	"github.com/joselitofilho/aws-terraform-generator/internal/generators/apigateway"
)

// apigatewayCmd represents the apigateway command
var apigatewayCmd = &cobra.Command{
	Use:   "apigateway",
	Short: "Manage APIGateway",
	Run: func(cmd *cobra.Command, _ []string) {
		config, err := cmd.Flags().GetString(flagConfig)
		if err != nil {
			printErrorAndExit(err)
		}

		output, err := cmd.Flags().GetString(flagOutput)
		if err != nil {
			printErrorAndExit(err)
		}

		err = apigateway.NewAPIGateway(config, output).Build()
		if err != nil {
			printErrorAndExit(err)
		}
	},
}

//nolint:gochecknoinits // That is the way the cobra library operates
func init() {
	rootCmd.AddCommand(apigatewayCmd)

	apigatewayCmd.Flags().StringP(flagConfig, "c", "",
		"Path to the configuration file. For example: ./apigateway.config.yaml")
	apigatewayCmd.Flags().StringP(flagOutput, "o", "",
		"Path to the output folder. For example: ./output")

	_ = apigatewayCmd.MarkFlagRequired(flagConfig)
	_ = apigatewayCmd.MarkFlagRequired(flagOutput)
}
