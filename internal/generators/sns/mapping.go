package sns

import (
	_ "embed"
)

//go:embed tmpls/sns.tf.tmpl
var snsTFTmpl []byte
