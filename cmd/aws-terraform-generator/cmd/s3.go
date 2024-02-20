package cmd

import (
	"github.com/spf13/cobra"

	"github.com/joselitofilho/aws-terraform-generator/internal/templates/s3"
)

// s3Cmd represents the sqs command
var s3Cmd = &cobra.Command{
	Use:   "s3",
	Short: "Manage S3",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := cmd.Flags().GetString("config")
		if err != nil {
			panic(err)
		}

		output, err := cmd.Flags().GetString("output")
		if err != nil {
			panic(err)
		}

		s3Tmpl := s3.NewS3(config, output)

		err = s3Tmpl.Build()
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(s3Cmd)

	s3Cmd.Flags().StringP("config", "c", "", "Path to the configuration file. For example: ./s3.config.yaml")
	s3Cmd.Flags().StringP("output", "o", "", "Path to the output file. For example: ./output/s3.tf")

	s3Cmd.MarkFlagRequired("config")
	s3Cmd.MarkFlagRequired("output")
}
