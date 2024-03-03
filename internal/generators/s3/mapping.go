package s3

import (
	_ "embed"
)

const filenameS3tf = "s3.tf"

//go:embed tmpls/s3.tf.tmpl
var tmplS3tf []byte

var defaultTfTemplateFiles = map[string]string{
	filenameS3tf: string(tmplS3tf),
}
