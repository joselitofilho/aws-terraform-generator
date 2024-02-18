package config

type Lambda struct {
	Name        string              `yaml:"name"`
	Description string              `yaml:"description"`
	Envars      []map[string]string `yaml:"envars"`
	SQSTriggers []SQSTrigger        `yaml:"sqs-triggers"`
	Cron        []Cron              `yaml:"crons"`
	Code        []Code              `yaml:"code"`
}

type SQSTrigger struct {
	SourceARN string `yaml:"source_arn"`
}

type Cron struct {
	ScheduleExpression string `yaml:"schedule_expression"`
	IsEnabled          string `yaml:"is_enabled"`
}
