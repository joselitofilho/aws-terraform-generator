package kinesis

import (
	_ "embed"
)

const filenameKinesisTf = "kinesis.tf"

//go:embed tmpls/kinesis.tf.tmpl
var tmplKinesisTf []byte

var defaultTfTemplateFiles = map[string]string{
	filenameKinesisTf: string(tmplKinesisTf),
}
