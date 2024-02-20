package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/joselitofilho/aws-terraform-generator/internal/drawio"
	"github.com/joselitofilho/aws-terraform-generator/internal/templates/config"
	"github.com/joselitofilho/aws-terraform-generator/internal/transformers"
)

// diagramCmd represents the structure command
var diagramCmd = &cobra.Command{
	Use:   "diagram",
	Short: "Manage Diagram",
	Run: func(cmd *cobra.Command, args []string) {
		stackName, err := cmd.Flags().GetString("stack")
		if err != nil {
			panic(err)
		}

		diagram, err := cmd.Flags().GetString("diagram")
		if err != nil {
			panic(err)
		}

		configFile, err := cmd.Flags().GetString("config")
		if err != nil {
			panic(err)
		}

		output, err := cmd.Flags().GetString("output")
		if err != nil {
			panic(err)
		}

		root, err := drawio.Parse(diagram)
		if err != nil {
			panic(err)
		}

		resources, err := drawio.ParseResources(root)
		if err != nil {
			panic(err)
		}

		yamlConfigOut, err := transformers.TransformDrawIOToYAML(stackName, resources)
		if err != nil {
			panic(err)
		}

		//

		if configFile != "" {
			yamlParser := config.NewYAML(configFile)

			yamlConfig, err := yamlParser.Parse()
			if err != nil {
				panic(err)
			}

			for i := range yamlConfigOut.Lambdas {
				yamlConfigOut.Lambdas[i].ModuleLambdaSource = yamlConfig.Diagram.Modules.Lambda
			}

			for i, g := range yamlConfigOut.APIGateways {
				for j := range g.Lambdas {
					yamlConfigOut.APIGateways[i].Lambdas[j].ModuleLambdaSource = yamlConfig.Diagram.Modules.Lambda
				}
			}
		}

		//

		data, err := yaml.Marshal(yamlConfigOut)
		if err != nil {
			panic(err)
		}

		err = os.WriteFile(output, data, os.ModePerm)
		if err != nil {
			panic(err)
		}

		fmt.Printf("YAML '%s' has been generated successfully\n", output)
	},
}

func init() {
	rootCmd.AddCommand(diagramCmd)

	diagramCmd.Flags().StringP("stack", "s", "", "Stack name")
	diagramCmd.Flags().StringP("diagram", "d", "", "Path to the xml file. For example: ./diagram.xml")
	diagramCmd.Flags().StringP("config", "c", "", "Path to the YAML config file. For example: ./diagram.config.yaml")
	diagramCmd.Flags().StringP("output", "o", "", "Path to the output file. For example: ./diagram.yaml")

	diagramCmd.MarkFlagRequired("stack")
	diagramCmd.MarkFlagRequired("diagram")
	diagramCmd.MarkFlagRequired("output")
}
