package sqs

import (
	_ "embed"

	"github.com/joselitofilho/aws-terraform-generator/internal/generators"
)

//go:embed tmpls/sqs.tf.tmpl
var sqsTFTmpl []byte

var defaultTemplatesMap = map[string]generators.TemplateMapValue{
	"sqs.tf": {TemplateName: "sqs-tf-template", Template: sqsTFTmpl},
}
