package cmd

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/joselitofilho/aws-terraform-generator/internal/guides"
	"github.com/spf13/cobra"

	"github.com/AlecAivazis/survey/v2"
)

const (
	rootCMDFlagWorkdir = "workdir"

	guideDiagram = "Generate a diagram config file"
)

var (
	ErrNoDiagramOrConfigFiles = errors.New("this directory does not contain any diagram or config files")
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "aws-terraform-generator",
	Short: "AWS terraform generator",
	Run: func(cmd *cobra.Command, args []string) {
		workdir, err := cmd.Flags().GetString(rootCMDFlagWorkdir)
		if err != nil {
			panic(err)
		}

		// Create a map to store file extensions and their respective files
		fileMap := make(map[string][]string)

		// Read files in the current directory
		files, err := os.ReadDir(workdir)
		if err != nil {
			panic(fmt.Errorf("error reading directory: %w", err))
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
			panic(ErrNoDiagramOrConfigFiles)
		}

		var commandName string
		if err := survey.AskOne(&survey.Select{
			Message: "What would you like to do?",
			Options: []string{
				guideDiagram,
				"Generate a structure",
				"Generate API gateway",
				"Generate Lambda",
				"Generate SQS queue",
				"Generate S3 bucket",
			},
		}, &commandName); err != nil {
			panic(err)
		}

		switch commandName {
		case guideDiagram:
			answers, err := guides.GuideDiagram(workdir, fileMap)
			if err != nil {
				panic(err)
			}

			_ = diagramCmd.Flags().Set(diagramCMDFlagDiagram, answers.Diagram)
			_ = diagramCmd.Flags().Set(diagramCMDFlagConfig, answers.Config)
			_ = diagramCmd.Flags().Set(diagramCMDFlagOutput, answers.Output)
			diagramCmd.Run(diagramCmd, []string{})
		}
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
	rootCmd.Flags().StringP(rootCMDFlagWorkdir, "d", ".",
		"Path to the directory where diagrams and configuration files are stored for the project. For example: ./example")
}
