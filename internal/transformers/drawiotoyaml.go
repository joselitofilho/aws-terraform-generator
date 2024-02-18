package transformers

import (
	"strings"

	"github.com/joselitofilho/aws-terraform-generator/internal/drawio"
	"github.com/joselitofilho/aws-terraform-generator/internal/templates/config"
)

func TransformDrawIOToYAML(stackName string, resources *drawio.ResourceCollection) (*config.Config, error) {
	endpointsPerAPIGateway := map[string]drawio.Endpoint{}
	apiGatewaysMap := map[string]drawio.APIGateway{}
	cronsMap := map[string]drawio.Cron{}

	for _, rel := range resources.Relationships {
		switch rel.Target.(type) {
		case drawio.Lambda:
			switch rel.Source.(type) {
			case drawio.Cron:
				cronsMap[rel.Target.ID()] = rel.Source.(drawio.Cron)
			}
		case drawio.APIGateway:
			apiGatewayID := rel.Target.ID()
			apiGatewaysMap[apiGatewayID] = rel.Target.(drawio.APIGateway)
			endpointsPerAPIGateway[apiGatewayID] = rel.Source.(drawio.Endpoint)
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

					if _, ok := apiGatewayLambdas[apiGatewayID]; !ok {
						apiGatewayLambdas[apiGatewayID] = []config.APIGatewayLambda{}
					}

					apiGatewayLambdas[apiGatewayID] = append(apiGatewayLambdas[apiGatewayID], config.APIGatewayLambda{
						Name:        lambda.Value(),
						Description: lambda.Description(),
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
			if cron, ok := cronsMap[lambda.ID()]; ok {
				crons = append(crons, config.Cron{
					ScheduleExpression: cron.Value(),
					IsEnabled:          "true",
				})
			}

			lambdas = append(lambdas, config.Lambda{
				Name:        lambda.Value(),
				Description: lambda.Description(),
				Code:        defaultCodes,
				Crons:       crons,
			})
		}
	}

	apiGateways := []config.APIGateway{}

	for id := range apiGatewaysMap {
		apiGateways = append(apiGateways, config.APIGateway{
			StackName: stackName,
			APIDomain: endpointsPerAPIGateway[id].Value(),
			APIG:      true,
			Lambdas:   apiGatewayLambdas[id],
		})
	}

	sqss := []config.SQS{}

	for _, sqs := range resources.SQSs {
		sqss = append(sqss, config.SQS{Name: sqs.Value(), MaxReceiveCount: 10})
	}

	return &config.Config{
		Lambdas:     lambdas,
		APIGateways: apiGateways,
		SQSs:        sqss,
	}, nil
}
