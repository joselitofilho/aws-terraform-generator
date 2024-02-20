package config

type Lambda struct {
	ModuleLambdaSource string              `yaml:"source"`
	Name               string              `yaml:"name"`
	Description        string              `yaml:"description"`
	Envars             []map[string]string `yaml:"envars,omitempty"`
	SQSTriggers        []SQSTrigger        `yaml:"sqs-triggers,omitempty"`
	Crons              []Cron              `yaml:"crons,omitempty"`
	Code               []Code              `yaml:"code,omitempty"`
}

type SQSTrigger struct {
	SourceARN string `yaml:"source_arn"`
}

type Cron struct {
	ScheduleExpression string `yaml:"schedule_expression"`
	IsEnabled          string `yaml:"is_enabled"`
}
