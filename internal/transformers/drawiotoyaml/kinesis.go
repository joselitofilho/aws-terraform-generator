package drawiotoyaml

import (
	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
	"github.com/joselitofilho/aws-terraform-generator/internal/resources"
)

func buildKinesisRelationship(source, target resources.Resource, envars map[string]map[string]string) {
	if source.ResourceType() == resources.LambdaType {
		buildLambdaToKinesis(source, target, envars)
	}
}

func buildKinesis(resourcesByTypeMap map[resources.ResourceType][]resources.Resource) []config.Kinesis {
	var kinesis []config.Kinesis

	for _, k := range resourcesByTypeMap[resources.KinesisType] {
		kinesis = append(kinesis, config.Kinesis{Name: k.Value(), RetentionPeriod: "24"})
	}

	return kinesis
}
