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
	kinesisTriggersByLambdaID := map[string][]drawio.Resource{}
	sqsTriggersByLambdaID := map[string][]drawio.Resource{}
	envars := map[string]map[string]string{}
	snsMap := map[string]config.SNS{}

	resourcesByTypeMap := buildResourcesByTypeMap(resources)

	for _, sns := range resourcesByTypeMap[drawio.SNSType] {
		snsMap[sns.ID()] = config.SNS{Name: sns.Value()}
	}

	for i := range resourcesByTypeMap[drawio.APIGatewayType] {
		apiGateway := resourcesByTypeMap[drawio.APIGatewayType][i]
		apiGatewaysByID[apiGateway.ID()] = apiGateway
	}

	buildResourceRelationships(
		resources, cronsByLambdaID, kinesisTriggersByLambdaID, sqsTriggersByLambdaID, snsMap, apiGatewaysByID,
		endpointsByAPIGatewayID, envars)

	lambdas, apiGatewayLambdas := buildLambdas(
		yamlConfig, resourcesByTypeMap, resources, envars, cronsByLambdaID,
		kinesisTriggersByLambdaID, sqsTriggersByLambdaID)
	apiGateways := buildAPIGateways(yamlConfig, apiGatewaysByID, endpointsByAPIGatewayID, apiGatewayLambdas)
	kinesis := buildKinesis(resourcesByTypeMap)
	snss := buildSNSs(snsMap)
	sqss := buildSQSs(resourcesByTypeMap)
	buckets := buildS3Buckets(resourcesByTypeMap)
	restfulAPIs := buildRestfulAPIs(resourcesByTypeMap)

	return &config.Config{
		Lambdas:     lambdas,
		APIGateways: apiGateways,
		Kinesis:     kinesis,
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
	kinesisTriggersByLambdaID map[string][]drawio.Resource,
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
			case drawio.KinesisType:
				buildKinesisToLambda(kinesisTriggersByLambdaID, source, target)
			case drawio.SQSType:
				buildSQSToLambda(sqsTriggersByLambdaID, source, target)
			case drawio.SNSType:
				buildSNSToLambda(snsMap, source, target)
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
				buildSNSToSQS(snsMap, source, target)
			}
		case drawio.DatabaseType:
			if source.ReseourceType() == drawio.LambdaType {
				buildLambdaToDatabase(envars, source, target)
			}
		case drawio.KinesisType:
			if source.ReseourceType() == drawio.LambdaType {
				buildLambdaToKinesis(envars, source, target)
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
				rel.Source.ReseourceType() == drawio.APIGatewayType {
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

func buildAPIGateways(
	yamlConfig *config.Config,
	apiGatewaysByID map[string]drawio.Resource,
	endpointsByAPIGatewayID map[string]drawio.Resource,
	apiGatewayLambdas map[string][]config.APIGatewayLambda,
) (apiGateways []config.APIGateway) {
	for id := range apiGatewaysByID {
		var apiDomainValue string
		if rsc, ok := endpointsByAPIGatewayID[id]; ok {
			apiDomainValue = rsc.Value()
		}

		apiGateways = append(apiGateways, config.APIGateway{
			StackName: yamlConfig.Diagram.StackName,
			APIG:      true,
			APIDomain: apiDomainValue,
			Lambdas:   apiGatewayLambdas[id],
		})
	}

	return apiGateways
}

func buildKinesis(resourcesByTypeMap map[drawio.ResourceType][]drawio.Resource) []config.Kinesis {
	var kinesis []config.Kinesis

	for _, k := range resourcesByTypeMap[drawio.KinesisType] {
		kinesis = append(kinesis, config.Kinesis{Name: k.Value(), RetentionPeriod: "24"})
	}

	return kinesis
}

func buildS3Buckets(resourcesByTypeMap map[drawio.ResourceType][]drawio.Resource) []config.S3 {
	var buckets []config.S3

	for _, bucket := range resourcesByTypeMap[drawio.S3Type] {
		buckets = append(buckets, config.S3{Name: bucket.Value(), ExpirationDays: 90})
	}

	return buckets
}

func buildSQSs(resourcesByTypeMap map[drawio.ResourceType][]drawio.Resource) []config.SQS {
	var sqss []config.SQS

	for _, sqs := range resourcesByTypeMap[drawio.SQSType] {
		sqss = append(sqss, config.SQS{Name: sqs.Value(), MaxReceiveCount: 10})
	}

	return sqss
}

func buildSNSs(snsMap map[string]config.SNS) []config.SNS {
	var snss []config.SNS

	for _, sns := range snsMap {
		snss = append(snss, sns)
	}

	return snss
}

func buildRestfulAPIs(resourcesByTypeMap map[drawio.ResourceType][]drawio.Resource) []config.RestfulAPI {
	var restfulAPIs []config.RestfulAPI

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
