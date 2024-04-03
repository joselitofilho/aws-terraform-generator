package resourcestoyaml

import (
	"fmt"
	"strings"

	"github.com/ettle/strcase"

	"github.com/diagram-code-generator/resources/pkg/resources"

	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
	awsresources "github.com/joselitofilho/aws-terraform-generator/internal/resources"
)

func (t *Transformer) buildLambdaRelationships(source, target resources.Resource) {
	switch awsresources.ParseResourceType(source.ResourceType()) {
	case awsresources.CronType:
		t.buildCronToLambda(source, target)
	case awsresources.KinesisType:
		t.buildKinesisToLambda(source, target)
	case awsresources.SQSType:
		t.buildSQSToLambda(source, target)
	case awsresources.SNSType:
		t.buildSNSToLambda(source, target)
	}
}

func (t *Transformer) buildLambdas() (
	lambdas []config.Lambda, apiGatewayLambdasByAPIGatewayID map[string][]config.APIGatewayLambda,
) {
	apiGatewayLambdasByAPIGatewayID = map[string][]config.APIGatewayLambda{}
	apiGatewayLambdaIDs := map[string]struct{}{}

	for _, rel := range t.resc.Relationships {
		isAPIGatewayLambda := awsresources.ParseResourceType(rel.Target.ResourceType()) == awsresources.LambdaType &&
			awsresources.ParseResourceType(rel.Source.ResourceType()) == awsresources.APIGatewayType

		if isAPIGatewayLambda {
			lambda := rel.Target
			apiGatewayID := rel.Source.ID()

			apiGatewayLambdasByAPIGatewayID[apiGatewayID] = append(
				apiGatewayLambdasByAPIGatewayID[apiGatewayID], config.APIGatewayLambda{
					Name:        lambda.Value(),
					Source:      t.yamlConfig.Diagram.Lambda.Source,
					RoleName:    t.yamlConfig.Diagram.Lambda.RoleName,
					Runtime:     t.yamlConfig.Diagram.Lambda.Runtime,
					Description: fmt.Sprintf("%s lambda", lambda.Value()),
					Envars:      t.envars[lambda.ID()],
					Verb:        strings.Split(rel.Source.Value(), " ")[0],
					Path:        strings.Split(rel.Source.Value(), " ")[1],
				})

			apiGatewayLambdaIDs[lambda.ID()] = struct{}{}
		}
	}

	for _, lambda := range t.resourcesByTypeMap[awsresources.LambdaType] {
		if _, ok := apiGatewayLambdaIDs[lambda.ID()]; ok {
			continue
		}

		crons := t.buildCrons(lambda)
		kinesisTriggers := t.buildKinesisTriggers(lambda)
		sqsTriggers := t.buildSQSTriggers(lambda)

		lambdas = append(lambdas, config.Lambda{
			Name:            lambda.Value(),
			Source:          t.yamlConfig.Diagram.Lambda.Source,
			RoleName:        t.yamlConfig.Diagram.Lambda.RoleName,
			Runtime:         t.yamlConfig.Diagram.Lambda.Runtime,
			Description:     fmt.Sprintf("%s lambda", lambda.Value()),
			Envars:          t.envars[lambda.ID()],
			KinesisTriggers: kinesisTriggers,
			SQSTriggers:     sqsTriggers,
			Crons:           crons,
		})
	}

	return lambdas, apiGatewayLambdasByAPIGatewayID
}

func (t *Transformer) buildCrons(lambda resources.Resource) []config.Cron {
	var crons []config.Cron
	if cron, ok := t.cronsByLambdaID[lambda.ID()]; ok {
		crons = append(crons, config.Cron{
			ScheduleExpression: cron.Value(),
			IsEnabled:          "true",
		})
	}

	return crons
}

func (t *Transformer) buildKinesisTriggers(lambda resources.Resource) []config.KinesisTrigger {
	var kinesisTriggers []config.KinesisTrigger
	for _, kinesisTrigger := range t.kinesisTriggersByLambdaID[lambda.ID()] {
		kinesisTriggers = append(kinesisTriggers, config.KinesisTrigger{
			SourceARN: fmt.Sprintf("aws_kinesis_stream.%s_kinesis.arn", strcase.ToSnake(kinesisTrigger.Value())),
		})
	}

	return kinesisTriggers
}

func (t *Transformer) buildSQSTriggers(lambda resources.Resource) []config.SQSTrigger {
	var sqsTriggers []config.SQSTrigger
	for _, sqsTrigger := range t.sqsTriggersByLambdaID[lambda.ID()] {
		sqsTriggers = append(sqsTriggers, config.SQSTrigger{
			SourceARN: fmt.Sprintf("aws_sqs_queue.%s_sqs.arn", strcase.ToSnake(sqsTrigger.Value())),
		})
	}

	return sqsTriggers
}
