package resourcestoyaml

import (
	"github.com/diagram-code-generator/resources/pkg/resources"

	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
	awsresources "github.com/joselitofilho/aws-terraform-generator/internal/resources"
)

func (t *Transformer) buildSNSRelationship(source, target resources.Resource) {
	if awsresources.ParseResourceType(source.ResourceType()) == awsresources.S3Type {
		t.buildS3ToSNS(source, target)
	}
}

func (t *Transformer) buildSNSs() []config.SNS {
	var snss []config.SNS

	for _, sns := range t.snsMap {
		snss = append(snss, sns)
	}

	return snss
}
