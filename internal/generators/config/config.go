package config

type Config struct {
	Diagram     Diagram      `yaml:"diagram,omitempty"`
	Structure   Structure    `yaml:"structure,omitempty"`
	Lambdas     []Lambda     `yaml:"lambdas,omitempty"`
	APIGateways []APIGateway `yaml:"apigateways,omitempty"`
	SQSs        []SQS        `yaml:"sqs,omitempty"`
	SNSs        []SNS        `yaml:"sns,omitempty"`
	Buckets     []S3         `yaml:"buckets,omitempty"`
	RestfulAPIs []RestfulAPI `yaml:"restfulapis,omitempty"`
}
