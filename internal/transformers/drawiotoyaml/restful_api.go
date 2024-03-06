package drawiotoyaml

import (
	"github.com/joselitofilho/aws-terraform-generator/internal/drawio"
	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
)

func buildRestfulAPIRelationship(source, target drawio.Resource, envars map[string]map[string]string) {
	if source.ResourceType() == drawio.LambdaType {
		buildLambdaToRestfulAPI(source, target, envars)
	}
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
