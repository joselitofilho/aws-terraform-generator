package config

type DriagramLambda struct {
	Source   string `yaml:"source"`
	RoleName string `yaml:"role_name"`
	Runtime  string `yaml:"runtime,omitempty"`
}

type Diagram struct {
	StackName string         `yaml:"stack_name"`
	Lambda    DriagramLambda `yaml:"lambda"`
}
