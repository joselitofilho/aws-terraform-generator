//nolint:dupl // That is a false positive
package cmd

import (
	"github.com/spf13/cobra"

	"github.com/joselitofilho/aws-terraform-generator/internal/generators/sqs"
)

// sqsCmd represents the sqs command
var sqsCmd = &cobra.Command{
	Use:   "sqs",
	Short: "Manage SQS",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := cmd.Flags().GetString(flagConfig)
		if err != nil {
			printErrorAndExit(err)
		}

		output, err := cmd.Flags().GetString(flagOutput)
		if err != nil {
			printErrorAndExit(err)
		}

		err = sqs.NewSQS(config, output).Build()
		if err != nil {
			printErrorAndExit(err)
		}
	},
}

//nolint:gochecknoinits // That is the way the cobra library operates
func init() {
	rootCmd.AddCommand(sqsCmd)

	sqsCmd.Flags().StringP(flagConfig, "c", "", "Path to the configuration file. For example: ./sqs.config.yaml")
	sqsCmd.Flags().StringP(flagOutput, "o", "", "Path to the output folder. For example: ./output")

	_ = sqsCmd.MarkFlagRequired(flagConfig)
	_ = sqsCmd.MarkFlagRequired(flagOutput)
}
