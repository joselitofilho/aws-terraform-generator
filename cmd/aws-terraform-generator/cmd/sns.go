package cmd

import (
	"github.com/joselitofilho/aws-terraform-generator/internal/generators/sns"
	"github.com/spf13/cobra"
)

// snsCmd represents the sqs command
var snsCmd = &cobra.Command{
	Use:   "sns",
	Short: "Manage SNS",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := cmd.Flags().GetString("config")
		if err != nil {
			panic(err)
		}

		output, err := cmd.Flags().GetString("output")
		if err != nil {
			panic(err)
		}

		err = sns.NewSNS(config, output).Build()
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(snsCmd)

	snsCmd.Flags().StringP("config", "c", "", "Path to the configuration file. For example: ./sns.config.yaml")
	snsCmd.Flags().StringP("output", "o", "", "Path to the output file. For example: ./output/sns.tf")

	snsCmd.MarkFlagRequired("config")
	snsCmd.MarkFlagRequired("output")
}
