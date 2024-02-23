package config

// SNSResource represents a Lambda function or SQS configuration.
type SNSResource struct {
	Name         string   `yaml:"name"`
	Events       []string `yaml:"events"`
	FilterPrefix string   `yaml:"filter_prefix,omitempty"`
	FilterSuffix string   `yaml:"filter_suffix,omitempty"`
}

// SNS represents the configuration for SNS (Simple Notification Service).
type SNS struct {
	Name       string        `yaml:"name"`
	BucketName string        `yaml:"bucket_name"`
	Lambdas    []SNSResource `yaml:"lambdas,omitempty"`
	SQSs       []SNSResource `yaml:"sqs,omitempty"`
}
