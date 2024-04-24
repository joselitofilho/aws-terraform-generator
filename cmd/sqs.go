package cmd

import (
	"github.com/spf13/cobra"

	"github.com/joselitofilho/aws-terraform-generator/internal/generators/sqs"
)

// sqsCmd represents the sqs command.
var sqsCmd = &cobra.Command{
	Use:   "sqs",
	Short: "Manage SQS",
	Run: func(cmd *cobra.Command, _ []string) {
		config, _ := cmd.Flags().GetString(flagConfig)
		output, _ := cmd.Flags().GetString(flagOutput)

		err := sqs.NewSQS(config, output).Build()
		if err != nil {
			printErrorAndExit(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(sqsCmd)

	sqsCmd.Flags().StringP(flagConfig, "c", "", "Path to the configuration file. For example: ./sqs.config.yaml")
	sqsCmd.Flags().StringP(flagOutput, "o", "", "Path to the output folder. For example: ./output")

	_ = sqsCmd.MarkFlagRequired(flagConfig)
	_ = sqsCmd.MarkFlagRequired(flagOutput)
}
