package s3

import (
	_ "embed"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/joselitofilho/aws-terraform-generator/internal/generators"
	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
	generatorserrs "github.com/joselitofilho/aws-terraform-generator/internal/generators/errors"
)

type Data struct {
	Name           string
	ExpirationDays int
}

type S3 struct {
	configFileName string
	output         string
}

func NewS3(configFileName, output string) *S3 {
	return &S3{configFileName: configFileName, output: output}
}

func (s *S3) Build() error {
	yamlParser := config.NewYAML(s.configFileName)

	yamlConfig, err := yamlParser.Parse()
	if err != nil {
		return fmt.Errorf("%w: %w", generatorserrs.ErrYAMLParse, err)
	}

	modPath := path.Join(s.output, "mod")
	_ = os.MkdirAll(modPath, os.ModePerm)

	tmplName := "s3-tf-template"
	result := make([]string, 0, len(yamlConfig.Buckets))

	for i := range yamlConfig.Buckets {
		conf := yamlConfig.Buckets[i]

		data := Data{
			Name:           conf.Name,
			ExpirationDays: conf.ExpirationDays,
		}

		if len(conf.Files) > 0 {
			filesConf := generators.CreateFilesMap(conf.Files)

			err = generators.GenerateFiles(nil, filesConf, data, modPath)
			if err != nil {
				return fmt.Errorf("%w", err)
			}

			fmt.Printf("S3 '%s' has been generated successfully\n", conf.Name)
		}

		output, err := generators.Build(data, tmplName, string(s3TFTmpl))
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		result = append(result, output)
	}

	if len(result) > 0 {
		outputFile := path.Join(modPath, "s3.tf")

		err := generators.GenerateFile(defaultTfTemplateFiles, tmplName, strings.Join(result, "\n"), outputFile, Data{})
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		fmt.Println("S3 has been generated successfully")
	}

	return nil
}
