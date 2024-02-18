package lambda

import (
	_ "embed"

	"github.com/joselitofilho/aws-terraform-generator/internal/templates"
)

type SQSTrigger struct {
	SourceARN string
}

type Cron struct {
	ScheduleExpression string
	IsEnabled          string
}

type Data struct {
	Name           string
	NameSnakeCase  string
	NamePascalCase string
	Description    string
	Envars         map[string]string
	SQSTriggers    []SQSTrigger
	Crons          []Cron
	Code           map[string]templates.Code
}
