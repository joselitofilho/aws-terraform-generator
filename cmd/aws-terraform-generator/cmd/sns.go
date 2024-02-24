package cmd

import (
	"github.com/joselitofilho/aws-terraform-generator/internal/generators/sns"
	"github.com/spf13/cobra"
)

const (
	snsCMDFlagConfig = "config"
	snsCMDFlagOutput = "output"
)

// snsCmd represents the sqs command
var snsCmd = &cobra.Command{
	Use:   "sns",
	Short: "Manage SNS",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := cmd.Flags().GetString(snsCMDFlagConfig)
		if err != nil {
			printErrorAndExit(err)
		}

		output, err := cmd.Flags().GetString(snsCMDFlagOutput)
		if err != nil {
			printErrorAndExit(err)
		}

		err = sns.NewSNS(config, output).Build()
		if err != nil {
			printErrorAndExit(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(snsCmd)

	snsCmd.Flags().StringP(snsCMDFlagConfig, "c", "", "Path to the configuration file. For example: ./sns.config.yaml")
	snsCmd.Flags().StringP(snsCMDFlagOutput, "o", "", "Path to the output file. For example: ./output/sns.tf")

	snsCmd.MarkFlagRequired(snsCMDFlagConfig)
	snsCmd.MarkFlagRequired(snsCMDFlagOutput)
}
