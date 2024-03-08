package drawiotoyaml

import (
	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
	"github.com/joselitofilho/aws-terraform-generator/internal/resources"
)

func buildRestfulAPIRelationship(source, target resources.Resource, envars map[string]map[string]string) {
	if source.ResourceType() == resources.LambdaType {
		buildLambdaToRestfulAPI(source, target, envars)
	}
}

func buildRestfulAPIs(resourcesByTypeMap map[resources.ResourceType][]resources.Resource) []config.RestfulAPI {
	var restfulAPIs []config.RestfulAPI

	restfulAPINames := map[string]struct{}{}

	for _, restfulAPI := range resourcesByTypeMap[resources.RestfulAPIType] {
		name := restfulAPI.Value()
		if _, ok := restfulAPINames[name]; !ok {
			restfulAPIs = append(restfulAPIs, config.RestfulAPI{Name: name})
			restfulAPINames[name] = struct{}{}
		}
	}

	return restfulAPIs
}
