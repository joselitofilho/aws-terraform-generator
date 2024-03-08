package drawiotoyaml

import (
	"fmt"
	"strings"

	"github.com/ettle/strcase"

	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
	"github.com/joselitofilho/aws-terraform-generator/internal/resources"
)

func buildLambdaRelationships(
	source, target resources.Resource, cronsByLambdaID map[string]resources.Resource,
	kinesisTriggersByLambdaID, sqsTriggersByLambdaID map[string][]resources.Resource, snsMap map[string]config.SNS) {
	switch source.ResourceType() {
	case resources.CronType:
		buildCronToLambda(cronsByLambdaID, source, target)
	case resources.KinesisType:
		buildKinesisToLambda(kinesisTriggersByLambdaID, source, target)
	case resources.SQSType:
		buildSQSToLambda(sqsTriggersByLambdaID, source, target)
	case resources.SNSType:
		buildSNSToLambda(snsMap, source, target)
	}
}

func buildLambdas(
	yamlConfig *config.Config,
	resourcesByTypeMap map[resources.ResourceType][]resources.Resource,
	rscs *resources.ResourceCollection, envars map[string]map[string]string,
	cronsByLambdaID map[string]resources.Resource,
	kinesisTriggersByLambdaID map[string][]resources.Resource,
	sqsTriggersByLambdaID map[string][]resources.Resource,
) (lambdas []config.Lambda, apiGatewayLambdasByAPIGatewayID map[string][]config.APIGatewayLambda) {
	apiGatewayLambdasByAPIGatewayID = map[string][]config.APIGatewayLambda{}
	apiGatewayLambdaIDs := map[string]struct{}{}

	for _, rel := range rscs.Relationships {
		isAPIGatewayLambda := rel.Target.ResourceType() == resources.LambdaType &&
			rel.Source.ResourceType() == resources.APIGatewayType

		if isAPIGatewayLambda {
			lambda := rel.Target
			apiGatewayID := rel.Source.ID()

			envarsList := buildEnvarsList(envars, lambda)

			apiGatewayLambdasByAPIGatewayID[apiGatewayID] = append(
				apiGatewayLambdasByAPIGatewayID[apiGatewayID], config.APIGatewayLambda{
					Name:        lambda.Value(),
					Source:      yamlConfig.Diagram.Lambda.Source,
					RoleName:    yamlConfig.Diagram.Lambda.RoleName,
					Runtime:     yamlConfig.Diagram.Lambda.Runtime,
					Description: fmt.Sprintf("%s lambda", lambda.Value()),
					Envars:      envarsList,
					Verb:        strings.Split(rel.Source.Value(), " ")[0],
					Path:        strings.Split(rel.Source.Value(), " ")[1],
				})

			apiGatewayLambdaIDs[lambda.ID()] = struct{}{}
		}
	}

	for _, lambda := range resourcesByTypeMap[resources.LambdaType] {
		if _, ok := apiGatewayLambdaIDs[lambda.ID()]; ok {
			continue
		}

		crons := buildCrons(cronsByLambdaID, lambda)
		envarsList := buildEnvarsList(envars, lambda)
		kinesisTriggers := buildKinesisTriggers(kinesisTriggersByLambdaID, lambda)
		sqsTriggers := buildSQSTriggers(sqsTriggersByLambdaID, lambda)

		lambdas = append(lambdas, config.Lambda{
			Name:            lambda.Value(),
			Source:          yamlConfig.Diagram.Lambda.Source,
			RoleName:        yamlConfig.Diagram.Lambda.RoleName,
			Runtime:         yamlConfig.Diagram.Lambda.Runtime,
			Description:     fmt.Sprintf("%s lambda", lambda.Value()),
			Envars:          envarsList,
			KinesisTriggers: kinesisTriggers,
			SQSTriggers:     sqsTriggers,
			Crons:           crons,
		})
	}

	return lambdas, apiGatewayLambdasByAPIGatewayID
}

func buildCrons(cronsByLambdaID map[string]resources.Resource, lambda resources.Resource) []config.Cron {
	var crons []config.Cron
	if cron, ok := cronsByLambdaID[lambda.ID()]; ok {
		crons = append(crons, config.Cron{
			ScheduleExpression: cron.Value(),
			IsEnabled:          "true",
		})
	}

	return crons
}

func buildEnvarsList(envars map[string]map[string]string, lambda resources.Resource) []map[string]string {
	var envarsList []map[string]string
	for key, value := range envars[lambda.ID()] {
		envarsList = append(envarsList, map[string]string{key: value})
	}

	return envarsList
}

func buildKinesisTriggers(
	kinesisTriggersByLambdaID map[string][]resources.Resource, lambda resources.Resource,
) []config.KinesisTrigger {
	var kinesisTriggers []config.KinesisTrigger
	for _, kinesisTrigger := range kinesisTriggersByLambdaID[lambda.ID()] {
		kinesisTriggers = append(kinesisTriggers, config.KinesisTrigger{
			SourceARN: fmt.Sprintf("aws_kinesis_stream.%s_kinesis.arn", strcase.ToSnake(kinesisTrigger.Value())),
		})
	}

	return kinesisTriggers
}

func buildSQSTriggers(sqsTriggersByLambdaID map[string][]resources.Resource, lambda resources.Resource) []config.SQSTrigger {
	var sqsTriggers []config.SQSTrigger
	for _, sqsTrigger := range sqsTriggersByLambdaID[lambda.ID()] {
		sqsTriggers = append(sqsTriggers, config.SQSTrigger{
			SourceARN: fmt.Sprintf("aws_sqs_queue.%s_sqs.arn", strcase.ToSnake(sqsTrigger.Value())),
		})
	}

	return sqsTriggers
}
