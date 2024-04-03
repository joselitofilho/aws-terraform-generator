package yamltoresources

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/diagram-code-generator/resources/pkg/resources"
	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
)

func Parse(filename string) (*resources.ResourceCollection, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	var cfg *config.Config

	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	result, err := NewTransformer(cfg).Transform()
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return result, nil
}
