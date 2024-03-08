package resources_to_yaml

import (
	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
	"github.com/joselitofilho/aws-terraform-generator/internal/resources"
)

func (t *Transformer) buildKinesisRelationship(source, target resources.Resource) {
	if source.ResourceType() == resources.LambdaType {
		t.buildLambdaToKinesis(source, target)
	}
}

func (t *Transformer) buildKinesis() []config.Kinesis {
	var kinesis []config.Kinesis

	for _, k := range t.resourcesByTypeMap[resources.KinesisType] {
		kinesis = append(kinesis, config.Kinesis{Name: k.Value(), RetentionPeriod: "24"})
	}

	return kinesis
}
