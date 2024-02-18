package cmd

import (
	"github.com/spf13/cobra"

	"github.com/joselitofilho/aws-terraform-generator/internal/templates/apigateway"
)

// apigatewayCmd represents the apigateway command
var apigatewayCmd = &cobra.Command{
	Use:   "apigateway",
	Short: "Manage APIGateway",
	Run: func(cmd *cobra.Command, args []string) {
		input, err := cmd.Flags().GetString("input")
		if err != nil {
			panic(err)
		}

		output, err := cmd.Flags().GetString("output")
		if err != nil {
			panic(err)
		}

		err = apigateway.NewAPIGateway(input, output).Build()
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(apigatewayCmd)

	apigatewayCmd.Flags().StringP("input", "i", "", "Path to the yaml file. For example: lambdas.yaml")
	apigatewayCmd.Flags().StringP("output", "o", "", "Path to the output folder. For example: ./output")

	apigatewayCmd.MarkFlagRequired("input")
	apigatewayCmd.MarkFlagRequired("output")
}
