package config

type OverrideDefaultTemplates struct {
	SNSs []FilenameTemplateMap `yaml:"sns,omitempty"`
	SQSs []FilenameTemplateMap `yaml:"sqs,omitempty"`
}
