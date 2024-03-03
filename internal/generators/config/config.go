package config

// Config represents a configuration object that can be populated from a YAML file.
type Config struct {
	OverrideDefaultTemplates OverrideDefaultTemplates `yaml:"override_default_templates,omitempty"`
	Diagram                  Diagram                  `yaml:"diagram,omitempty"`
	Structure                Structure                `yaml:"structure,omitempty"`
	APIGateways              []APIGateway             `yaml:"apigateways,omitempty"`
	Kinesis                  []Kinesis                `yaml:"kinesis,omitempty"`
	Lambdas                  []Lambda                 `yaml:"lambdas,omitempty"`
	Buckets                  []S3                     `yaml:"buckets,omitempty"`
	SNSs                     []SNS                    `yaml:"sns,omitempty"`
	SQSs                     []SQS                    `yaml:"sqs,omitempty"`
	RestfulAPIs              []RestfulAPI             `yaml:"restfulapis,omitempty"`
}
