package config

type Modules struct {
	Lambda string `yaml:"lambda"`
}

type Diagram struct {
	StackName string  `yaml:"stack_name"`
	Modules   Modules `yaml:"modules"`
}
