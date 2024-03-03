package sns

import (
	_ "embed"
)

const filenameSNStf = "sns.tf"

//go:embed tmpls/sns.tf.tmpl
var tmplSNStf []byte

var defaultTfTemplateFiles = map[string]string{
	filenameSNStf: string(tmplSNStf),
}
