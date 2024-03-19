package config

type Lambda struct {
	Name            string            `yaml:"name"`
	Source          string            `yaml:"source"`
	RoleName        string            `yaml:"role_name,omitempty"`
	Runtime         string            `yaml:"runtime,omitempty"`
	Description     string            `yaml:"description"`
	Envars          map[string]string `yaml:"envars,omitempty"`
	KinesisTriggers []KinesisTrigger  `yaml:"kinesis-triggers,omitempty"`
	SQSTriggers     []SQSTrigger      `yaml:"sqs-triggers,omitempty"`
	Crons           []Cron            `yaml:"crons,omitempty"`
	Files           []File            `yaml:"files,omitempty"`
}

func (r *Lambda) GetName() string { return r.Name }

type SQSTrigger struct {
	SourceARN string `yaml:"source_arn"`
}

type Cron struct {
	ScheduleExpression string `yaml:"schedule_expression"`
	IsEnabled          string `yaml:"is_enabled"`
}

type KinesisTrigger struct {
	SourceARN string `yaml:"source_arn"`
}
