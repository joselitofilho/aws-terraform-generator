package s3

import (
	_ "embed"
)

//go:embed tmpls/s3.tf.tmpl
var s3TFTmpl []byte

var defaultTfTemplateFiles = map[string]string{
	"s3.tf": string(s3TFTmpl),
}
