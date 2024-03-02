package transformers

import (
	"github.com/joselitofilho/aws-terraform-generator/internal/drawio"
	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
)

func buildSNSRelationship(source, target drawio.Resource, snsMap map[string]config.SNS) {
	if source.ResourceType() == drawio.S3Type {
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
