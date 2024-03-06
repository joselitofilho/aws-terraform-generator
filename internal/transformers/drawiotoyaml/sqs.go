package drawiotoyaml

import (
	"github.com/joselitofilho/aws-terraform-generator/internal/drawio"
	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
)

func buildSQSRelationships(
	source, target drawio.Resource, envars map[string]map[string]string, snsMap map[string]config.SNS,
) {
	switch source.ResourceType() {
	case drawio.LambdaType:
		buildLambdaToSQS(source, target, envars)
	case drawio.SNSType:
		buildSNSToSQS(snsMap, source, target)
	}
}

func buildSQSs(resourcesByTypeMap map[drawio.ResourceType][]drawio.Resource) []config.SQS {
	var sqss []config.SQS

	for _, sqs := range resourcesByTypeMap[drawio.SQSType] {
		sqss = append(sqss, config.SQS{Name: sqs.Value(), MaxReceiveCount: 10})
	}

	return sqss
}
