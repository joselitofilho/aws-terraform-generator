package resources_to_yaml

import (
	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
	"github.com/joselitofilho/aws-terraform-generator/internal/resources"
)

func (t *Transformer) buildS3Relationship(source, target resources.Resource) {
	if source.ResourceType() == resources.LambdaType {
		t.buildLambdaToS3(source, target)
	}
}

func (t *Transformer) buildS3Buckets() []config.S3 {
	var buckets []config.S3

	for _, bucket := range t.resourcesByTypeMap[resources.S3Type] {
		buckets = append(buckets, config.S3{Name: bucket.Value(), ExpirationDays: 90})
	}

	return buckets
}
