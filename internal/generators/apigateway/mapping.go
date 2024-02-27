package apigateway

import (
	_ "embed"
)

const (
	filenameTfAPIG   = "apig.tf"
	filenameTfLambda = "lambda.tf"
	filenameGoLambda = "lambda.go"
	filenameGoMain   = "main.go"
)

var (
	//go:embed tmpls/apig.tf.tmpl
	apigTFTmpl []byte

	//go:embed tmpls/lambda.go.tmpl
	lambdaGoTmpl []byte

	//go:embed tmpls/lambda.tf.tmpl
	lambdaTFTmpl []byte

	//go:embed tmpls/main.go.tmpl
	mainGoTmpl []byte
)

var defaultTfTemplateFiles = map[string]string{
	filenameTfAPIG:   string(apigTFTmpl),
	filenameTfLambda: string(lambdaTFTmpl),
}

var defaultGoTemplateFiles = map[string]string{
	filenameGoLambda: string(lambdaGoTmpl),
	filenameGoMain:   string(mainGoTmpl),
}
