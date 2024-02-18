package lambda

import (
	_ "embed"

	"github.com/joselitofilho/aws-terraform-generator/internal/templates"
)

type templateMapValue struct {
	templateName string
	template     []byte
}

var (
	//go:embed tmpls/lambda.tf.tmpl
	lambdaTFTmpl []byte

	//go:embed tmpls/lambda.go.tmpl
	lambdaGoTmpl []byte

	//go:embed tmpls/main.go.tmpl
	mainGoTmpl []byte
)

var defaultTemplatesMap = map[string]templates.TemplateMapValue{
	"main":   {TemplateName: "main-go-template", Template: mainGoTmpl},
	"lambda": {TemplateName: "lambda-go-template", Template: lambdaGoTmpl},
}
