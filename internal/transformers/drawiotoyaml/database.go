package drawiotoyaml

import "github.com/joselitofilho/aws-terraform-generator/internal/drawio"

func buildDatabaseRelationship(source, target drawio.Resource, envars map[string]map[string]string) {
	if source.ResourceType() == drawio.LambdaType {
		buildLambdaToDatabase(source, target, envars)
	}
}
