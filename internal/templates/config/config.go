package config

type Config struct {
	Structure   Structure    `yaml:"structure,omitempty"`
	Lambdas     []Lambda     `yaml:"lambdas,omitempty"`
	APIGateways []APIGateway `yaml:"apigateways,omitempty"`
	SQSs        []SQS        `yaml:"sqs,omitempty"`
	Buckets     []S3         `yaml:"buckets,omitempty"`
}
