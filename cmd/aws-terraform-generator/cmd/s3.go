package cmd

import (
	"github.com/spf13/cobra"

	"github.com/joselitofilho/aws-terraform-generator/internal/generators/s3"
)

const (
	s3CMDFlagConfig = "config"
	s3CMDFlagOutput = "output"
)

// s3Cmd represents the sqs command
var s3Cmd = &cobra.Command{
	Use:   "s3",
	Short: "Manage S3",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := cmd.Flags().GetString(s3CMDFlagConfig)
		if err != nil {
			printErrorAndExit(err)
		}

		output, err := cmd.Flags().GetString(s3CMDFlagOutput)
		if err != nil {
			printErrorAndExit(err)
		}

		s3Tmpl := s3.NewS3(config, output)

		err = s3Tmpl.Build()
		if err != nil {
			printErrorAndExit(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(s3Cmd)

	s3Cmd.Flags().StringP(s3CMDFlagConfig, "c", "", "Path to the configuration file. For example: ./s3.config.yaml")
	s3Cmd.Flags().StringP(s3CMDFlagOutput, "o", "", "Path to the output file. For example: ./output/s3.tf")

	s3Cmd.MarkFlagRequired(s3CMDFlagConfig)
	s3Cmd.MarkFlagRequired(s3CMDFlagOutput)
}
