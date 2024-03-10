package resourcestoyaml

import (
	"github.com/joselitofilho/aws-terraform-generator/internal/resources"
)

func (t *Transformer) buildDatabaseRelationship(source, target resources.Resource) {
	if source.ResourceType() == resources.LambdaType {
		t.buildLambdaToDatabase(source, target)
	}
}
