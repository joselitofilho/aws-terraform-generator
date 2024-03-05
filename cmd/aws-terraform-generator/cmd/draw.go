package cmd

import (
	"encoding/xml"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/joselitofilho/aws-terraform-generator/internal/generators/terraform"
)

// drawCmd represents the sqs command
var drawCmd = &cobra.Command{
	Use:   "draw",
	Short: "Manage Draw",
	Run: func(cmd *cobra.Command, _ []string) {
		directory, err := cmd.Flags().GetString("directory")
		if err != nil {
			panic(err)
		}

		output, err := cmd.Flags().GetString("output")
		if err != nil {
			panic(err)
		}

		config, err := terraform.ParseTerraformFiles(directory)
		if err != nil {
			fmt.Printf("Error parsing Terraform files: %s\n", err)
			return
		}

		xmlData, err := xml.MarshalIndent(config, "", "    ")
		if err != nil {
			fmt.Printf("Error marshalling XML: %s\n", err)
			return
		}

		err = os.WriteFile(output, xmlData, os.ModePerm)
		if err != nil {
			fmt.Printf("Error writing XML to file: %s\n", err)
			return
		}

		fmt.Println("XML representation of Terraform configuration has been written to", output)
	},
}

func init() {
	rootCmd.AddCommand(drawCmd)

	drawCmd.Flags().StringP("directory", "d", "", "Path to the infra folder. For example: ./output/mystack")
	drawCmd.Flags().StringP("output", "o", "", "Path to the output file. For example: ./output/diagram.xml")

	_ = drawCmd.MarkFlagRequired("output")
}
