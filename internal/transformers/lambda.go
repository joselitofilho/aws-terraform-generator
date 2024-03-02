package transformers

import (
	"fmt"
	"strings"

	"github.com/ettle/strcase"

	"github.com/joselitofilho/aws-terraform-generator/internal/drawio"
	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
)

func buildLambdaRelationships(
	source, target drawio.Resource, cronsByLambdaID map[string]drawio.Resource,
	kinesisTriggersByLambdaID, sqsTriggersByLambdaID map[string][]drawio.Resource, snsMap map[string]config.SNS) {
	switch source.ResourceType() {
	case drawio.CronType:
		buildCronToLambda(cronsByLambdaID, source, target)
	case drawio.KinesisType:
		buildKinesisToLambda(kinesisTriggersByLambdaID, source, target)
	case drawio.SQSType:
		buildSQSToLambda(sqsTriggersByLambdaID, source, target)
	case drawio.SNSType:
		buildSNSToLambda(snsMap, source, target)
	}
}

func buildLambdas(
	yamlConfig *config.Config,
	resourcesByTypeMap map[drawio.ResourceType][]drawio.Resource,
	resources *drawio.ResourceCollection, envars map[string]map[string]string,
	cronsByLambdaID map[string]drawio.Resource,
	kinesisTriggersByLambdaID map[string][]drawio.Resource,
	sqsTriggersByLambdaID map[string][]drawio.Resource,
) (lambdas []config.Lambda, apiGatewayLambdasByAPIGatewayID map[string][]config.APIGatewayLambda) {
	apiGatewayLambdasByAPIGatewayID = map[string][]config.APIGatewayLambda{}
	apiGatewayLambdaIDs := map[string]struct{}{}

	for _, rel := range resources.Relationships {
		isAPIGatewayLambda := rel.Target.ResourceType() == drawio.LambdaType &&
			rel.Source.ResourceType() == drawio.APIGatewayType

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

	for _, lambda := range resourcesByTypeMap[drawio.LambdaType] {
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

func buildCrons(cronsByLambdaID map[string]drawio.Resource, lambda drawio.Resource) []config.Cron {
	var crons []config.Cron
	if cron, ok := cronsByLambdaID[lambda.ID()]; ok {
		crons = append(crons, config.Cron{
			ScheduleExpression: cron.Value(),
			IsEnabled:          "true",
		})
	}

	return crons
}

func buildEnvarsList(envars map[string]map[string]string, lambda drawio.Resource) []map[string]string {
	var envarsList []map[string]string
	for key, value := range envars[lambda.ID()] {
		envarsList = append(envarsList, map[string]string{key: value})
	}

	return envarsList
}

func buildKinesisTriggers(
	kinesisTriggersByLambdaID map[string][]drawio.Resource, lambda drawio.Resource,
) []config.KinesisTrigger {
	var kinesisTriggers []config.KinesisTrigger
	for _, kinesisTrigger := range kinesisTriggersByLambdaID[lambda.ID()] {
		kinesisTriggers = append(kinesisTriggers, config.KinesisTrigger{
			SourceARN: fmt.Sprintf("aws_kinesis_stream.%s_kinesis.arn", strcase.ToSnake(kinesisTrigger.Value())),
		})
	}

	return kinesisTriggers
}

func buildSQSTriggers(sqsTriggersByLambdaID map[string][]drawio.Resource, lambda drawio.Resource) []config.SQSTrigger {
	var sqsTriggers []config.SQSTrigger
	for _, sqsTrigger := range sqsTriggersByLambdaID[lambda.ID()] {
		sqsTriggers = append(sqsTriggers, config.SQSTrigger{
			SourceARN: fmt.Sprintf("aws_sqs_queue.%s_sqs.arn", strcase.ToSnake(sqsTrigger.Value())),
		})
	}

	return sqsTriggers
}
