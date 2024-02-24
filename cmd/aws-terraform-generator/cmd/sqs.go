package cmd

import (
	"github.com/joselitofilho/aws-terraform-generator/internal/generators/sqs"
	"github.com/spf13/cobra"
)

const (
	sqsCMDFlagConfig = "config"
	sqsCMDFlagOutput = "output"
)

// sqsCmd represents the sqs command
var sqsCmd = &cobra.Command{
	Use:   "sqs",
	Short: "Manage SQS",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := cmd.Flags().GetString(sqsCMDFlagConfig)
		if err != nil {
			printErrorAndExit(err)
		}

		output, err := cmd.Flags().GetString(sqsCMDFlagOutput)
		if err != nil {
			printErrorAndExit(err)
		}

		err = sqs.NewSQS(config, output).Build()
		if err != nil {
			printErrorAndExit(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(sqsCmd)

	sqsCmd.Flags().StringP(sqsCMDFlagConfig, "c", "", "Path to the configuration file. For example: ./sqs.config.yaml")
	sqsCmd.Flags().StringP(sqsCMDFlagOutput, "o", "", "Path to the output folder. For example: ./output")

	sqsCmd.MarkFlagRequired(sqsCMDFlagConfig)
	sqsCmd.MarkFlagRequired(sqsCMDFlagOutput)
}
