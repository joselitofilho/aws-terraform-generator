package transformers

import (
	"fmt"
	"strings"

	"github.com/ettle/strcase"
	"github.com/joselitofilho/aws-terraform-generator/internal/drawio"
	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
)

func TransformDrawIOToYAML(yamlConfig *config.Config, resources *drawio.ResourceCollection) (*config.Config, error) {
	apiGatewaysByID := map[string]drawio.Resource{}
	endpointsByAPIGatewayID := map[string]drawio.Resource{}
	cronsByLambdaID := map[string]drawio.Resource{}
	sqsTriggersByLambdaID := map[string][]drawio.Resource{}
	envars := map[string]map[string]string{}
	snsMap := map[string]config.SNS{}

	buildResourceRelationships(
		resources, cronsByLambdaID, sqsTriggersByLambdaID, snsMap, apiGatewaysByID, endpointsByAPIGatewayID, envars)

	resourcesByTypeMap := buildResourcesByTypeMap(resources)

	lambdas, apiGatewayLambdas := buildLambdas(
		resourcesByTypeMap, resources, envars, yamlConfig, cronsByLambdaID, sqsTriggersByLambdaID)
	apiGateways := buildAPIGateways(yamlConfig, apiGatewaysByID, endpointsByAPIGatewayID, apiGatewayLambdas)
	snss := buildSNSs(snsMap)
	sqss := buildSQSs(resourcesByTypeMap)
	buckets := buildBuckets(resourcesByTypeMap)
	restfulAPIs := buildRestfulAPIs(resourcesByTypeMap)

	return &config.Config{
		Lambdas:     lambdas,
		APIGateways: apiGateways,
		SNSs:        snss,
		SQSs:        sqss,
		Buckets:     buckets,
		RestfulAPIs: restfulAPIs,
	}, nil
}

func buildResourcesByTypeMap(resources *drawio.ResourceCollection) map[drawio.ResourceType][]drawio.Resource {
	resourcesByTypeMap := map[drawio.ResourceType][]drawio.Resource{}

	for _, resource := range resources.Resources {
		resourcesByTypeMap[resource.ReseourceType()] = append(resourcesByTypeMap[resource.ReseourceType()], resource)
	}

	return resourcesByTypeMap
}

//nolint:gocyclo // Reducing complexity will make it unreadable
func buildResourceRelationships(
	resources *drawio.ResourceCollection,
	cronsByLambdaID map[string]drawio.Resource,
	sqsTriggersByLambdaID map[string][]drawio.Resource,
	snsMap map[string]config.SNS,
	apiGatewaysByID map[string]drawio.Resource,
	endpointsByAPIGatewayID map[string]drawio.Resource,
	envars map[string]map[string]string,
) {
	for _, rel := range resources.Relationships {
		target := rel.Target
		source := rel.Source

		switch target.ReseourceType() {
		case drawio.LambdaType:
			switch source.ReseourceType() {
			case drawio.CronType:
				buildCronToLambda(cronsByLambdaID, source, target)
			case drawio.SQSType:
				buildSQSToLambda(sqsTriggersByLambdaID, source, target)
			case drawio.SNSType:
				buildSNSToLambda(snsMap, source)
			}
		case drawio.APIGatewayType:
			if source.ReseourceType() == drawio.EndpointType {
				buildEndpointToAPIGateway(apiGatewaysByID, endpointsByAPIGatewayID, source, target)
			}
		case drawio.SQSType:
			switch source.ReseourceType() {
			case drawio.LambdaType:
				buildLambdaToSQS(envars, source, target)
			case drawio.SNSType:
				buildSNSToSQS(snsMap, source)
			}
		case drawio.DatabaseType:
			if source.ReseourceType() == drawio.LambdaType {
				buildLambdaToDatabase(envars, source)
			}
		case drawio.RestfulAPIType:
			if source.ReseourceType() == drawio.LambdaType {
				buildLambdaToRestfulAPI(envars, source, target)
			}
		case drawio.S3Type:
			if source.ReseourceType() == drawio.LambdaType {
				buildLambdaToS3(envars, source, target)
			}
		case drawio.SNSType:
			if source.ReseourceType() == drawio.S3Type {
				buildS3ToSNS(snsMap, source, target)
			}
		}
	}
}

func buildLambdas(
	resourcesByTypeMap map[drawio.ResourceType][]drawio.Resource,
	resources *drawio.ResourceCollection, envars map[string]map[string]string,
	yamlConfig *config.Config,
	cronsByLambdaID map[string]drawio.Resource,
	sqsTriggersByLambdaID map[string][]drawio.Resource,
) (lambdas []config.Lambda, apiGatewayLambdas map[string][]config.APIGatewayLambda) {
	apiGatewayLambdas = map[string][]config.APIGatewayLambda{}
	defaultFiles := []config.File{{Name: "lambda.go"}, {Name: "main.go"}}

	for _, lambda := range resourcesByTypeMap[drawio.LambdaType] {
		isAPIGatewayLambda := false

		for _, rel := range resources.Relationships {
			if rel.Target.ID() == lambda.ID() {
				if rel.Source.ReseourceType() == drawio.APIGatewayType {
					isAPIGatewayLambda = true

					apiGatewayID := rel.Source.ID()

					envarsList := []map[string]string{}
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
						Files:       defaultFiles,
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
				Source:      yamlConfig.Diagram.Lambda.Source,
				RoleName:    yamlConfig.Diagram.Lambda.RoleName,
				Runtime:     yamlConfig.Diagram.Lambda.Runtime,
				Description: fmt.Sprintf("%s lambda", lambda.Value()),
				Envars:      envarsList,
				SQSTriggers: sqsTriggers,
				Files:       defaultFiles,
				Crons:       crons,
			})
		}
	}

	return lambdas, apiGatewayLambdas
}

func buildAPIGateways(
	yamlConfig *config.Config,
	apiGatewaysByID map[string]drawio.Resource,
	endpointsByAPIGatewayID map[string]drawio.Resource,
	apiGatewayLambdas map[string][]config.APIGatewayLambda,
) []config.APIGateway {
	apiGateways := []config.APIGateway{
		{
			StackName: yamlConfig.Diagram.StackName,
			APIG:      true,
		},
	}

	for id := range apiGatewaysByID {
		apiGateways[0].APIDomain = endpointsByAPIGatewayID[id].Value()
		apiGateways[0].Lambdas = append(apiGateways[0].Lambdas, apiGatewayLambdas[id]...)
	}
	return apiGateways
}

func buildSQSs(resourcesByTypeMap map[drawio.ResourceType][]drawio.Resource) []config.SQS {
	sqss := []config.SQS{}

	for _, sqs := range resourcesByTypeMap[drawio.SQSType] {
		sqss = append(sqss, config.SQS{Name: sqs.Value(), MaxReceiveCount: 10})
	}

	return sqss
}

func buildBuckets(resourcesByTypeMap map[drawio.ResourceType][]drawio.Resource) []config.S3 {
	buckets := []config.S3{}

	for _, bucket := range resourcesByTypeMap[drawio.S3Type] {
		buckets = append(buckets, config.S3{Name: bucket.Value(), ExpirationDays: 90})
	}

	return buckets
}

func buildRestfulAPIs(resourcesByTypeMap map[drawio.ResourceType][]drawio.Resource) []config.RestfulAPI {
	restfulAPIs := []config.RestfulAPI{}
	restfulAPINames := map[string]struct{}{}

	for _, restfulAPI := range resourcesByTypeMap[drawio.RestfulAPIType] {
		name := restfulAPI.Value()
		if _, ok := restfulAPINames[name]; !ok {
			restfulAPIs = append(restfulAPIs, config.RestfulAPI{Name: name})
			restfulAPINames[name] = struct{}{}
		}
	}

	return restfulAPIs
}

func buildSNSs(snsMap map[string]config.SNS) []config.SNS {
	snss := make([]config.SNS, 0, len(snsMap))

	for _, sns := range snsMap {
		snss = append(snss, sns)
	}

	return snss
}
