package cmd

import (
	"github.com/spf13/cobra"

	"github.com/joselitofilho/aws-terraform-generator/internal/generators/draw"
)

// drawCmd represents the sqs command.
var drawCmd = &cobra.Command{
	Use:   "draw",
	Short: "Manage Draw",
	Run: func(cmd *cobra.Command, _ []string) {
		workdirs, err := cmd.Flags().GetStringArray(flagWorkdir)
		if err != nil {
			printErrorAndExit(err)
		}

		files, err := cmd.Flags().GetStringArray(flagFile)
		if err != nil {
			printErrorAndExit(err)
		}

		configFilename, err := cmd.Flags().GetString(flagConfig)
		if err != nil {
			printErrorAndExit(err)
		}

		output, err := cmd.Flags().GetString(flagOutput)
		if err != nil {
			printErrorAndExit(err)
		}

		err = draw.NewDraw(workdirs, files, configFilename, output).Build()
		if err != nil {
			printErrorAndExit(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(drawCmd)

	drawCmd.Flags().StringArrayP(flagWorkdir, "", []string{"."},
		"Path to the folder where the terraform files are. For example: ./workdir")
	drawCmd.Flags().StringArrayP(flagFile, "", nil,
		"Path to the specific terraform file. For example: ./workdir/sqs.tf")
	drawCmd.Flags().StringP(flagConfig, "c", "",
		"Path to the YAML config file. For example: ./draw.config.yaml")
	drawCmd.Flags().StringP(flagOutput, "o", "", "Path to the output folder. For example: ./output")

	_ = drawCmd.MarkFlagRequired(flagConfig)
	_ = drawCmd.MarkFlagRequired(flagOutput)
}
