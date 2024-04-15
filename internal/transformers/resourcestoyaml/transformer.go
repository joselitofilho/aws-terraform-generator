package resourcestoyaml

import (
	"github.com/diagram-code-generator/resources/pkg/resources"
	awsresources "github.com/joselitofilho/aws-terraform-generator/internal/resources"

	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
)

type Transformer struct {
	yamlConfig *config.Config
	resc       *resources.ResourceCollection

	apiGatewaysByID           map[string]resources.Resource
	cronsByLambdaID           map[string]resources.Resource
	endpointsByAPIGatewayID   map[string]resources.Resource
	kinesisTriggersByLambdaID map[string][]resources.Resource
	lambdasBySNSID            map[string][]resources.Resource
	s3BucketsBySNSID          map[string]resources.Resource
	sqssBySNSID               map[string][]resources.Resource
	sqsTriggersByLambdaID     map[string][]resources.Resource

	envars map[string]map[string]string

	resourcesByTypeMap map[awsresources.ResourceType][]resources.Resource
}

func NewTransformer(yamlConfig *config.Config, resc *resources.ResourceCollection) *Transformer {
	return &Transformer{
		yamlConfig: yamlConfig,
		resc:       resc,

		apiGatewaysByID:           map[string]resources.Resource{},
		cronsByLambdaID:           map[string]resources.Resource{},
		endpointsByAPIGatewayID:   map[string]resources.Resource{},
		kinesisTriggersByLambdaID: map[string][]resources.Resource{},
		lambdasBySNSID:            map[string][]resources.Resource{},
		s3BucketsBySNSID:          map[string]resources.Resource{},
		sqsTriggersByLambdaID:     map[string][]resources.Resource{},
		sqssBySNSID:               map[string][]resources.Resource{},

		envars: map[string]map[string]string{},

		resourcesByTypeMap: map[awsresources.ResourceType][]resources.Resource{},
	}
}

func (t *Transformer) Transform() (*config.Config, error) {
	t.buildResourcesByTypeMap()

	for i := range t.resourcesByTypeMap[awsresources.APIGatewayType] {
		apiGateway := t.resourcesByTypeMap[awsresources.APIGatewayType][i]
		t.apiGatewaysByID[apiGateway.ID()] = apiGateway
	}

	t.buildResourceRelationships()

	lambdas, apiGatewayLambdasByAPIGatewayID := t.buildLambdas()
	apiGateways := t.buildAPIGateways(apiGatewayLambdasByAPIGatewayID)
	kinesis := t.buildKinesis()
	snss := t.buildSNSs()
	sqss := t.buildSQSs()
	buckets := t.buildS3Buckets()
	restfulAPIs := t.buildRestfulAPIs()

	return &config.Config{
		Lambdas:     lambdas,
		APIGateways: apiGateways,
		Kinesis:     kinesis,
		SNSs:        snss,
		SQSs:        sqss,
		Buckets:     buckets,
		RestfulAPIs: restfulAPIs,
	}, nil
}

func (t *Transformer) buildResourcesByTypeMap() {
	for _, resource := range t.resc.Resources {
		resType := awsresources.ParseResourceType(resource.ResourceType())
		t.resourcesByTypeMap[resType] = append(t.resourcesByTypeMap[resType], resource)
	}
}

func (t *Transformer) buildResourceRelationships() {
	for _, rel := range t.resc.Relationships {
		target := rel.Target
		source := rel.Source

		switch awsresources.ParseResourceType(target.ResourceType()) {
		case awsresources.APIGatewayType:
			t.buildAPIGatewayRelationship(source, target)
		case awsresources.GoogleBQType:
			t.buildGoogleBQRelationship(source, target)
		case awsresources.DatabaseType:
			t.buildDatabaseRelationship(source, target)
		case awsresources.KinesisType:
			t.buildKinesisRelationship(source, target)
		case awsresources.LambdaType:
			t.buildLambdaRelationships(source, target)
		case awsresources.RestfulAPIType:
			t.buildRestfulAPIRelationship(source, target)
		case awsresources.S3Type:
			t.buildS3Relationship(source, target)
		case awsresources.SNSType:
			t.buildSNSRelationship(source, target)
		case awsresources.SQSType:
			t.buildSQSRelationships(source, target)
		}
	}
}
