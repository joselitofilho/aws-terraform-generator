package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// lambdaCmd represents the lambda command
var lambdaCmd = &cobra.Command{
	Use:   "lambda",
	Short: "Manage Lambda",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("lambda called")
	},
}

func init() {
	rootCmd.AddCommand(lambdaCmd)

	lambdaCmd.Flags().StringP("name", "n", "", "Name of the Lambda")
	lambdaCmd.Flags().StringP("description", "d", "", "Description of the Lambda")

	lambdaCmd.MarkFlagRequired("name")
}
