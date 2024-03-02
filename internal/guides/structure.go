package guides

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"

	surveyasker "github.com/joselitofilho/aws-terraform-generator/internal/survey"
)

type StructureAnswers struct {
	Config string
	Output string
}

func GuideStructure(
	surveyAsker surveyasker.Asker, workdir string, fileMap map[string][]string,
) (*StructureAnswers, error) {
	if len(fileMap["config"]) == 0 {
		return nil, ErrDirDoesNotContainAnyConfigFile
	}

	configOptions := make([]string, 0, len(fileMap["config"]))
	configOptions = append(configOptions, fileMap["config"]...)

	answers := StructureAnswers{}

	var defaultConfigOption int

	for i := range configOptions {
		if strings.Contains(configOptions[i], "structure.config") {
			defaultConfigOption = i
		}
	}

	if err := surveyAsker.AskOne(&survey.Select{
		Message: "Choose a config:",
		Default: defaultConfigOption,
		Options: configOptions,
	}, &answers.Config); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	answers.Config = replaceDoubleSlash(fmt.Sprintf("%s/%s", workdir, answers.Config))

	if err := surveyAsker.AskOne(
		&survey.Input{
			Message: "Enter the output folder:",
			Default: "./output",
		},
		&answers.Output); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	answers.Output = replaceDoubleSlash(answers.Output)

	return &answers, nil
}
