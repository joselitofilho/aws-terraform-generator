package config

type S3 struct {
	Name           string `yaml:"name"`
	ExpirationDays int    `yaml:"expiration-days,omitempty"`
}
