package yamltoresources

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
	"github.com/joselitofilho/aws-terraform-generator/internal/resources"
)

func Parse(filename string) (*resources.ResourceCollection, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	var config *config.Config

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	result, err := NewTransformer(config).Transform()
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return result, nil
}
