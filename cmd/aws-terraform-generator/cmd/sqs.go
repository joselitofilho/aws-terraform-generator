package cmd

import (
	"github.com/joselitofilho/aws-terraform-generator/internal/templates/sqs"
	"github.com/spf13/cobra"
)

// sqsCmd represents the sqs command
var sqsCmd = &cobra.Command{
	Use:   "sqs",
	Short: "Manage SQS",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := cmd.Flags().GetString("config")
		if err != nil {
			panic(err)
		}

		output, err := cmd.Flags().GetString("output")
		if err != nil {
			panic(err)
		}

		sqsTmpl := sqs.NewSQS(config, output)

		err = sqsTmpl.Build()
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(sqsCmd)

	sqsCmd.Flags().StringP("config", "c", "", "Path to the configuration file. For example: ./sqs.config.yaml")
	sqsCmd.Flags().StringP("output", "o", "", "Path to the output file. For example: ./output/sqs.tf")

	sqsCmd.MarkFlagRequired("config")
	sqsCmd.MarkFlagRequired("output")
}
