package yamltoresources

import (
	_ "embed"
	"testing"

	"github.com/diagram-code-generator/resources/pkg/resources"

	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
	awsresources "github.com/joselitofilho/aws-terraform-generator/internal/resources"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

//go:embed testdata/diagram.yaml
var diagramData []byte

var (
	endpointResource = resources.NewGenericResource(
		"1", "mystack-api.domain-${var.environment}.com", awsresources.EndpointType.String())
	apiGateway       = resources.NewGenericResource("2", "POST /v1/examples", awsresources.APIGatewayType.String())
	lambdaAPIGateway = resources.NewGenericResource("3", "exampleApiReceiver", awsresources.LambdaType.String())
	kinesis          = resources.NewGenericResource("4", "MyKinesis", awsresources.KinesisType.String())
	lambda           = resources.NewGenericResource("5", "exampleReceiver", awsresources.LambdaType.String())
	cron             = resources.NewGenericResource("6", "cron(0 1 * * ? *)", awsresources.CronType.String())
	restAPI          = resources.NewGenericResource("7", "MyApi", awsresources.RestfulAPIType.String())
	s3Bucket         = resources.NewGenericResource("8", "my-bucket", awsresources.S3Type.String())
	sns              = resources.NewGenericResource("9", "example", awsresources.SNSType.String())
	targetSQS        = resources.NewGenericResource("10", "target", awsresources.SQSType.String())
	sourceSQS        = resources.NewGenericResource("11", "source", awsresources.SQSType.String())

	wantResourceCollection = &resources.ResourceCollection{
		Resources: []resources.Resource{
			endpointResource,
			apiGateway,
			lambdaAPIGateway,
			kinesis,
			lambda,
			cron,
			restAPI,
			s3Bucket,
			sns,
			targetSQS,
			sourceSQS,
		},
		Relationships: []resources.Relationship{
			{
				Source: apiGateway,
				Target: lambdaAPIGateway,
			},
			{
				Source: endpointResource,
				Target: apiGateway,
			},
			{
				Source: cron,
				Target: lambda,
			},
			{
				Source: lambda,
				Target: targetSQS,
			},
			{
				Source: kinesis,
				Target: lambda,
			},
			{
				Source: sourceSQS,
				Target: lambda,
			},
		},
	}
)

func setupEmpty(
	id *int, wantResources *[]resources.Resource, wantRelationships *[]resources.Relationship,
) func() {
	return func() {
		*id = 2
		*wantResources = []resources.Resource{}
		*wantRelationships = []resources.Relationship{}
	}
}

func setupWithLambda(
	lambdaResource, targetResource resources.Resource,
	id *int, wantResources *[]resources.Resource, wantRelationships *[]resources.Relationship,
) func() {
	return func() {
		*id = 2
		*wantResources = []resources.Resource{targetResource}
		*wantRelationships = []resources.Relationship{{
			Source: lambdaResource,
			Target: targetResource,
		}}
	}
}

func TestTransformer_Transform(t *testing.T) {
	type fields struct {
		yamlConfig *config.Config
	}

	var diagramYAML *config.Config

	err := yaml.Unmarshal(diagramData, &diagramYAML)
	require.NoError(t, err)

	tests := []struct {
		name      string
		fields    fields
		want      *resources.ResourceCollection
		targetErr error
	}{
		{
			name:   "happy path",
			fields: fields{yamlConfig: diagramYAML},
			want:   wantResourceCollection,
		},
		{
			name:      "when YAML is invalid or empty should return an error",
			fields:    fields{yamlConfig: nil},
			want:      nil,
			targetErr: ErrEmptyConfig,
		},
	}

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			tr := NewTransformer(tc.fields.yamlConfig)

			got, err := tr.Transform()

			require.ErrorIs(t, err, tc.targetErr)
			if tc.want == nil {
				require.Nil(t, got)
			} else {
				require.True(t, tc.want.Equal(got))
			}
		})
	}
}

func TestTransformer_transformLambdaEnvars(t *testing.T) {
	type fields struct {
		yamlConfig *config.Config
	}

	type args struct {
		res           *config.Lambda
		lambda        *resources.GenericResource
		lambdaARN     awsresources.ResourceARN
		resources     *[]resources.Resource
		relationships *[]resources.Relationship
		id            *int
	}

	lambdaResource := resources.NewGenericResource("1", "myReceiver", awsresources.LambdaType.String())

	var (
		id                int
		wantResources     []resources.Resource
		wantRelationships []resources.Relationship
	)

	tests := []struct {
		name   string
		fields fields
		args   args
		setup  func()
		wantID int
	}{
		{
			name: "lambda and database",
			fields: fields{
				yamlConfig: &config.Config{},
			},
			args: args{
				res:           &config.Lambda{Envars: map[string]string{"DOCDB_HOST": "var.doc_db_host"}},
				lambda:        lambdaResource,
				resources:     &[]resources.Resource{},
				relationships: &[]resources.Relationship{},
				id:            &id,
			},
			setup: setupWithLambda(lambdaResource, resources.NewGenericResource("2", "doc", awsresources.DatabaseType.String()),
				&id, &wantResources, &wantRelationships),
			wantID: 3,
		},
		{
			name: "lambda and google BQ",
			fields: fields{
				yamlConfig: &config.Config{},
			},
			args: args{
				res:           &config.Lambda{Envars: map[string]string{"GOOGLE_BQ_PROJECT_ID": "google"}},
				lambda:        lambdaResource,
				resources:     &[]resources.Resource{},
				relationships: &[]resources.Relationship{},
				id:            &id,
			},
			setup: setupWithLambda(lambdaResource, resources.NewGenericResource("2", "google", awsresources.GoogleBQType.String()),
				&id, &wantResources, &wantRelationships),
			wantID: 3,
		},
		{
			name: "lambda and kinesis",
			fields: fields{
				yamlConfig: &config.Config{},
			},
			args: args{
				res:           &config.Lambda{Envars: map[string]string{"MY_KINESIS_KINESIS_STREAM_URL": "MyKinesis"}},
				lambda:        lambdaResource,
				resources:     &[]resources.Resource{},
				relationships: &[]resources.Relationship{},
				id:            &id,
			},
			setup:  setupEmpty(&id, &wantResources, &wantRelationships),
			wantID: 2,
		},
		{
			name: "lambda and s3 bucket S3_BUCKET",
			fields: fields{
				yamlConfig: &config.Config{},
			},
			args: args{
				res:           &config.Lambda{Envars: map[string]string{"PAYLOADS_S3_BUCKET": "payloads"}},
				lambda:        lambdaResource,
				resources:     &[]resources.Resource{},
				relationships: &[]resources.Relationship{},
				id:            &id,
			},
			setup:  setupEmpty(&id, &wantResources, &wantRelationships),
			wantID: 2,
		},
		{
			name: "lambda and s3 bucket BUCKET_NAME",
			fields: fields{
				yamlConfig: &config.Config{},
			},
			args: args{
				res:           &config.Lambda{Envars: map[string]string{"PAYLOADS_BUCKET_NAME": "payloads"}},
				lambda:        lambdaResource,
				resources:     &[]resources.Resource{},
				relationships: &[]resources.Relationship{},
				id:            &id,
			},
			setup:  setupEmpty(&id, &wantResources, &wantRelationships),
			wantID: 2,
		},
		{
			name: "lambda and restful API",
			fields: fields{
				yamlConfig: &config.Config{},
			},
			args: args{
				res:           &config.Lambda{Envars: map[string]string{"MY_REST_API_BASE_URL": "MyRest"}},
				lambda:        lambdaResource,
				resources:     &[]resources.Resource{},
				relationships: &[]resources.Relationship{},
				id:            &id,
			},
			setup: setupWithLambda(lambdaResource, resources.NewGenericResource("2", "MyRest", awsresources.RestfulAPIType.String()),
				&id, &wantResources, &wantRelationships),
			wantID: 3,
		},
	}

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			tr := NewTransformer(tc.fields.yamlConfig)

			tr.transformLambdaEnvars(tc.args.res, tc.args.lambda, tc.args.lambdaARN,
				tc.args.resources, tc.args.relationships, tc.args.id)

			require.Equal(t, tc.wantID, *tc.args.id)
			require.Equal(t, &wantResources, tc.args.resources)
			require.Equal(t, &wantRelationships, tc.args.relationships)
		})
	}
}
