package sqs

import (
	_ "embed"
	"fmt"

	"github.com/joselitofilho/aws-terraform-generator/internal/generators"
	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
	"github.com/joselitofilho/aws-terraform-generator/internal/utils"
)

type Data struct {
	Name            string
	MaxReceiveCount int32
}

type SQS struct {
	configFileName string
	output         string
}

func NewSQS(configFileName, output string) *SQS {
	return &SQS{configFileName: configFileName, output: output}
}

func (s *SQS) Build() error {
	yamlParser := config.NewYAML(s.configFileName)

	yamlConfig, err := yamlParser.Parse()
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	tmplName := "sqs-tf-template"
	result := ""

	for i := range yamlConfig.SQSs {
		conf := yamlConfig.SQSs[i]

		data := Data{
			Name:            conf.Name,
			MaxReceiveCount: conf.MaxReceiveCount,
		}

		if len(conf.Files) > 0 {
			filesConf := generators.CreateFilesMap(conf.Files)

			err = generators.GenerateFiles(defaultTemplatesMap, filesConf, fmt.Sprintf("%s/mod", s.output), data)
			if err != nil {
				return fmt.Errorf("%w", err)
			}

			fmt.Printf("SQS '%s' has been generated successfully\n", conf.Name)

			continue
		}

		output, err := generators.Build(data, tmplName, string(sqsTFTmpl))
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		result = fmt.Sprintf("%s\n%s", result, output)
	}

	if result != "" {
		outputFile := fmt.Sprintf("%s/mod/sqs.tf", s.output)

		if err := generators.BuildFile(Data{}, tmplName, result, outputFile); err != nil {
			return fmt.Errorf("%w", err)
		}

		if err := utils.TerraformFormat(outputFile); err != nil {
			fmt.Println(err)
		}

		fmt.Println("SQS has been generated successfully")
	}

	return nil
}
