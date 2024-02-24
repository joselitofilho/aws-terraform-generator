package sns

import (
	_ "embed"
)

//go:embed tmpls/sns.tf.tmpl
var snsTFTmpl []byte

var defaultTemplatesMap = map[string]templates.TemplateMapValue{
	"sns.tf": {TemplateName: "sns-tf-template", Template: snsTFTmpl},
}
