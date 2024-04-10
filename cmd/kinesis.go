package cmd

import (
	"github.com/spf13/cobra"

	"github.com/joselitofilho/aws-terraform-generator/internal/generators/kinesis"
)

// kinesisCmd represents the sqs command.
var kinesisCmd = &cobra.Command{
	Use:   "kinesis",
	Short: "Manage Kinesis streams",
	Run: func(cmd *cobra.Command, _ []string) {
		config, _ := cmd.Flags().GetString(flagConfig)
		output, _ := cmd.Flags().GetString(flagOutput)

		if err := kinesis.NewKinesis(config, output).Build(); err != nil {
			printErrorAndExit(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(kinesisCmd)

	kinesisCmd.Flags().StringP(flagConfig, "c", "", "Path to the configuration file. For example: ./kinesis.config.yaml")
	kinesisCmd.Flags().StringP(flagOutput, "o", "", "Path to the output folder. For example: ./output")

	_ = kinesisCmd.MarkFlagRequired(flagConfig)
	_ = kinesisCmd.MarkFlagRequired(flagOutput)
}
