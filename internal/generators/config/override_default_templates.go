package config

type OverrideDefaultTemplates struct {
	SQSs []FilenameTemplateMap `yaml:"sqs,omitempty"`
}
