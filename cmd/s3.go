package cmd

import (
	"github.com/spf13/cobra"

	"github.com/joselitofilho/aws-terraform-generator/internal/generators/s3"
)

// s3Cmd represents the sqs command.
var s3Cmd = &cobra.Command{
	Use:   "s3",
	Short: "Manage S3",
	Run: func(cmd *cobra.Command, _ []string) {
		config, _ := cmd.Flags().GetString(flagConfig)
		output, _ := cmd.Flags().GetString(flagOutput)

		err := s3.NewS3(config, output).Build()
		if err != nil {
			printErrorAndExit(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(s3Cmd)

	s3Cmd.Flags().StringP(flagConfig, "c", "", "Path to the configuration file. For example: ./s3.config.yaml")
	s3Cmd.Flags().StringP(flagOutput, "o", "", "Path to the output folder. For example: ./output")

	_ = s3Cmd.MarkFlagRequired(flagConfig)
	_ = s3Cmd.MarkFlagRequired(flagOutput)
}
