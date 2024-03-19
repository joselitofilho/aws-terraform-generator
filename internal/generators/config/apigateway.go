package config

type APIGatewayLambda struct {
	Name        string            `yaml:"name"`
	Source      string            `yaml:"source"`
	RoleName    string            `yaml:"role_name,omitempty"`
	Runtime     string            `yaml:"runtime,omitempty"`
	Description string            `yaml:"description"`
	Envars      map[string]string `yaml:"envars,omitempty"`
	Verb        string            `yaml:"verb"`
	Path        string            `yaml:"path"`
	Files       []File            `yaml:"files,omitempty"`
}

func (r *APIGatewayLambda) GetName() string { return r.Name }

type APIGateway struct {
	StackName string             `yaml:"stack_name"`
	APIDomain string             `yaml:"api_domain"`
	APIG      bool               `yaml:"apig"`
	Lambdas   []APIGatewayLambda `yaml:"lambdas"`
}
