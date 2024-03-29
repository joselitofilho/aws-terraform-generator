package lambda

import (
	_ "embed"

	"github.com/joselitofilho/aws-terraform-generator/internal/generators"
)

type KinesisTrigger struct {
	SourceARN string
}

type SQSTrigger struct {
	SourceARN string
}

type Cron struct {
	ScheduleExpression string
	IsEnabled          string
}

type Data struct {
	Name            string
	AsModule        bool
	Source          string
	RoleName        string
	Runtime         string
	Description     string
	Envars          map[string]string
	KinesisTriggers []KinesisTrigger
	SQSTriggers     []SQSTrigger
	Crons           []Cron
	Files           map[string]generators.File
}
