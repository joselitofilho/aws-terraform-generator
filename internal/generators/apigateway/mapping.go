package apigateway

import (
	_ "embed"

	"github.com/joselitofilho/aws-terraform-generator/internal/generators"
)

var (
	//go:embed tmpls/apig.tf.tmpl
	apigTFTmpl []byte

	//go:embed tmpls/lambda.go.tmpl
	lambdaGoTmpl []byte

	//go:embed tmpls/lambda.tf.tmpl
	lambdaTFTmpl []byte

	//go:embed tmpls/main.go.tmpl
	mainGoTmpl []byte
)

var defaultTemplatesMap = map[string]generators.TemplateMapValue{
	"apig.tf":   {TemplateName: "apig-tf-template", Template: apigTFTmpl},
	"lambda.go": {TemplateName: "lambda-go-template", Template: lambdaGoTmpl},
	"lambda.tf": {TemplateName: "lambda-tf-template", Template: lambdaTFTmpl},
	"main.go":   {TemplateName: "main-go-template", Template: mainGoTmpl},
}
