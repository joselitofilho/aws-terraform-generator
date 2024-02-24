package lambda

import (
	_ "embed"

	"github.com/joselitofilho/aws-terraform-generator/internal/generators"
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
	Description        string
	Envars             map[string]string
	SQSTriggers        []SQSTrigger
	Crons              []Cron
	Files              map[string]generators.File
}
