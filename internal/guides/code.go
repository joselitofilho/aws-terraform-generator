package guides

import (
	"fmt"
	"path"
	"strings"

	"github.com/AlecAivazis/survey/v2"

	surveyasker "github.com/joselitofilho/aws-terraform-generator/internal/survey"
)

type CodeAnswers struct {
	StackName string
	Config    string
	Output    string
}

func GuideCode(surveyAsker surveyasker.SurveyAsker, workdir string, fileMap map[string][]string) (*CodeAnswers, error) {
	if len(fileMap["config"]) == 0 {
		return nil, ErrDirDoesNotContainAnyConfigFile
	}

	configOptions := make([]string, 0, len(fileMap["config"]))
	configOptions = append(configOptions, fileMap["config"]...)

	answers := CodeAnswers{}

	_, defaultStackName := path.Split(workdir)

	if err := surveyAsker.Ask([]*survey.Question{
		{
			Name: "stackName",
			Prompt: &survey.Input{
				Message: "Enter the stack name:",
				Default: defaultStackName,
			},
			Validate: survey.Required,
		},
	}, &answers); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	var defaultConfigOption int
	for i := range configOptions {
		if strings.Contains(configOptions[i], "diagram.yaml") ||
			strings.Contains(configOptions[i], "diagram.yml") {
			defaultConfigOption = i
		}
	}

	if err := surveyAsker.Ask([]*survey.Question{
		{
			Name: "config",
			Prompt: &survey.Select{
				Message: "Choose a config:",
				Default: defaultConfigOption,
				Options: configOptions,
			},
		},
	}, &answers); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	answers.Config = replaceDoubleSlash(path.Join(workdir, answers.Config))

	if err := surveyAsker.Ask([]*survey.Question{
		{
			Name: "output",
			Prompt: &survey.Input{
				Message: "Enter the output folder:",
				Default: "./output",
			},
		},
	}, &answers); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return &answers, nil
}
