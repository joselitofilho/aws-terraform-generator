package cmd

import (
	"github.com/spf13/cobra"

	"github.com/joselitofilho/aws-terraform-generator/internal/generators/apigateway"
)

// apigatewayCmd represents the apigateway command
var apigatewayCmd = &cobra.Command{
	Use:   "apigateway",
	Short: "Manage APIGateway",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := cmd.Flags().GetString("config")
		if err != nil {
			panic(err)
		}

		output, err := cmd.Flags().GetString("output")
		if err != nil {
			panic(err)
		}

		err = apigateway.NewAPIGateway(config, output).Build()
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(apigatewayCmd)

	apigatewayCmd.Flags().StringP("config", "c", "", "Path to the configuration file. For example: ./apigateway.config.yaml")
	apigatewayCmd.Flags().StringP("output", "o", "", "Path to the output folder. For example: ./output")

	apigatewayCmd.MarkFlagRequired("config")
	apigatewayCmd.MarkFlagRequired("output")
}
