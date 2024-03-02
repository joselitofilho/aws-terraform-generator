package kinesis

import (
	_ "embed"
)

//go:embed tmpls/kinesis.tf.tmpl
var kinesisTFTmpl []byte

var defaultTfTemplateFiles = map[string]string{
	"kinesis.tf": string(kinesisTFTmpl),
}