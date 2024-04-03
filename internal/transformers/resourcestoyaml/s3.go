package resourcestoyaml

import (
	"github.com/diagram-code-generator/resources/pkg/resources"

	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
	awsresources "github.com/joselitofilho/aws-terraform-generator/internal/resources"
)

func (t *Transformer) buildS3Relationship(source, target resources.Resource) {
	if awsresources.ParseResourceType(source.ResourceType()) == awsresources.LambdaType {
		t.buildLambdaToS3(source, target)
	}
}

func (t *Transformer) buildS3Buckets() []config.S3 {
	var buckets []config.S3

	for _, bucket := range t.resourcesByTypeMap[awsresources.S3Type] {
		buckets = append(buckets, config.S3{Name: bucket.Value(), ExpirationDays: 90})
	}

	return buckets
}
