package sqs

import (
	_ "embed"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/joselitofilho/aws-terraform-generator/internal/fmtcolor"
	"github.com/joselitofilho/aws-terraform-generator/internal/generators"
	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
	generatorserrs "github.com/joselitofilho/aws-terraform-generator/internal/generators/errors"
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
		return fmt.Errorf("%w: %s", generatorserrs.ErrYAMLParser, err)
	}

	modPath := path.Join(s.output, "mod")
	_ = os.MkdirAll(modPath, os.ModePerm)

	result := make([]string, 0, len(yamlConfig.SQSs))

	templates := utils.MergeStringMap(defaultTfTemplateFiles,
		generators.CreateTemplatesMap(yamlConfig.OverrideDefaultTemplates.SQS))

	tg := generators.NewGenerator()

	for i := range yamlConfig.SQSs {
		conf := yamlConfig.SQSs[i]

		data := Data{
			Name:            conf.Name,
			MaxReceiveCount: conf.MaxReceiveCount,
		}

		if len(conf.Files) > 0 {
			filesConf := generators.CreateFilesMap(conf.Files)

			generators.MustGenerateFiles(tg, nil, filesConf, data, modPath)

			fmtcolor.White.Printf("SQS '%s' has been generated successfully\n", conf.Name)

			continue
		}

		output, err := tg.Build(data, "sqs-tf-template", templates[filenameSQStf])
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		result = append(result, output)
	}

	if len(result) > 0 {
		outputFile := path.Join(modPath, filenameSQStf)

		generators.MustGenerateFile(tg, nil, filenameSQStf, strings.Join(result, "\n"), outputFile, Data{})

		fmtcolor.White.Println("SQS has been generated successfully")
	}

	return nil
}
