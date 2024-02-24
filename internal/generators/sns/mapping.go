package sns

import (
	_ "embed"

	"github.com/joselitofilho/aws-terraform-generator/internal/generators"
)

//go:embed tmpls/sns.tf.tmpl
var snsTFTmpl []byte

var defaultTemplatesMap = map[string]generators.TemplateMapValue{
	"sns.tf": {TemplateName: "sns-tf-template", Template: snsTFTmpl},
}
