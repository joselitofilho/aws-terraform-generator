package resourcestoyaml

import (
	"github.com/diagram-code-generator/resources/pkg/resources"

	awsresources "github.com/joselitofilho/aws-terraform-generator/internal/resources"
)

func (t *Transformer) buildGoogleBQRelationship(source, target resources.Resource) {
	if awsresources.ParseResourceType(source.ResourceType()) == awsresources.LambdaType {
		t.buildLambdaToGoogleBQ(source, target)
	}
}
