package config

type SQS struct {
	Name            string `yaml:"name"`
	MaxReceiveCount int32  `yaml:"max_receive_count"`
}
