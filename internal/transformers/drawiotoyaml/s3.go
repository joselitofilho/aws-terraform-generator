package drawiotoyaml

import (
	"github.com/joselitofilho/aws-terraform-generator/internal/drawio"
	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
)

func buildS3Relationship(source, target drawio.Resource, envars map[string]map[string]string) {
	if source.ResourceType() == drawio.LambdaType {
		buildLambdaToS3(source, target, envars)
	}
}

func buildS3Buckets(resourcesByTypeMap map[drawio.ResourceType][]drawio.Resource) []config.S3 {
	var buckets []config.S3

	for _, bucket := range resourcesByTypeMap[drawio.S3Type] {
		buckets = append(buckets, config.S3{Name: bucket.Value(), ExpirationDays: 90})
	}

	return buckets
}
