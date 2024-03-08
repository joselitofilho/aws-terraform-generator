package drawiotoyaml

import (
	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
	"github.com/joselitofilho/aws-terraform-generator/internal/resources"
)

func buildSNSRelationship(source, target resources.Resource, snsMap map[string]config.SNS) {
	if source.ResourceType() == resources.S3Type {
		buildS3ToSNS(snsMap, source, target)
	}
}

func buildSNSs(snsMap map[string]config.SNS) []config.SNS {
	var snss []config.SNS

	for _, sns := range snsMap {
		snss = append(snss, sns)
	}

	return snss
}
