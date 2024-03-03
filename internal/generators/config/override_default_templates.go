package config

type OverrideDefaultTemplates struct {
	S3Buckets []FilenameTemplateMap `yaml:"bucket,omitempty"`
	SNSs      []FilenameTemplateMap `yaml:"sns,omitempty"`
	SQSs      []FilenameTemplateMap `yaml:"sqs,omitempty"`
}
