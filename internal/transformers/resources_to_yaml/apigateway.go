package resources_to_yaml

import (
	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
	"github.com/joselitofilho/aws-terraform-generator/internal/resources"
)

func buildAPIGatewayRelationship(
	source, target resources.Resource, apiGatewaysByID, endpointsByAPIGatewayID map[string]resources.Resource,
) {
	if source.ResourceType() == resources.EndpointType {
		buildEndpointToAPIGateway(apiGatewaysByID, endpointsByAPIGatewayID, source, target)
	}
}

func buildAPIGateways(
	yamlConfig *config.Config,
	apiGatewaysByID map[string]resources.Resource,
	endpointsByAPIGatewayID map[string]resources.Resource,
	apiGatewayLambdasByAPIGatewayID map[string][]config.APIGatewayLambda,
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
			Lambdas:   apiGatewayLambdasByAPIGatewayID[id],
		})
	}

	return apiGateways
}
