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
	tmplAPIGtf []byte

	//go:embed tmpls/lambda.go.tmpl
	tmplLambdaGo []byte

	//go:embed tmpls/lambda.tf.tmpl
	tmplLambdaTf []byte

	//go:embed tmpls/main.go.tmpl
	tmplMainGo []byte
)

var defaultGoTemplateFiles = map[string]string{
	filenameGoLambda: string(tmplLambdaGo),
	filenameGoMain:   string(tmplMainGo),
}
