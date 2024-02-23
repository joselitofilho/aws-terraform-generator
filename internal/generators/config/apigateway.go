package config

type APIGatewayLambda struct {
	Source      string              `yaml:"source"`
	Name        string              `yaml:"name"`
	Description string              `yaml:"description"`
	Envars      []map[string]string `yaml:"envars,omitempty"`
	Verb        string              `yaml:"verb"`
	Path        string              `yaml:"path"`
	Files       []File              `yaml:"files,omitempty"`
}

type APIGateway struct {
	StackName string             `yaml:"stack_name"`
	APIDomain string             `yaml:"api_domain"`
	APIG      bool               `yaml:"apig"`
	Lambdas   []APIGatewayLambda `yaml:"lambdas"`
}