package resourcestoyaml

import (
	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
	"github.com/joselitofilho/aws-terraform-generator/internal/resources"
)

func (t *Transformer) buildRestfulAPIRelationship(source, target resources.Resource) {
	if source.ResourceType() == resources.LambdaType {
		t.buildLambdaToRestfulAPI(source, target)
	}
}

func (t *Transformer) buildRestfulAPIs() []config.RestfulAPI {
	var restfulAPIs []config.RestfulAPI

	restfulAPINames := map[string]struct{}{}

	for _, restfulAPI := range t.resourcesByTypeMap[resources.RestfulAPIType] {
		name := restfulAPI.Value()
		if _, ok := restfulAPINames[name]; !ok {
			restfulAPIs = append(restfulAPIs, config.RestfulAPI{Name: name})
			restfulAPINames[name] = struct{}{}
		}
	}

	return restfulAPIs
}
