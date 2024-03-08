package resources_to_yaml

import (
	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
	"github.com/joselitofilho/aws-terraform-generator/internal/resources"
)

type Transformer struct {
	yamlConfig *config.Config
	resc       *resources.ResourceCollection

	apiGatewaysByID           map[string]resources.Resource
	cronsByLambdaID           map[string]resources.Resource
	endpointsByAPIGatewayID   map[string]resources.Resource
	kinesisTriggersByLambdaID map[string][]resources.Resource
	sqsTriggersByLambdaID     map[string][]resources.Resource

	snsMap map[string]config.SNS

	envars map[string]map[string]string

	resourcesByTypeMap map[resources.ResourceType][]resources.Resource
}

func NewTransformer(yamlConfig *config.Config, resc *resources.ResourceCollection) *Transformer {
	return &Transformer{
		yamlConfig: yamlConfig,
		resc:       resc,

		apiGatewaysByID:           map[string]resources.Resource{},
		cronsByLambdaID:           map[string]resources.Resource{},
		endpointsByAPIGatewayID:   map[string]resources.Resource{},
		kinesisTriggersByLambdaID: map[string][]resources.Resource{},
		sqsTriggersByLambdaID:     map[string][]resources.Resource{},

		snsMap: map[string]config.SNS{},

		envars: map[string]map[string]string{},

		resourcesByTypeMap: map[resources.ResourceType][]resources.Resource{},
	}
}

func (t *Transformer) Transform() (*config.Config, error) {
	t.buildResourcesByTypeMap()

	for _, sns := range t.resourcesByTypeMap[resources.SNSType] {
		t.snsMap[sns.ID()] = config.SNS{Name: sns.Value()}
	}

	for i := range t.resourcesByTypeMap[resources.APIGatewayType] {
		apiGateway := t.resourcesByTypeMap[resources.APIGatewayType][i]
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
		t.resourcesByTypeMap[resource.ResourceType()] = append(t.resourcesByTypeMap[resource.ResourceType()], resource)
	}
}

func (t *Transformer) buildResourceRelationships() {
	for _, rel := range t.resc.Relationships {
		target := rel.Target
		source := rel.Source

		switch target.ResourceType() {
		case resources.APIGatewayType:
			t.buildAPIGatewayRelationship(source, target)
		case resources.GoogleBQType:
			t.buildGoogleBQRelationship(source, target)
		case resources.DatabaseType:
			t.buildDatabaseRelationship(source, target)
		case resources.KinesisType:
			t.buildKinesisRelationship(source, target)
		case resources.LambdaType:
			t.buildLambdaRelationships(source, target)
		case resources.RestfulAPIType:
			t.buildRestfulAPIRelationship(source, target)
		case resources.S3Type:
			t.buildS3Relationship(source, target)
		case resources.SNSType:
			t.buildSNSRelationship(source, target)
		case resources.SQSType:
			t.buildSQSRelationships(source, target)
		}
	}
}
