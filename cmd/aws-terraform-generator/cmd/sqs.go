package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// sqsCmd represents the sqs command
var sqsCmd = &cobra.Command{
	Use:   "sqs",
	Short: "Manage SQS",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("sqs called")
	},
}

func init() {
	rootCmd.AddCommand(sqsCmd)

	sqsCmd.Flags().StringP("name", "n", "", "Name of the SQS")

	sqsCmd.MarkFlagRequired("name")
}
