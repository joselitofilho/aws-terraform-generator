package config

// Config represents a configuration object that can be populated from a YAML file
type Config struct {
	Diagram     Diagram      `yaml:"diagram,omitempty"`
	Structure   Structure    `yaml:"structure,omitempty"`
	Lambdas     []Lambda     `yaml:"lambdas,omitempty"`
	APIGateways []APIGateway `yaml:"apigateways,omitempty"`
	SNSs        []SNS        `yaml:"sns,omitempty"`
	SQSs        []SQS        `yaml:"sqs,omitempty"`
	Buckets     []S3         `yaml:"buckets,omitempty"`
	RestfulAPIs []RestfulAPI `yaml:"restfulapis,omitempty"`
}
