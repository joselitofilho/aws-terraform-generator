package sqs

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/ettle/strcase"

	"github.com/joselitofilho/aws-terraform-generator/internal/templates"
)

//go:embed tmpls/sqs.tf.tmpl
var templateContent []byte

type Data struct {
	Name          string
	NameWithSpace string
	NameSnakeCase string
}

type SQS struct {
	name   string
	output string
}

func NewSQS(name, output string) *SQS {
	return &SQS{name: name, output: output}
}

func (s *SQS) Build() error {
	data := Data{
		Name:          s.name,
		NameWithSpace: strings.ReplaceAll(s.name, "-", " "),
		NameSnakeCase: strcase.ToSnake(s.name),
	}

	if s.output == "" {
		output, err := templates.Build(data, string(templateContent))
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		fmt.Println("SQS output:")
		fmt.Println(output)

		return nil
	}

	err := templates.BuildFile(data, string(templateContent), s.output)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	fmt.Println("SQS has been generated successfully")

	return nil
}
