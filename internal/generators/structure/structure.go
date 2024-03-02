package structure

import (
	_ "embed"
	"fmt"
	"os"
	"path"

	"github.com/joselitofilho/aws-terraform-generator/internal/generators"
	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
	generatorserrs "github.com/joselitofilho/aws-terraform-generator/internal/generators/errors"
)

type Structure struct {
	configFileName string
	output         string
}

func NewStructure(configFileName, output string) *Structure {
	return &Structure{configFileName: configFileName, output: output}
}

func (s *Structure) Build() error {
	yamlParser := config.NewYAML(s.configFileName)

	yamlConfig, err := yamlParser.Parse()
	if err != nil {
		return fmt.Errorf("%w: %w", generatorserrs.ErrYAMLParse, err)
	}

	defaultTemplatesMap := map[string]string{}

	for i := range yamlConfig.Structure.DefaultTemplates {
		for k, v := range yamlConfig.Structure.DefaultTemplates[i] {
			defaultTemplatesMap[k] = v
		}
	}

	for i := range yamlConfig.Structure.Stacks {
		conf := yamlConfig.Structure.Stacks[i]

		data := Data{
			StackName: conf.Name,
		}

		for _, folder := range conf.Folders {
			output := path.Join(s.output, conf.Name, folder.Name)
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
