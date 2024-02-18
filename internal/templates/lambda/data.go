package lambda

import (
	_ "embed"
)

var (
	//go:embed tmpls/lambda.tf.tmpl
	lambdaTFTmpl []byte

	//go:embed tmpls/lambda.go.tmpl
	lambdaGoTmpl []byte

	//go:embed tmpls/main.go.tmpl
	mainGoTmpl []byte
)

type SQSTrigger struct {
	SourceARN string
}

type Cron struct {
	ScheduleExpression string
	IsEnabled          string
}

type Code struct {
	Tmpl    string
	Imports []string
}

type Data struct {
	Name           string
	NameSnakeCase  string
	NamePascalCase string
	Description    string
	Envars         map[string]string
	SQSTriggers    []SQSTrigger
	Crons          []Cron
	Code           map[string]Code
}
