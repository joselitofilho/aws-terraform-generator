package sqs

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/ettle/strcase"

	"github.com/joselitofilho/aws-terraform-generator/internal/templates"
	"github.com/joselitofilho/aws-terraform-generator/internal/templates/config"
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
	input  string
	output string
}

func NewSQS(input, output string) *SQS {
	return &SQS{input: input, output: output}
}

func (s *SQS) Build() error {
	yamlParser := config.NewYAML(s.input)

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

		output, err := templates.Build(data, tmplName, string(sqsTFTmpl))
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		result = fmt.Sprintf("%s\n%s", result, output)
	}

	if result != "" {
		err = templates.BuildFile(Data{}, tmplName, result, s.output)
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		err = utils.TerraformFormat(s.output)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("SQS has been generated successfully")
	}

	return nil
}
