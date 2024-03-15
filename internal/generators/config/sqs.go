package config

type SQS struct {
	Name            string `yaml:"name"`
	MaxReceiveCount int32  `yaml:"max_receive_count"`
	Files           []File `yaml:"files,omitempty"`
}

func (r *SQS) GetName() string { return r.Name }
