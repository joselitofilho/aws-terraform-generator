package sqs

import (
	_ "embed"
)

//go:embed tmpls/sqs.tf.tmpl
var sqsTFTmpl []byte
