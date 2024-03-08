package resources_to_yaml

import (
	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
	"github.com/joselitofilho/aws-terraform-generator/internal/resources"
)

func buildS3Relationship(source, target resources.Resource, envars map[string]map[string]string) {
	if source.ResourceType() == resources.LambdaType {
		buildLambdaToS3(source, target, envars)
	}
}

func buildS3Buckets(resourcesByTypeMap map[resources.ResourceType][]resources.Resource) []config.S3 {
	var buckets []config.S3

	for _, bucket := range resourcesByTypeMap[resources.S3Type] {
		buckets = append(buckets, config.S3{Name: bucket.Value(), ExpirationDays: 90})
	}

	return buckets
}
