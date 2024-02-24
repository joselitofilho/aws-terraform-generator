package guides

import (
	"fmt"
	"path"
	"strings"

	"github.com/AlecAivazis/survey/v2"
)

type CodeAnswers struct {
	StackName string
	Config    string
	Output    string
}

func GuideCode(workdir string, fileMap map[string][]string) (*CodeAnswers, error) {
	if len(fileMap["config"]) == 0 {
		return nil, ErrDirDoesNotContainAnyConfigFile
	}

	configOptions := make([]string, 0, len(fileMap["config"]))
	configOptions = append(configOptions, fileMap["config"]...)

	answers := CodeAnswers{}

	_, defaultStackName := path.Split(workdir)

	if err := survey.Ask([]*survey.Question{
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

	if err := survey.Ask([]*survey.Question{
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

	answers.Config = replaceDoubleSlash(fmt.Sprintf("%s/%s", workdir, answers.Config))

	if err := survey.AskOne(
		&survey.Input{
			Message: "Enter the output folder:",
			Default: "./output",
		},
		&answers.Output); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	answers.Output = replaceDoubleSlash(fmt.Sprintf("%s/%s", answers.Output, answers.StackName))

	return &answers, nil
}
