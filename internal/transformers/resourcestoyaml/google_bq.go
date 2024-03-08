package resourcestoyaml

import (
	"github.com/joselitofilho/aws-terraform-generator/internal/resources"
)

func (t *Transformer) buildGoogleBQRelationship(source, target resources.Resource) {
	if source.ResourceType() == resources.LambdaType {
		t.buildLambdaToGoogleBQ(source, target)
	}
}
