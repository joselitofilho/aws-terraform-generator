package config

type Modules struct {
	Lambda string `yaml:"lambda"`
}

type Diagram struct {
	Modules Modules `yaml:"modules"`
}
