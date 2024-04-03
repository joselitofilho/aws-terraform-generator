package resourcestoyaml

import (
	"github.com/diagram-code-generator/resources/pkg/resources"

	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
	awsresources "github.com/joselitofilho/aws-terraform-generator/internal/resources"
)

func (t *Transformer) buildKinesisRelationship(source, target resources.Resource) {
	if awsresources.ParseResourceType(source.ResourceType()) == awsresources.LambdaType {
		t.buildLambdaToKinesis(source, target)
	}
}

func (t *Transformer) buildKinesis() []config.Kinesis {
	var kinesis []config.Kinesis

	for _, k := range t.resourcesByTypeMap[awsresources.KinesisType] {
		kinesis = append(kinesis, config.Kinesis{Name: k.Value(), RetentionPeriod: "24"})
	}

	return kinesis
}
