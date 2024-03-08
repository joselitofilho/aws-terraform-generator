package resources_to_yaml

import (
	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
	"github.com/joselitofilho/aws-terraform-generator/internal/resources"
)

func buildSQSRelationships(
	source, target resources.Resource, envars map[string]map[string]string, snsMap map[string]config.SNS,
) {
	switch source.ResourceType() {
	case resources.LambdaType:
		buildLambdaToSQS(source, target, envars)
	case resources.SNSType:
		buildSNSToSQS(snsMap, source, target)
	}
}

func buildSQSs(resourcesByTypeMap map[resources.ResourceType][]resources.Resource) []config.SQS {
	var sqss []config.SQS

	for _, sqs := range resourcesByTypeMap[resources.SQSType] {
		sqss = append(sqss, config.SQS{Name: sqs.Value(), MaxReceiveCount: 10})
	}

	return sqss
}
