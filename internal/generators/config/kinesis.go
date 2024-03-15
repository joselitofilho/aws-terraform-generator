package config

type Kinesis struct {
	Name            string `yaml:"name"`
	RetentionPeriod string `yaml:"retention_period,omitempty"`
	KMSKeyID        string `yaml:"kms_key_id,omitempty"`
	Files           []File `yaml:"files,omitempty"`
}

func (r *Kinesis) GetName() string { return r.Name }
