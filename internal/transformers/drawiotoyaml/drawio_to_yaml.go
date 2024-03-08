package drawiotoyaml

import (
	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
	"github.com/joselitofilho/aws-terraform-generator/internal/resources"
)

func TransformDrawIOToYAML(yamlConfig *config.Config, resc *resources.ResourceCollection) (*config.Config, error) {
	apiGatewaysByID := map[string]resources.Resource{}
	cronsByLambdaID := map[string]resources.Resource{}
	endpointsByAPIGatewayID := map[string]resources.Resource{}
	kinesisTriggersByLambdaID := map[string][]resources.Resource{}
	sqsTriggersByLambdaID := map[string][]resources.Resource{}

	snsMap := map[string]config.SNS{}

	envars := map[string]map[string]string{}

	resourcesByTypeMap := buildResourcesByTypeMap(resc)

	for _, sns := range resourcesByTypeMap[resources.SNSType] {
		snsMap[sns.ID()] = config.SNS{Name: sns.Value()}
	}

	for i := range resourcesByTypeMap[resources.APIGatewayType] {
		apiGateway := resourcesByTypeMap[resources.APIGatewayType][i]
		apiGatewaysByID[apiGateway.ID()] = apiGateway
	}

	buildResourceRelationships(resc, envars,
		apiGatewaysByID, cronsByLambdaID, endpointsByAPIGatewayID,
		kinesisTriggersByLambdaID, sqsTriggersByLambdaID,
		snsMap)

	lambdas, apiGatewayLambdasByAPIGatewayID := buildLambdas(
		yamlConfig, resourcesByTypeMap, resc, envars, cronsByLambdaID,
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

func buildResourcesByTypeMap(resc *resources.ResourceCollection) map[resources.ResourceType][]resources.Resource {
	resourcesByTypeMap := map[resources.ResourceType][]resources.Resource{}

	for _, resource := range resc.Resources {
		resourcesByTypeMap[resource.ResourceType()] = append(resourcesByTypeMap[resource.ResourceType()], resource)
	}

	return resourcesByTypeMap
}

func buildResourceRelationships(
	resc *resources.ResourceCollection,
	envars map[string]map[string]string,
	apiGatewaysByID, cronsByLambdaID, endpointsByAPIGatewayID map[string]resources.Resource,
	kinesisTriggersByLambdaID, sqsTriggersByLambdaID map[string][]resources.Resource,
	snsMap map[string]config.SNS,
) {
	for _, rel := range resc.Relationships {
		target := rel.Target
		source := rel.Source

		switch target.ResourceType() {
		case resources.APIGatewayType:
			buildAPIGatewayRelationship(source, target, apiGatewaysByID, endpointsByAPIGatewayID)
		case resources.GoogleBQType:
			buildGoogleBQRelationship(source, target, envars)
		case resources.DatabaseType:
			buildDatabaseRelationship(source, target, envars)
		case resources.KinesisType:
			buildKinesisRelationship(source, target, envars)
		case resources.LambdaType:
			buildLambdaRelationships(
				source, target, cronsByLambdaID, kinesisTriggersByLambdaID, sqsTriggersByLambdaID, snsMap)
		case resources.RestfulAPIType:
			buildRestfulAPIRelationship(source, target, envars)
		case resources.S3Type:
			buildS3Relationship(source, target, envars)
		case resources.SNSType:
			buildSNSRelationship(source, target, snsMap)
		case resources.SQSType:
			buildSQSRelationships(source, target, envars, snsMap)
		}
	}
}
