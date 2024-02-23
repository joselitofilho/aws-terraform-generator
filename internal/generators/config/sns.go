package config

// SNSLambda represents a Lambda function configuration.
type SNSLambda struct {
	Name         string   `yaml:"name"`
	Events       []string `yaml:"events"`
	FilterPrefix string   `yaml:"filter_prefix"`
	FilterSuffix string   `yaml:"filter_suffix"`
}

// SNS represents the configuration for SNS (Simple Notification Service).
type SNS struct {
	Name       string      `yaml:"name"`
	BucketName string      `yaml:"bucket_name"`
	Lambdas    []SNSLambda `yaml:"lambda"`
}
