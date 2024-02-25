package structure

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/joselitofilho/aws-terraform-generator/internal/generators"
	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
)

type Structure struct {
	config string
	output string
}

func NewStructure(config, output string) *Structure {
	return &Structure{config: config, output: output}
}

func (s *Structure) Build() error {
	yamlParser := config.NewYAML(s.config)

	yamlConfig, err := yamlParser.Parse()
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	defaultTemplatesMap := map[string]string{}
	for i := range yamlConfig.Structure.DefaultTemplates {
		defaultTemplatesMap = yamlConfig.Structure.DefaultTemplates[i]
	}

	for i := range yamlConfig.Structure.Stacks {
		conf := yamlConfig.Structure.Stacks[i]

		data := Data{
			StackName: conf.Name,
		}

		for _, folder := range conf.Folders {
			output := fmt.Sprintf("%s/%s/%s", s.output, conf.Name, folder.Name)
			_ = os.MkdirAll(output, os.ModePerm)

			for _, file := range folder.Files {
				outputFile := fmt.Sprintf("%s/%s", output, file.Name)

				err = generators.GenerateFile(defaultTemplatesMap, file.Name, file.Tmpl, outputFile, data)
				if err != nil {
					return fmt.Errorf("%w", err)
				}
			}
		}

		fmt.Printf("Structure '%s' has been generated successfully\n", conf.Name)
	}

	return nil
}
