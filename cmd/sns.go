package cmd

import (
	"github.com/spf13/cobra"

	"github.com/joselitofilho/aws-terraform-generator/internal/generators/sns"
)

// snsCmd represents the sqs command.
var snsCmd = &cobra.Command{
	Use:   "sns",
	Short: "Manage SNS",
	Run: func(cmd *cobra.Command, _ []string) {
		config, _ := cmd.Flags().GetString(flagConfig)
		output, _ := cmd.Flags().GetString(flagOutput)

		err := sns.NewSNS(config, output).Build()
		if err != nil {
			printErrorAndExit(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(snsCmd)

	snsCmd.Flags().StringP(flagConfig, "c", "", "Path to the configuration file. For example: ./sns.config.yaml")
	snsCmd.Flags().StringP(flagOutput, "o", "", "Path to the output folder. For example: ./output")

	_ = snsCmd.MarkFlagRequired(flagConfig)
	_ = snsCmd.MarkFlagRequired(flagOutput)
}
