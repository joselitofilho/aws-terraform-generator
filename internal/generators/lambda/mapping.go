package lambda

import (
	_ "embed"
)

const (
	filenameTfLambda = "lambda.tf"
	filenameGoLambda = "lambda.go"
	filenameGoMain   = "main.go"
)

var (
	//go:embed tmpls/lambda.tf.tmpl
	lambdaTFTmpl []byte

	//go:embed tmpls/lambda.go.tmpl
	lambdaGoTmpl []byte

	//go:embed tmpls/main.go.tmpl
	mainGoTmpl []byte
)

var (
	defaultTfTemplatesMap = map[string]string{
		filenameTfLambda: string(lambdaTFTmpl),
	}

	defaultGoTemplatesMap = map[string]string{
		filenameGoLambda: string(lambdaGoTmpl),
		filenameGoMain:   string(mainGoTmpl),
	}
)
