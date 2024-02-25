package cmd

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/fatih/color"
	"github.com/joselitofilho/aws-terraform-generator/internal/guides"
	"github.com/spf13/cobra"

	"github.com/AlecAivazis/survey/v2"
)

const (
	rootCMDFlagWorkdir = "workdir"

	optionGuideDiagram          = "Generate a diagram config file"
	optionGuideInitialStructure = "Generate the initial structure"
	optionGuideCode             = "Generate code"
	optionExit                  = "Exit"
)

var (
	ErrNoDiagramOrConfigFiles = errors.New("this directory does not contain any diagram or config files")

	fmtRed = color.New(color.FgHiRed)
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "aws-terraform-generator",
	Short: "AWS terraform generator",
	Run: func(cmd *cobra.Command, args []string) {
		workdir, err := cmd.Flags().GetString(rootCMDFlagWorkdir)
		if err != nil {
			printErrorAndExit(err)
		}

		fmt.Println(`

		 ██████╗ ██████╗ ██████╗ ███████╗     ██████╗ ███████╗███╗   ██╗
		██╔════╝██╔═══██╗██╔══██╗██╔════╝    ██╔════╝ ██╔════╝████╗  ██║
		██║     ██║   ██║██║  ██║█████╗      ██║  ███╗█████╗  ██╔██╗ ██║
		██║     ██║   ██║██║  ██║██╔══╝      ██║   ██║██╔══╝  ██║╚██╗██║
		╚██████╗╚██████╔╝██████╔╝███████╗    ╚██████╔╝███████╗██║ ╚████║
		 ╚═════╝ ╚═════╝ ╚═════╝ ╚══════╝     ╚═════╝ ╚══════╝╚═╝  ╚═══╝

`)

		shouldContinue := true
		for shouldContinue {
			// Create a map to store file extensions and their respective files
			fileMap := make(map[string][]string)

			// Read files in the current directory
			files, err := os.ReadDir(workdir)
			if err != nil {
				printErrorAndExit(fmt.Errorf("error reading directory: %w", err))
			}

			// Iterate through the files and populate the map
			for _, file := range files {
				if !file.IsDir() {
					ext := strings.ToLower(path.Ext(file.Name()))

					switch ext {
					case ".xml":
						fileMap[diagramCMDFlagDiagram] = append(fileMap[diagramCMDFlagDiagram], file.Name())
					case ".yaml", ".yml":
						fileMap[diagramCMDFlagConfig] = append(fileMap[diagramCMDFlagConfig], file.Name())
					}
				}
			}

			if len(fileMap) == 0 {
				printErrorAndExit(ErrNoDiagramOrConfigFiles)
			}

			var commandName string
			if err := survey.AskOne(&survey.Select{
				Message: "What would you like to do?",
				Options: []string{
					optionGuideDiagram,
					optionGuideInitialStructure,
					optionGuideCode,
					optionExit,
				},
			}, &commandName); err != nil {
				printErrorAndExit(err)
			}

			switch commandName {
			case optionGuideDiagram:
				answers, err := guides.GuideDiagram(workdir, fileMap)
				if err != nil {
					printErrorAndExit(err)
				}

				_ = diagramCmd.Flags().Set(diagramCMDFlagDiagram, answers.Diagram)
				_ = diagramCmd.Flags().Set(diagramCMDFlagConfig, answers.Config)
				_ = diagramCmd.Flags().Set(diagramCMDFlagOutput, answers.Output)
				diagramCmd.Run(diagramCmd, []string{})
			case optionGuideInitialStructure:
				answers, err := guides.GuideStructure(workdir, fileMap)
				if err != nil {
					printErrorAndExit(err)
				}

				_ = structureCmd.Flags().Set(structureCMDFlagConfig, answers.Config)
				_ = structureCmd.Flags().Set(structureCMDFlagOutput, answers.Output)
				structureCmd.Run(structureCmd, []string{})
			case optionGuideCode:
				answers, err := guides.GuideCode(workdir, fileMap)
				if err != nil {
					printErrorAndExit(err)
				}

				stackOutput := fmt.Sprintf("%s/%s", answers.Output, answers.StackName)

				fmt.Println("→ Generating API Gateway code...")
				_ = apigatewayCmd.Flags().Set(apigatewayCMDFlagConfig, answers.Config)
				_ = apigatewayCmd.Flags().Set(apigatewayCMDFlagOutput, answers.Output)
				apigatewayCmd.Run(apigatewayCmd, []string{})
				fmt.Println()

				fmt.Println("→ Generating Lambda code...")
				_ = lambdaCmd.Flags().Set(lambdaCMDFlagConfig, answers.Config)
				_ = lambdaCmd.Flags().Set(lambdaCMDFlagOutput, stackOutput)
				lambdaCmd.Run(lambdaCmd, []string{})
				fmt.Println()

				fmt.Println("→ Generating S3 code...")
				_ = s3Cmd.Flags().Set(s3CMDFlagConfig, answers.Config)
				_ = s3Cmd.Flags().Set(s3CMDFlagOutput, stackOutput)
				s3Cmd.Run(s3Cmd, []string{})
				fmt.Println()

				fmt.Println("→ Generating SNS code...")
				_ = snsCmd.Flags().Set(snsCMDFlagConfig, answers.Config)
				_ = snsCmd.Flags().Set(snsCMDFlagOutput, stackOutput)
				snsCmd.Run(snsCmd, []string{})
				fmt.Println()

				fmt.Println("→ Generating SQS code...")
				_ = sqsCmd.Flags().Set(sqsCMDFlagConfig, answers.Config)
				_ = sqsCmd.Flags().Set(sqsCMDFlagOutput, stackOutput)
				sqsCmd.Run(sqsCmd, []string{})
			default:
				shouldContinue = false
			}

			if shouldContinue {
				fmt.Println()
			}
		}

		fmt.Println("👋 Goodbye. Until next time!")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP(rootCMDFlagWorkdir, "", ".",
		"Path to the directory where diagrams and configuration files are stored for the project. For example: ./example")
}

func printErrorAndExit(err error) {
	fmtRed.Printf("🚨 %s\n", err)
	os.Exit(1)
}
