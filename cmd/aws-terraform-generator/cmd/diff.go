package cmd

import (
	"github.com/spf13/cobra"

	"github.com/joselitofilho/aws-terraform-generator/internal/resources"
	"github.com/joselitofilho/aws-terraform-generator/internal/transformers/yamltoresources"
)

// diffCmd represents the sqs command.
var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Manage Diff",
	Run: func(cmd *cobra.Command, _ []string) {
		left, err := cmd.Flags().GetString(flagLeft)
		if err != nil {
			printErrorAndExit(err)
		}

		right, err := cmd.Flags().GetString(flagRight)
		if err != nil {
			printErrorAndExit(err)
		}

		leftRc, err := yamltoresources.Parse(left)
		if err != nil {
			printErrorAndExit(err)
		}

		rightRc, err := yamltoresources.Parse(right)
		if err != nil {
			printErrorAndExit(err)
		}

		resources.PrintDiff(leftRc, rightRc)
	},
}

func init() {
	rootCmd.AddCommand(diffCmd)

	diffCmd.Flags().StringP(flagLeft, "l", "", "Path to the left YAML config file. For example: ./left.yaml")
	diffCmd.Flags().StringP(flagRight, "r", "", "Path to the right YAML config file. For example: ./right.yaml")

	_ = diffCmd.MarkFlagRequired(flagLeft)
	_ = diffCmd.MarkFlagRequired(flagRight)
}
