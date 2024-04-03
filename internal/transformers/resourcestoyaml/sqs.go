package resourcestoyaml

import (
	"github.com/diagram-code-generator/resources/pkg/resources"

	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
	awsresources "github.com/joselitofilho/aws-terraform-generator/internal/resources"
)

func (t *Transformer) buildSQSRelationships(source, target resources.Resource) {
	switch awsresources.ParseResourceType(source.ResourceType()) {
	case awsresources.LambdaType:
		t.buildLambdaToSQS(source, target)
	case awsresources.SNSType:
		t.buildSNSToSQS(source, target)
	}
}

func (t *Transformer) buildSQSs() []config.SQS {
	var sqss []config.SQS

	for _, sqs := range t.resourcesByTypeMap[awsresources.SQSType] {
		sqss = append(sqss, config.SQS{Name: sqs.Value(), MaxReceiveCount: 10})
	}

	return sqss
}
