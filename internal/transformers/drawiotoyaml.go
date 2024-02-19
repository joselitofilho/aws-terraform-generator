package transformers

import (
	"fmt"
	"strings"

	"github.com/ettle/strcase"
	"github.com/joselitofilho/aws-terraform-generator/internal/drawio"
	"github.com/joselitofilho/aws-terraform-generator/internal/templates/config"
)

func TransformDrawIOToYAML(stackName string, resources *drawio.ResourceCollection) (*config.Config, error) {
	endpointsByAPIGatewayID := map[string]drawio.Resource{}
	apiGatewaysByID := map[string]drawio.Resource{}
	cronsByLambdaID := map[string]drawio.Resource{}
	envars := map[string]map[string]string{}
	sqsTriggersByLambdaID := map[string][]drawio.Resource{}

	for _, rel := range resources.Relationships {
		switch rel.Target.(type) {
		case drawio.Lambda:
			lambdaID := rel.Target.ID()

			switch rel.Source.(type) {
			case drawio.Cron:
				cronsByLambdaID[lambdaID] = rel.Source
			case drawio.SQS:
				sqsTriggersByLambdaID[lambdaID] = append(sqsTriggersByLambdaID[lambdaID], rel.Source)
			}
		case drawio.APIGateway:
			apiGatewayID := rel.Target.ID()
			apiGatewaysByID[apiGatewayID] = rel.Target
			endpointsByAPIGatewayID[apiGatewayID] = rel.Source
		case drawio.SQS:
			switch rel.Source.(type) {
			case drawio.Lambda:
				if _, ok := envars[rel.Source.ID()]; !ok {
					envars[rel.Source.ID()] = map[string]string{}
				}

				envars[rel.Source.ID()]["SQS_QUEUE_URL"] =
					fmt.Sprintf("aws_sqs_queue.%s_sqs.id", strcase.ToSnake(rel.Target.Value()))
			}
		case drawio.Database:
			switch rel.Source.(type) {
			case drawio.Lambda:
				if _, ok := envars[rel.Source.ID()]; !ok {
					envars[rel.Source.ID()] = map[string]string{}
				}

				envars[rel.Source.ID()]["DOCDB_HOST"] = "var.docdb_host"
				envars[rel.Source.ID()]["DOCDB_USER"] = "var.docdb_user"
				envars[rel.Source.ID()]["DOCDB_PASSWORD_SECRET"] = "var.docdb_password_secret"
			}
		}
	}

	defaultCodes := []config.Code{{Key: "lambda"}, {Key: "main"}}
	lambdas := []config.Lambda{}
	apiGatewayLambdas := map[string][]config.APIGatewayLambda{}

	for i := range resources.Lambdas {
		lambda := resources.Lambdas[i]

		isAPIGatewayLambda := false

		for _, rel := range resources.Relationships {
			if rel.Target.ID() == lambda.ID() {
				switch rel.Source.(type) {
				case drawio.APIGateway:
					isAPIGatewayLambda = true

					apiGatewayID := rel.Source.ID()

					envarsList := []map[string]string{}
					for key, value := range envars[lambda.ID()] {
						envarsList = append(envarsList, map[string]string{key: value})
					}

					apiGatewayLambdas[apiGatewayID] = append(apiGatewayLambdas[apiGatewayID], config.APIGatewayLambda{
						Name:        lambda.Value(),
						Description: fmt.Sprintf("%s lambda", lambda.Value()),
						Envars:      envarsList,
						Verb:        strings.Split(rel.Source.Value(), " ")[0],
						Path:        strings.Split(rel.Source.Value(), " ")[1],
						Code:        defaultCodes,
					})
				}
			}

			if isAPIGatewayLambda {
				break
			}
		}

		if !isAPIGatewayLambda {
			crons := []config.Cron{}
			if cron, ok := cronsByLambdaID[lambda.ID()]; ok {
				crons = append(crons, config.Cron{
					ScheduleExpression: cron.Value(),
					IsEnabled:          "true",
				})
			}

			envarsList := []map[string]string{}
			for key, value := range envars[lambda.ID()] {
				envarsList = append(envarsList, map[string]string{key: value})
			}

			sqsTriggers := []config.SQSTrigger{}
			for _, sqsTrigger := range sqsTriggersByLambdaID[lambda.ID()] {
				sqsTriggers = append(sqsTriggers, config.SQSTrigger{
					SourceARN: fmt.Sprintf("aws_sqs_queue.%s_sqs.arn", strcase.ToSnake(sqsTrigger.Value())),
				})
			}

			lambdas = append(lambdas, config.Lambda{
				Name:        lambda.Value(),
				Description: fmt.Sprintf("%s lambda", lambda.Value()),
				Envars:      envarsList,
				SQSTriggers: sqsTriggers,
				Code:        defaultCodes,
				Crons:       crons,
			})
		}
	}

	apiGateways := []config.APIGateway{}

	for id := range apiGatewaysByID {
		apiGateways = append(apiGateways, config.APIGateway{
			StackName: stackName,
			APIDomain: endpointsByAPIGatewayID[id].Value(),
			APIG:      true,
			Lambdas:   apiGatewayLambdas[id],
		})
	}

	sqss := []config.SQS{}

	for _, sqs := range resources.SQSs {
		sqss = append(sqss, config.SQS{Name: sqs.Value(), MaxReceiveCount: 10})
	}

	buckets := []config.S3{}

	for _, bucket := range resources.Buckets {
		buckets = append(buckets, config.S3{Name: bucket.Value()})
	}

	return &config.Config{
		Lambdas:     lambdas,
		APIGateways: apiGateways,
		SQSs:        sqss,
		Buckets:     buckets,
	}, nil
}
