package config

type OverrideDefaultTemplates struct {
	Kinesis   []FilenameTemplateMap `yaml:"kinesis,omitempty"`
	S3Buckets []FilenameTemplateMap `yaml:"bucket,omitempty"`
	SNSs      []FilenameTemplateMap `yaml:"sns,omitempty"`
	SQSs      []FilenameTemplateMap `yaml:"sqs,omitempty"`
}
