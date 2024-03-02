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
) (lambdas []config.Lambda, apiGatewayLambdas map[string][]config.APIGatewayLambda) {
	apiGatewayLambdas = map[string][]config.APIGatewayLambda{}

	for _, lambda := range resourcesByTypeMap[drawio.LambdaType] {
		isAPIGatewayLambda := false

		for _, rel := range resources.Relationships {
			if rel.Target.ID() == lambda.ID() &&
				rel.Source.ResourceType() == drawio.APIGatewayType {
				isAPIGatewayLambda = true

				apiGatewayID := rel.Source.ID()

				var envarsList []map[string]string
				for key, value := range envars[lambda.ID()] {
					envarsList = append(envarsList, map[string]string{key: value})
				}

				apiGatewayLambdas[apiGatewayID] = append(apiGatewayLambdas[apiGatewayID], config.APIGatewayLambda{
					Name:        lambda.Value(),
					Source:      yamlConfig.Diagram.Lambda.Source,
					RoleName:    yamlConfig.Diagram.Lambda.RoleName,
					Runtime:     yamlConfig.Diagram.Lambda.Runtime,
					Description: fmt.Sprintf("%s lambda", lambda.Value()),
					Envars:      envarsList,
					Verb:        strings.Split(rel.Source.Value(), " ")[0],
					Path:        strings.Split(rel.Source.Value(), " ")[1],
				})
			}

			if isAPIGatewayLambda {
				break
			}
		}

		if !isAPIGatewayLambda {
			var crons []config.Cron
			if cron, ok := cronsByLambdaID[lambda.ID()]; ok {
				crons = append(crons, config.Cron{
					ScheduleExpression: cron.Value(),
					IsEnabled:          "true",
				})
			}

			var envarsList []map[string]string
			for key, value := range envars[lambda.ID()] {
				envarsList = append(envarsList, map[string]string{key: value})
			}

			var kinesisTriggers []config.KinesisTrigger
			for _, kinesisTrigger := range kinesisTriggersByLambdaID[lambda.ID()] {
				kinesisTriggers = append(kinesisTriggers, config.KinesisTrigger{
					SourceARN: fmt.Sprintf("aws_kinesis_stream.%s_kinesis.arn", strcase.ToSnake(kinesisTrigger.Value())),
				})
			}

			var sqsTriggers []config.SQSTrigger
			for _, sqsTrigger := range sqsTriggersByLambdaID[lambda.ID()] {
				sqsTriggers = append(sqsTriggers, config.SQSTrigger{
					SourceARN: fmt.Sprintf("aws_sqs_queue.%s_sqs.arn", strcase.ToSnake(sqsTrigger.Value())),
				})
			}

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
	}

	return lambdas, apiGatewayLambdas
}
