package drawiotoyaml

import "github.com/joselitofilho/aws-terraform-generator/internal/drawio"

func buildGoogleBQRelationship(source, target drawio.Resource, envars map[string]map[string]string) {
	if source.ResourceType() == drawio.LambdaType {
		buildLambdaToGoogleBQ(source, target, envars)
	}
}
