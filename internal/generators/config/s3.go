package config

type S3 struct {
	Name           string `yaml:"name"`
	ExpirationDays int    `yaml:"expiration-days,omitempty"`
	Files          []File `yaml:"files,omitempty"`
}

func (r *S3) GetName() string { return r.Name }
