package yamltoresources

import (
	_ "embed"
	"testing"

	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
	"github.com/joselitofilho/aws-terraform-generator/internal/resources"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

//go:embed testdata/diagram.yaml
var diagramData []byte

var (
	endpointResource = resources.NewGenericResource("1", "mystack-api.domain-${var.environment}.com",
		resources.EndpointType)
	apiGateway        = resources.NewGenericResource("2", "POST /v1/examples", resources.APIGatewayType)
	lambdaAPIGateway  = resources.NewGenericResource("3", "exampleAPIReceiver", resources.LambdaType)
	kinesisResource   = resources.NewGenericResource("4", "MyKinesis", resources.KinesisType)
	lambdaResource    = resources.NewGenericResource("5", "exampleReceiver", resources.LambdaType)
	cronResource      = resources.NewGenericResource("6", "cron(0 1 * * ? *)", resources.CronType)
	sourceSQSResource = resources.NewGenericResource("7", "source", resources.SQSType)

	wantResourceCollection = &resources.ResourceCollection{
		Resources: []resources.Resource{
			endpointResource,
			apiGateway,
			lambdaAPIGateway,
			kinesisResource,
			lambdaResource,
			cronResource,
			sourceSQSResource,
			resources.NewGenericResource("8", "MyApi", resources.RestfulAPIType),
			resources.NewGenericResource("9", "my-bucket", resources.S3Type),
			resources.NewGenericResource("10", "target", resources.SQSType),
			resources.NewGenericResource("11", "example", resources.SNSType),
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
				Source: cronResource,
				Target: lambdaResource,
			},
			{
				Source: kinesisResource,
				Target: lambdaResource,
			},
			{
				Source: sourceSQSResource,
				Target: lambdaResource,
			},
		},
	}
)

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
			require.Equal(t, tc.want, got)
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
		resources     *[]resources.Resource
		relationships *[]resources.Relationship
		id            *int
	}

	lambdaResource := resources.NewGenericResource("1", "myReceiver", resources.LambdaType)

	var (
		id             int
		targetResource resources.Resource
	)

	tests := []struct {
		name   string
		fields fields
		args   args
		setup  func()
	}{
		{
			name: "lambda and sqs",
			fields: fields{
				yamlConfig: &config.Config{},
			},
			args: args{
				res: &config.Lambda{Envars: []map[string]string{{
					"DOCDB_HOST": "var.doc_db_host",
				}}},
				lambda:        lambdaResource,
				resources:     &[]resources.Resource{},
				relationships: &[]resources.Relationship{},
				id:            &id,
			},
			setup: func() {
				id = 2
				targetResource = resources.NewGenericResource("2", "doc", resources.DatabaseType)
			},
		},
		{
			name: "lambda and google BQ",
			fields: fields{
				yamlConfig: &config.Config{},
			},
			args: args{
				res: &config.Lambda{Envars: []map[string]string{{
					"GOOGLE_BQ_PROJECT_ID": "google",
				}}},
				lambda:        lambdaResource,
				resources:     &[]resources.Resource{},
				relationships: &[]resources.Relationship{},
				id:            &id,
			},
			setup: func() {
				id = 2
				targetResource = resources.NewGenericResource("2", "google", resources.GoogleBQType)
			},
		},
		{
			name: "lambda and kinesis",
			fields: fields{
				yamlConfig: &config.Config{},
			},
			args: args{
				res: &config.Lambda{Envars: []map[string]string{{
					"MY_KINESIS_KINESIS_STREAM_URL": "MyKinesis",
				}}},
				lambda:        lambdaResource,
				resources:     &[]resources.Resource{},
				relationships: &[]resources.Relationship{},
				id:            &id,
			},
			setup: func() {
				id = 2
				targetResource = resources.NewGenericResource("2", "MyKinesis", resources.KinesisType)
			},
		},
		{
			name: "lambda and s3 bucket S3_BUCKET",
			fields: fields{
				yamlConfig: &config.Config{},
			},
			args: args{
				res: &config.Lambda{Envars: []map[string]string{{
					"PAYLOADS_S3_BUCKET": "payloads",
				}}},
				lambda:        lambdaResource,
				resources:     &[]resources.Resource{},
				relationships: &[]resources.Relationship{},
				id:            &id,
			},
			setup: func() {
				id = 2
				targetResource = resources.NewGenericResource("2", "payloads", resources.S3Type)
			},
		},
		{
			name: "lambda and s3 bucket BUCKET_NAME",
			fields: fields{
				yamlConfig: &config.Config{},
			},
			args: args{
				res: &config.Lambda{Envars: []map[string]string{{
					"PAYLOADS_BUCKET_NAME": "payloads",
				}}},
				lambda:        lambdaResource,
				resources:     &[]resources.Resource{},
				relationships: &[]resources.Relationship{},
				id:            &id,
			},
			setup: func() {
				id = 2
				targetResource = resources.NewGenericResource("2", "payloads", resources.S3Type)
			},
		},
		{
			name: "lambda and restful API",
			fields: fields{
				yamlConfig: &config.Config{},
			},
			args: args{
				res: &config.Lambda{Envars: []map[string]string{{
					"MY_REST_API_BASE_URL": "MyRest",
				}}},
				lambda:        lambdaResource,
				resources:     &[]resources.Resource{},
				relationships: &[]resources.Relationship{},
				id:            &id,
			},
			setup: func() {
				id = 2
				targetResource = resources.NewGenericResource("2", "MyRest", resources.RestfulAPIType)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			tr := NewTransformer(tt.fields.yamlConfig)

			tr.transformLambdaEnvars(tt.args.res, tt.args.lambda, tt.args.resources, tt.args.relationships, tt.args.id)

			require.Equal(t, 3, *tt.args.id)
			require.Equal(t, &[]resources.Resource{targetResource}, tt.args.resources)
			require.Equal(t, &[]resources.Relationship{{
				Source: lambdaResource,
				Target: targetResource,
			}}, tt.args.relationships)
		})
	}
}
