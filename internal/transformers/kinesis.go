package transformers

import (
	"github.com/joselitofilho/aws-terraform-generator/internal/drawio"
	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
)

func buildKinesisRelationship(source, target drawio.Resource, envars map[string]map[string]string) {
	if source.ResourceType() == drawio.LambdaType {
		buildLambdaToKinesis(source, target, envars)
	}
}

func buildKinesis(resourcesByTypeMap map[drawio.ResourceType][]drawio.Resource) []config.Kinesis {
	var kinesis []config.Kinesis

	for _, k := range resourcesByTypeMap[drawio.KinesisType] {
		kinesis = append(kinesis, config.Kinesis{Name: k.Value(), RetentionPeriod: "24"})
	}

	return kinesis
}
