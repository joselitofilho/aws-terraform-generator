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

	snsResourceEvents := []string{"s3:ObjectCreated:*"}

	for _, s := range t.resourcesByTypeMap[awsresources.SNSType] {
		var bucketName string
		if s3Bucket, ok := t.s3BucketsBySNSID[s.ID()]; ok {
			bucketName = s3Bucket.Value()
		}

		var lambdas []config.SNSResource
		for _, l := range t.lambdasBySNSID[s.ID()] {
			lambdas = append(lambdas, config.SNSResource{
				Name:   l.Value(),
				Events: snsResourceEvents,
			})
		}

		var sqss []config.SNSResource
		for _, sqs := range t.sqssBySNSID[s.ID()] {
			sqss = append(sqss, config.SNSResource{
				Name:   sqs.Value(),
				Events: snsResourceEvents,
			})
		}

		snss = append(snss, config.SNS{
			Name:       s.Value(),
			BucketName: bucketName,
			Lambdas:    lambdas,
			SQSs:       sqss,
		})
	}

	return snss
}
