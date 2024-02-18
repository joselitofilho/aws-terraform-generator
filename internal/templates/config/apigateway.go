package config

type APIGatewayLambda struct {
	Name        string              `yaml:"name"`
	Description string              `yaml:"description"`
	Envars      []map[string]string `yaml:"envars"`
	Verb        string              `yaml:"verb"`
	Path        string              `yaml:"path"`
	Code        []Code              `yaml:"code"`
}

type APIGateway struct {
	StackName string             `yaml:"stack_name"`
	APIDomain string             `yaml:"api_domain"`
	APIG      bool               `yaml:"apig"`
	Lambdas   []APIGatewayLambda `yaml:"lambdas"`
}
