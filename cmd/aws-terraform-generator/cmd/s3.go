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
		input, err := cmd.Flags().GetString("input")
		if err != nil {
			panic(err)
		}

		output, err := cmd.Flags().GetString("output")
		if err != nil {
			panic(err)
		}

		s3Tmpl := s3.NewS3(input, output)

		err = s3Tmpl.Build()
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(s3Cmd)

	s3Cmd.Flags().StringP("input", "i", "", "Path to the yaml file. For example: lambdas.yaml")
	s3Cmd.Flags().StringP("output", "o", "", "Path to the output file. For example: ./output/s3.tf")

	s3Cmd.MarkFlagRequired("input")
	s3Cmd.MarkFlagRequired("output")
}
