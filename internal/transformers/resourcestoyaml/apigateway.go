package resourcestoyaml

import (
	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
	"github.com/joselitofilho/aws-terraform-generator/internal/resources"
)

func (t *Transformer) buildAPIGatewayRelationship(source, target resources.Resource) {
	if source.ResourceType() == resources.EndpointType {
		t.buildEndpointToAPIGateway(source, target)
	}
}

func (t *Transformer) buildAPIGateways(
	apiGatewayLambdasByAPIGatewayID map[string][]config.APIGatewayLambda,
) (apiGateways []config.APIGateway) {
	for id := range t.apiGatewaysByID {
		var apiDomainValue string
		if rsc, ok := t.endpointsByAPIGatewayID[id]; ok {
			apiDomainValue = rsc.Value()
		}

		apiGateways = append(apiGateways, config.APIGateway{
			StackName: t.yamlConfig.Diagram.StackName,
			APIG:      true,
			APIDomain: apiDomainValue,
			Lambdas:   apiGatewayLambdasByAPIGatewayID[id],
		})
	}

	return apiGateways
}
