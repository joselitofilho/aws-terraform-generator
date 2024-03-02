package sqs

import (
	_ "embed"
)

//go:embed tmpls/sqs.tf.tmpl
var sqsTFTmpl []byte

var defaultTfTemplateFiles = map[string]string{
	"sqs.tf": string(sqsTFTmpl),
}
