package lambda

import (
	_ "embed"

	templates "github.com/joselitofilho/aws-terraform-generator/internal/generators"
)

type SQSTrigger struct {
	SourceARN string
}

type Cron struct {
	ScheduleExpression string
	IsEnabled          string
}

type Data struct {
	ModuleLambdaSource string
	Name               string
	NameSnakeCase      string
	NamePascalCase     string
	Description        string
	Envars             map[string]string
	SQSTriggers        []SQSTrigger
	Crons              []Cron
	Files              map[string]templates.File
}