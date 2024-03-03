package config

type OverrideDefaultTemplates struct {
	Kinesis  []FilenameTemplateMap `yaml:"kinesis,omitempty"`
	S3Bucket []FilenameTemplateMap `yaml:"bucket,omitempty"`
	SNS      []FilenameTemplateMap `yaml:"sns,omitempty"`
	SQS      []FilenameTemplateMap `yaml:"sqs,omitempty"`
}
