package transformers

import (
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

	buildResourceRelationships(resources, envars,
		apiGatewaysByID, cronsByLambdaID, endpointsByAPIGatewayID,
		kinesisTriggersByLambdaID, sqsTriggersByLambdaID,
		snsMap)

	lambdas, apiGatewayLambdasByAPIGatewayID := buildLambdas(
		yamlConfig, resourcesByTypeMap, resources, envars, cronsByLambdaID,
		kinesisTriggersByLambdaID, sqsTriggersByLambdaID)
	apiGateways := buildAPIGateways(
		yamlConfig, apiGatewaysByID, endpointsByAPIGatewayID, apiGatewayLambdasByAPIGatewayID)
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
		resourcesByTypeMap[resource.ResourceType()] = append(resourcesByTypeMap[resource.ResourceType()], resource)
	}

	return resourcesByTypeMap
}

func buildResourceRelationships(
	resources *drawio.ResourceCollection,
	envars map[string]map[string]string,
	apiGatewaysByID, cronsByLambdaID, endpointsByAPIGatewayID map[string]drawio.Resource,
	kinesisTriggersByLambdaID, sqsTriggersByLambdaID map[string][]drawio.Resource,
	snsMap map[string]config.SNS,
) {
	for _, rel := range resources.Relationships {
		target := rel.Target
		source := rel.Source

		switch target.ResourceType() {
		case drawio.APIGatewayType:
			buildAPIGatewayRelationship(source, target, apiGatewaysByID, endpointsByAPIGatewayID)
		case drawio.DatabaseType:
			buildDatabaseRelationship(source, target, envars)
		case drawio.KinesisType:
			buildKinesisRelationship(source, target, envars)
		case drawio.LambdaType:
			buildLambdaRelationships(
				source, target, cronsByLambdaID, kinesisTriggersByLambdaID, sqsTriggersByLambdaID, snsMap)
		case drawio.RestfulAPIType:
			buildRestfulAPIRelationship(source, target, envars)
		case drawio.S3Type:
			buildS3Relationship(source, target, envars)
		case drawio.SNSType:
			buildSNSRelationship(source, target, snsMap)
		case drawio.SQSType:
			buildSQSRelationships(source, target, envars, snsMap)
		}
	}
}
