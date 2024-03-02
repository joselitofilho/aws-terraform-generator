package transformers

import (
	"github.com/joselitofilho/aws-terraform-generator/internal/drawio"
	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
)

func buildAPIGatewayRelationship(
	source, target drawio.Resource, apiGatewaysByID, endpointsByAPIGatewayID map[string]drawio.Resource,
) {
	if source.ResourceType() == drawio.EndpointType {
		buildEndpointToAPIGateway(apiGatewaysByID, endpointsByAPIGatewayID, source, target)
	}
}

func buildAPIGateways(
	yamlConfig *config.Config,
	apiGatewaysByID map[string]drawio.Resource,
	endpointsByAPIGatewayID map[string]drawio.Resource,
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
