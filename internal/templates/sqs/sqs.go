package sqs

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/ettle/strcase"

	"github.com/joselitofilho/aws-terraform-generator/internal/templates"
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
	name            string
	maxReceiveCount int32
	output          string
}

func NewSQS(name string, maxReceiveCount int32, output string) *SQS {
	return &SQS{name: name, maxReceiveCount: maxReceiveCount, output: output}
}

func (s *SQS) Build() error {
	data := Data{
		Name:            s.name,
		NameWithSpace:   strings.ReplaceAll(s.name, "-", " "),
		NameSnakeCase:   strcase.ToSnake(s.name),
		MaxReceiveCount: s.maxReceiveCount,
	}

	tmplName := "sqs-tf-template"

	if s.output == "" {
		output, err := templates.Build(data, tmplName, string(sqsTFTmpl))
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		fmt.Println("SQS output:")
		fmt.Println(output)

		return nil
	}

	err := templates.BuildFile(data, tmplName, string(sqsTFTmpl), s.output)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	err = utils.TerraformFormat(s.output)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("SQS has been generated successfully")

	return nil
}
