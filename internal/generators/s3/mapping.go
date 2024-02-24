package s3

import (
	_ "embed"

	"github.com/joselitofilho/aws-terraform-generator/internal/generators"
)

//go:embed tmpls/s3.tf.tmpl
var s3TFTmpl []byte

var defaultTemplatesMap = map[string]generators.TemplateMapValue{
	"s3.tf": {TemplateName: "s3-tf-template", Template: s3TFTmpl},
}
