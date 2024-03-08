package resourcestoyaml

import (
	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
	"github.com/joselitofilho/aws-terraform-generator/internal/resources"
)

func (t *Transformer) buildSQSRelationships(source, target resources.Resource) {
	switch source.ResourceType() {
	case resources.LambdaType:
		t.buildLambdaToSQS(source, target)
	case resources.SNSType:
		t.buildSNSToSQS(source, target)
	}
}

func (t *Transformer) buildSQSs() []config.SQS {
	var sqss []config.SQS

	for _, sqs := range t.resourcesByTypeMap[resources.SQSType] {
		sqss = append(sqss, config.SQS{Name: sqs.Value(), MaxReceiveCount: 10})
	}

	return sqss
}
