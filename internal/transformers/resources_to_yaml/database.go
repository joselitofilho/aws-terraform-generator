package resources_to_yaml

import (
	"github.com/joselitofilho/aws-terraform-generator/internal/resources"
)

func buildDatabaseRelationship(source, target resources.Resource, envars map[string]map[string]string) {
	if source.ResourceType() == resources.LambdaType {
		buildLambdaToDatabase(source, target, envars)
	}
}
