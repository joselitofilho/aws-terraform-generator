package config

type Config struct {
	Structure   Structure    `yaml:"structure"`
	Lambdas     []Lambda     `yaml:"lambdas"`
	APIGateways []APIGateway `yaml:"apigateways"`
}
