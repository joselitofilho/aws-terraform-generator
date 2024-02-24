package sqs

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/ettle/strcase"

	"github.com/joselitofilho/aws-terraform-generator/internal/generators"
	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
	"github.com/joselitofilho/aws-terraform-generator/internal/utils"
)

//go:embed tmpls/sqs.tf.tmpl
var sqsTFTmpl []byte

type Data struct {
	Name            string
	NameWithSpace   string
	NameSnakeCase   string
	MaxReceiveCount int32
}

type SQS struct {
	config string
	output string
}

func NewSQS(config, output string) *SQS {
	return &SQS{config: config, output: output}
}

func (s *SQS) Build() error {
	yamlParser := config.NewYAML(s.config)

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
			NameWithSpace:   strings.ReplaceAll(conf.Name, "-", " "),
			NameSnakeCase:   strcase.ToSnake(conf.Name),
			MaxReceiveCount: conf.MaxReceiveCount,
		}

		output, err := generators.Build(data, tmplName, string(sqsTFTmpl))
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		result = fmt.Sprintf("%s\n%s", result, output)
	}

	if result != "" {
		outputFile := fmt.Sprintf("%s/sqs.tf", s.output)

		err = generators.BuildFile(Data{}, tmplName, result, outputFile)
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		err = utils.TerraformFormat(outputFile)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("SQS has been generated successfully")
	}

	return nil
}
