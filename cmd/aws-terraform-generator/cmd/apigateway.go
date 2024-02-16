package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// apigatewayCmd represents the apigateway command
var apigatewayCmd = &cobra.Command{
	Use:   "apigateway",
	Short: "Manage API Gateway",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("apigateway called")
	},
}

func init() {
	rootCmd.AddCommand(apigatewayCmd)
}
