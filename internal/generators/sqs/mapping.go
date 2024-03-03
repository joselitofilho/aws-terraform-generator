package sqs

import (
	_ "embed"
)

const filenameSQStf = "sqs.tf"

//go:embed tmpls/sqs.tf.tmpl
var tmplSQStf []byte

var defaultTfTemplateFiles = map[string]string{
	filenameSQStf: string(tmplSQStf),
}
