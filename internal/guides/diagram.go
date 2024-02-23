package guides

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
)

var (
	ErrDirDoesNotContainAnyConfigFile  = errors.New("this directory does not contain any config (.yaml|.yml) file")
	ErrDirDoesNotContainAnyDiagramFile = errors.New("this directory does not contain any diagram (.xml) file")
)

type DiagramAnswers struct {
	Diagram string
	Config  string
	Output  string
}

func GuideDiagram(workdir string, fileMap map[string][]string) (*DiagramAnswers, error) {
	if len(fileMap["diagram"]) == 0 {
		return nil, ErrDirDoesNotContainAnyDiagramFile
	}

	if len(fileMap["config"]) == 0 {
		fmt.Println("This directory does not contain any config (.yaml|.yml) file")
		os.Exit(1)
	}

	diagramOptions := make([]string, 0, len(fileMap["diagram"]))
	diagramOptions = append(diagramOptions, fileMap["diagram"]...)

	configOptions := make([]string, 0, len(fileMap["config"]))
	configOptions = append(configOptions, fileMap["config"]...)

	answers := &DiagramAnswers{}

	if err := survey.Ask([]*survey.Question{
		{
			Name: "diagram",
			Prompt: &survey.Select{
				Message: "Choose a diagram:",
				Options: diagramOptions,
			},
			Validate: survey.Required,
		},
		{
			Name: "config",
			Prompt: &survey.Select{
				Message: "Choose a config:",
				Options: configOptions,
			},
		},
	}, &answers); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	if err := survey.AskOne(&survey.Input{
		Message: "Enter the output file:",
		Default: replaceDoubleSlash(fmt.Sprintf("%s/%s.yaml", workdir, strings.TrimSuffix(answers.Diagram, ".xml")))},
		&answers.Output); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	answers.Diagram = replaceDoubleSlash(fmt.Sprintf("%s/%s", workdir, answers.Diagram))
	answers.Config = replaceDoubleSlash(fmt.Sprintf("%s/%s", workdir, answers.Config))

	return answers, nil
}

func replaceDoubleSlash(str string) string {
	return strings.ReplaceAll(str, "//", "/")
}
