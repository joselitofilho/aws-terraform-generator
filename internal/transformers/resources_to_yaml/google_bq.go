package resources_to_yaml

import (
	"github.com/joselitofilho/aws-terraform-generator/internal/resources"
)

func buildGoogleBQRelationship(source, target resources.Resource, envars map[string]map[string]string) {
	if source.ResourceType() == resources.LambdaType {
		buildLambdaToGoogleBQ(source, target, envars)
	}
}
