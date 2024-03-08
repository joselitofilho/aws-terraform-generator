package config

type Draw struct {
	Orientation string  `yaml:"orientation,omitempty"`
	Filters     Filters `yaml:"filters,omitempty"`
}
