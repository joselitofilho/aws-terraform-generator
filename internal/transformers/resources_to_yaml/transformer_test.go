package resources_to_yaml

import (
	"testing"

	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
	"github.com/joselitofilho/aws-terraform-generator/internal/resources"

	"github.com/stretchr/testify/require"
)

var diagramConfig = &config.Config{
	Diagram: config.Diagram{
		StackName: "my-stack",
		Lambda: config.DriagramLambda{
			Source:   "git@",
			RoleName: "execute_lambda",
		},
	},
}

func TestTransformDrawIOToYAML_APIGateway(t *testing.T) {
	type args struct {
		yamlConfig *config.Config
		resources  *resources.ResourceCollection
	}

	endpointResource := resources.NewGenericResource("id1", "https://my-domain.com", resources.EndpointType)
	apiGatewayResource := resources.NewGenericResource("id2", "POST /examples", resources.APIGatewayType)
	lambdaResource := resources.NewGenericResource("id3", "my-lambda", resources.LambdaType)

	tests := []struct {
		name      string
		args      args
		want      *config.Config
		targetErr error
	}{
		{
			name: "only API Gateway",
			args: args{
				yamlConfig: diagramConfig,
				resources: &resources.ResourceCollection{
					Resources: []resources.Resource{apiGatewayResource},
				},
			},
			want: &config.Config{
				APIGateways: []config.APIGateway{
					{
						StackName: "my-stack",
						APIG:      true,
					},
				},
			},
		},
		{
			name: "API Gateway full example",
			args: args{
				yamlConfig: diagramConfig,
				resources: &resources.ResourceCollection{
					Resources: []resources.Resource{
						endpointResource,
						apiGatewayResource,
						lambdaResource,
					},
					Relationships: []resources.Relationship{
						{Source: endpointResource, Target: apiGatewayResource},
						{Source: apiGatewayResource, Target: lambdaResource},
					},
				},
			},
			want: &config.Config{
				APIGateways: []config.APIGateway{
					{
						StackName: "my-stack",
						APIG:      true,
						APIDomain: "https://my-domain.com",
						Lambdas: []config.APIGatewayLambda{
							{
								Name:        "my-lambda",
								Source:      "git@",
								RoleName:    "execute_lambda",
								Description: "my-lambda lambda",
								Verb:        "POST",
								Path:        "/examples",
							},
						},
					},
				},
			},
		},
	}

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			got, err := TransformDrawIOToYAML(tc.args.yamlConfig, tc.args.resources)

			if tc.targetErr == nil {
				require.NoError(t, err)
				require.Equal(t, tc.want, got)
			} else {
				require.ErrorIs(t, err, tc.targetErr)
			}
		})
	}
}

func TestTransformDrawIOToYAML_Database(t *testing.T) {
	type args struct {
		yamlConfig *config.Config
		resources  *resources.ResourceCollection
	}

	database := resources.NewGenericResource("id1", "my-database", resources.DatabaseType)
	lambda := resources.NewGenericResource("id2", "myReceiver", resources.LambdaType)

	tests := []struct {
		name      string
		args      args
		want      *config.Config
		targetErr error
	}{
		{
			name: "only database",
			args: args{
				yamlConfig: diagramConfig,
				resources:  &resources.ResourceCollection{Resources: []resources.Resource{database}},
			},
			want: &config.Config{},
		},
		{
			name: "database receives data from a Lambda",
			args: args{
				yamlConfig: diagramConfig,
				resources: &resources.ResourceCollection{
					Resources:     []resources.Resource{database, lambda},
					Relationships: []resources.Relationship{{Source: lambda, Target: database}},
				},
			},
			want: &config.Config{
				Lambdas: []config.Lambda{
					{
						Name:        "myReceiver",
						Source:      "git@",
						RoleName:    "execute_lambda",
						Description: "myReceiver lambda",
						Envars: []map[string]string{
							{"MY_DATABASE_DB_HOST": "var.my_database_db_host"},
							{"MY_DATABASE_DB_USER": "var.my_database_db_user"},
							{"MY_DATABASE_DB_PASSWORD_SECRET": "var.my_database_db_password_secret"},
						},
					},
				},
			},
		},
	}

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			got, err := TransformDrawIOToYAML(tc.args.yamlConfig, tc.args.resources)

			if tc.targetErr == nil {
				require.NoError(t, err)
				require.Equal(t, tc.want, got)
			} else {
				require.ErrorIs(t, err, tc.targetErr)
			}
		})
	}
}

func TestTransformDrawIOToYAML_GoogleBQ(t *testing.T) {
	type args struct {
		yamlConfig *config.Config
		resources  *resources.ResourceCollection
	}

	googleBQ := resources.NewGenericResource("id1", "google", resources.GoogleBQType)
	lambda := resources.NewGenericResource("id2", "myReceiver", resources.LambdaType)

	tests := []struct {
		name      string
		args      args
		want      *config.Config
		targetErr error
	}{
		{
			name: "only google BQ",
			args: args{
				yamlConfig: diagramConfig,
				resources:  &resources.ResourceCollection{Resources: []resources.Resource{googleBQ}},
			},
			want: &config.Config{},
		},
		{
			name: "google BQ receives data from a Lambda",
			args: args{
				yamlConfig: diagramConfig,
				resources: &resources.ResourceCollection{
					Resources:     []resources.Resource{googleBQ, lambda},
					Relationships: []resources.Relationship{{Source: lambda, Target: googleBQ}},
				},
			},
			want: &config.Config{
				Lambdas: []config.Lambda{
					{
						Name:        "myReceiver",
						Source:      "git@",
						RoleName:    "execute_lambda",
						Description: "myReceiver lambda",
						Envars: []map[string]string{
							{"GOOGLE_BQ_PROJECT_ID": "var.google_bq_project_id"},
							{"GOOGLE_BQ_API_KEY_SECRET": "var.google_bq_api_key_secret"},
							{"GOOGLE_BQ_PARTITION_FIELD": "var.google_bq_partition_field"},
							{"GOOGLE_BQ_CLUSTERING_FIELDS": "var.google_bq_clustering_fields"},
						},
					},
				},
			},
		},
	}

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			got, err := TransformDrawIOToYAML(tc.args.yamlConfig, tc.args.resources)

			if tc.targetErr == nil {
				require.NoError(t, err)
				require.Equal(t, tc.want, got)
			} else {
				require.ErrorIs(t, err, tc.targetErr)
			}
		})
	}
}

func TestTransformDrawIOToYAML_Kinesis(t *testing.T) {
	type args struct {
		yamlConfig *config.Config
		resources  *resources.ResourceCollection
	}

	kinesis := resources.NewGenericResource("id1", "my-stream", resources.KinesisType)
	lambda := resources.NewGenericResource("id2", "myReceiver", resources.LambdaType)

	tests := []struct {
		name      string
		args      args
		want      *config.Config
		targetErr error
	}{
		{
			name: "only Kinesis",
			args: args{
				yamlConfig: diagramConfig,
				resources: &resources.ResourceCollection{
					Resources: []resources.Resource{kinesis},
				},
			},
			want: &config.Config{
				Kinesis: []config.Kinesis{{Name: "my-stream", RetentionPeriod: "24"}},
			},
		},
		{
			name: "enqueue a message to an Kinesis stream published by Lambda",
			args: args{
				yamlConfig: diagramConfig,
				resources: &resources.ResourceCollection{
					Resources:     []resources.Resource{kinesis, lambda},
					Relationships: []resources.Relationship{{Source: lambda, Target: kinesis}},
				},
			},
			want: &config.Config{
				Lambdas: []config.Lambda{
					{
						Name:        "myReceiver",
						Source:      "git@",
						RoleName:    "execute_lambda",
						Description: "myReceiver lambda",
						Envars: []map[string]string{
							{"MY_STREAM_KINESIS_STREAM_URL": "aws_kinesis_stream.my_stream_kinesis.name"},
						},
					},
				},
				Kinesis: []config.Kinesis{{Name: "my-stream", RetentionPeriod: "24"}},
			},
		},
	}

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			got, err := TransformDrawIOToYAML(tc.args.yamlConfig, tc.args.resources)

			if tc.targetErr == nil {
				require.NoError(t, err)
				require.Equal(t, tc.want, got)
			} else {
				require.ErrorIs(t, err, tc.targetErr)
			}
		})
	}
}

func TestTransformDrawIOToYAML_Lambda(t *testing.T) {
	type args struct {
		yamlConfig *config.Config
		resources  *resources.ResourceCollection
	}

	lambda := resources.NewGenericResource("id1", "myReceiver", resources.LambdaType)
	cron := resources.NewGenericResource("id2", "cron(0 2 * * ? *)", resources.CronType)
	sqs := resources.NewGenericResource("id3", "my-queue", resources.SQSType)
	sns := resources.NewGenericResource("id4", "my-notification", resources.SNSType)
	s3Bucket := resources.NewGenericResource("id5", "my-bucket", resources.S3Type)
	kinesis := resources.NewGenericResource("id6", "my-stream", resources.KinesisType)

	tests := []struct {
		name      string
		args      args
		want      *config.Config
		targetErr error
	}{
		{
			name: "Lambda as module",
			args: args{
				yamlConfig: diagramConfig,
				resources: &resources.ResourceCollection{
					Resources: []resources.Resource{lambda},
				},
			},
			want: &config.Config{
				Lambdas: []config.Lambda{
					{
						Name:        "myReceiver",
						Source:      "git@",
						RoleName:    "execute_lambda",
						Description: "myReceiver lambda",
					},
				},
			},
		},
		{
			name: "Lambda as resource",
			args: args{
				yamlConfig: &config.Config{
					Diagram: config.Diagram{
						Lambda: config.DriagramLambda{
							Source:   "../artefacts",
							RoleName: "execute_lambda",
							Runtime:  "go1.x",
						},
					},
				},
				resources: &resources.ResourceCollection{
					Resources: []resources.Resource{lambda},
				},
			},
			want: &config.Config{
				Lambdas: []config.Lambda{
					{
						Name:        "myReceiver",
						Source:      "../artefacts",
						RoleName:    "execute_lambda",
						Runtime:     "go1.x",
						Description: "myReceiver lambda",
					},
				},
			},
		},
		{
			name: "invoke a Lambda on a schedule using a cron expression",
			args: args{
				yamlConfig: diagramConfig,
				resources: &resources.ResourceCollection{
					Resources:     []resources.Resource{lambda, cron},
					Relationships: []resources.Relationship{{Source: cron, Target: lambda}},
				},
			},
			want: &config.Config{
				Lambdas: []config.Lambda{
					{
						Name:        "myReceiver",
						Source:      "git@",
						RoleName:    "execute_lambda",
						Description: "myReceiver lambda",
						Crons:       []config.Cron{{ScheduleExpression: "cron(0 2 * * ? *)", IsEnabled: "true"}},
					},
				},
			},
		},
		{
			name: "invoke a Lambda to receive messages from an Kinesis stream",
			args: args{
				yamlConfig: diagramConfig,
				resources: &resources.ResourceCollection{
					Resources:     []resources.Resource{lambda, kinesis},
					Relationships: []resources.Relationship{{Source: kinesis, Target: lambda}},
				},
			},
			want: &config.Config{
				Lambdas: []config.Lambda{
					{
						Name:            "myReceiver",
						Source:          "git@",
						RoleName:        "execute_lambda",
						Description:     "myReceiver lambda",
						KinesisTriggers: []config.KinesisTrigger{{SourceARN: "aws_kinesis_stream.my_stream_kinesis.arn"}},
					},
				},
				Kinesis: []config.Kinesis{{Name: "my-stream", RetentionPeriod: "24"}},
			},
		},
		{
			name: "invoke a Lambda to receive messages from an SQS queue",
			args: args{
				yamlConfig: diagramConfig,
				resources: &resources.ResourceCollection{
					Resources:     []resources.Resource{lambda, sqs},
					Relationships: []resources.Relationship{{Source: sqs, Target: lambda}},
				},
			},
			want: &config.Config{
				Lambdas: []config.Lambda{
					{
						Name:        "myReceiver",
						Source:      "git@",
						RoleName:    "execute_lambda",
						Description: "myReceiver lambda",
						SQSTriggers: []config.SQSTrigger{{SourceARN: "aws_sqs_queue.my_queue_sqs.arn"}},
					},
				},
				SQSs: []config.SQS{{Name: "my-queue", MaxReceiveCount: 10}},
			},
		},
		{
			name: "invoke a Lambda to process SNS notifications",
			args: args{
				yamlConfig: diagramConfig,
				resources: &resources.ResourceCollection{
					Resources: []resources.Resource{lambda, sns, s3Bucket},
					Relationships: []resources.Relationship{
						{Source: s3Bucket, Target: sns},
						{Source: sns, Target: lambda},
					},
				},
			},
			want: &config.Config{
				Lambdas: []config.Lambda{
					{
						Name:        "myReceiver",
						Source:      "git@",
						RoleName:    "execute_lambda",
						Description: "myReceiver lambda",
					},
				},
				SNSs: []config.SNS{{
					Name:       "my-notification",
					BucketName: "my-bucket",
					Lambdas:    []config.SNSResource{{Name: "myReceiver", Events: []string{"s3:ObjectCreated:*"}}},
				}},
				Buckets: []config.S3{{Name: "my-bucket", ExpirationDays: 90}},
			},
		},
	}

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			got, err := TransformDrawIOToYAML(tc.args.yamlConfig, tc.args.resources)

			require.ErrorIs(t, err, tc.targetErr)
			require.Equal(t, tc.want, got)
		})
	}
}

func TestTransformDrawIOToYAML_S3Bucket(t *testing.T) {
	type args struct {
		yamlConfig *config.Config
		resources  *resources.ResourceCollection
	}

	s3Bucket := resources.NewGenericResource("id1", "my-bucket", resources.S3Type)
	lambda := resources.NewGenericResource("id2", "myReceiver", resources.LambdaType)

	tests := []struct {
		name      string
		args      args
		want      *config.Config
		targetErr error
	}{
		{
			name: "only s3 bucket",
			args: args{
				yamlConfig: diagramConfig,
				resources: &resources.ResourceCollection{
					Resources: []resources.Resource{s3Bucket},
				},
			},
			want: &config.Config{
				Buckets: []config.S3{{Name: "my-bucket", ExpirationDays: 90}},
			},
		},
		{
			name: "s3 bucket stores object from Lambda",
			args: args{
				yamlConfig: diagramConfig,
				resources: &resources.ResourceCollection{
					Resources:     []resources.Resource{s3Bucket, lambda},
					Relationships: []resources.Relationship{{Source: lambda, Target: s3Bucket}},
				},
			},
			want: &config.Config{
				Lambdas: []config.Lambda{
					{
						Name:        "myReceiver",
						Source:      "git@",
						RoleName:    "execute_lambda",
						Description: "myReceiver lambda",
						Envars: []map[string]string{
							{"MY_BUCKET_S3_BUCKET": "aws_s3_bucket.my_bucket_bucket.bucket"},
							{"MY_BUCKET_S3_DIRECTORY": `"my_receiver_files"`},
						},
					},
				},
				Buckets: []config.S3{{Name: "my-bucket", ExpirationDays: 90}},
			},
		},
	}

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			got, err := TransformDrawIOToYAML(tc.args.yamlConfig, tc.args.resources)

			if tc.targetErr == nil {
				require.NoError(t, err)
				require.Equal(t, tc.want, got)
			} else {
				require.ErrorIs(t, err, tc.targetErr)
			}
		})
	}
}

func TestTransformDrawIOToYAML_SQS(t *testing.T) {
	type args struct {
		yamlConfig *config.Config
		resources  *resources.ResourceCollection
	}

	sqs := resources.NewGenericResource("id1", "my-queue", resources.SQSType)
	lambda := resources.NewGenericResource("id2", "myReceiver", resources.LambdaType)
	sns := resources.NewGenericResource("id3", "my-notification", resources.SNSType)
	s3Bucket := resources.NewGenericResource("id4", "my-bucket", resources.S3Type)

	tests := []struct {
		name      string
		args      args
		want      *config.Config
		targetErr error
	}{
		{
			name: "only SQS",
			args: args{
				yamlConfig: diagramConfig,
				resources: &resources.ResourceCollection{
					Resources: []resources.Resource{sqs},
				},
			},
			want: &config.Config{
				SQSs: []config.SQS{{Name: "my-queue", MaxReceiveCount: 10}},
			},
		},
		{
			name: "enqueue a message to an SQS queue published by Lambda",
			args: args{
				yamlConfig: diagramConfig,
				resources: &resources.ResourceCollection{
					Resources:     []resources.Resource{sqs, lambda},
					Relationships: []resources.Relationship{{Source: lambda, Target: sqs}},
				},
			},
			want: &config.Config{
				Lambdas: []config.Lambda{
					{
						Name:        "myReceiver",
						Source:      "git@",
						RoleName:    "execute_lambda",
						Description: "myReceiver lambda",
						Envars: []map[string]string{
							{"MY_QUEUE_SQS_QUEUE_URL": "aws_sqs_queue.my_queue_sqs.name"},
						},
					},
				},
				SQSs: []config.SQS{{Name: "my-queue", MaxReceiveCount: 10}},
			},
		},
		{
			name: "enqueue a message to an SQS queue triggered by SNS notification",
			args: args{
				yamlConfig: diagramConfig,
				resources: &resources.ResourceCollection{
					Resources: []resources.Resource{sqs, sns, s3Bucket},
					Relationships: []resources.Relationship{
						{Source: s3Bucket, Target: sns},
						{Source: sns, Target: sqs},
					},
				},
			},
			want: &config.Config{
				SQSs: []config.SQS{{Name: "my-queue", MaxReceiveCount: 10}},
				SNSs: []config.SNS{{
					Name:       "my-notification",
					BucketName: "my-bucket",
					SQSs:       []config.SNSResource{{Name: "my-queue", Events: []string{"s3:ObjectCreated:*"}}},
				}},
				Buckets: []config.S3{{Name: "my-bucket", ExpirationDays: 90}},
			},
		},
	}

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			got, err := TransformDrawIOToYAML(tc.args.yamlConfig, tc.args.resources)

			if tc.targetErr == nil {
				require.NoError(t, err)
				require.Equal(t, tc.want, got)
			} else {
				require.ErrorIs(t, err, tc.targetErr)
			}
		})
	}
}

func TestTransformDrawIOToYAML_SNS(t *testing.T) {
	type args struct {
		yamlConfig *config.Config
		resources  *resources.ResourceCollection
	}

	sns := resources.NewGenericResource("id1", "my-notification", resources.SNSType)
	s3Bucket := resources.NewGenericResource("id1", "my-bucket", resources.S3Type)

	tests := []struct {
		name      string
		args      args
		want      *config.Config
		targetErr error
	}{
		{
			name: "only SNS",
			args: args{
				yamlConfig: diagramConfig,
				resources:  &resources.ResourceCollection{Resources: []resources.Resource{sns}},
			},
			want: &config.Config{SNSs: []config.SNS{{Name: "my-notification"}}},
		},
		{
			name: "receive an SNS notification triggered by an S3 Bucket event",
			args: args{
				yamlConfig: diagramConfig,
				resources: &resources.ResourceCollection{
					Resources:     []resources.Resource{sns, s3Bucket},
					Relationships: []resources.Relationship{{Source: s3Bucket, Target: sns}},
				},
			},
			want: &config.Config{
				SNSs: []config.SNS{{
					Name:       "my-notification",
					BucketName: "my-bucket",
				}},
				Buckets: []config.S3{{Name: "my-bucket", ExpirationDays: 90}},
			},
		},
	}

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			got, err := TransformDrawIOToYAML(tc.args.yamlConfig, tc.args.resources)

			if tc.targetErr == nil {
				require.NoError(t, err)
				require.Equal(t, tc.want, got)
			} else {
				require.ErrorIs(t, err, tc.targetErr)
			}
		})
	}
}

func TestTransformDrawIOToYAML_RestfulAPI(t *testing.T) {
	type args struct {
		yamlConfig *config.Config
		resources  *resources.ResourceCollection
	}

	restfulAPI := resources.NewGenericResource("id1", "my-api", resources.RestfulAPIType)
	lambda := resources.NewGenericResource("id2", "myReceiver", resources.LambdaType)

	tests := []struct {
		name      string
		args      args
		want      *config.Config
		targetErr error
	}{
		{
			name: "only restful API",
			args: args{
				yamlConfig: diagramConfig,
				resources:  &resources.ResourceCollection{Resources: []resources.Resource{restfulAPI}},
			},
			want: &config.Config{RestfulAPIs: []config.RestfulAPI{{Name: "my-api"}}},
		},
		{
			name: "restful API receives a request from a Lambda",
			args: args{
				yamlConfig: diagramConfig,
				resources: &resources.ResourceCollection{
					Resources:     []resources.Resource{restfulAPI, lambda},
					Relationships: []resources.Relationship{{Source: lambda, Target: restfulAPI}},
				},
			},
			want: &config.Config{
				Lambdas: []config.Lambda{
					{
						Name:        "myReceiver",
						Source:      "git@",
						RoleName:    "execute_lambda",
						Description: "myReceiver lambda",
						Envars: []map[string]string{
							{"MY_API_API_BASE_URL": "var.my_api_api_base_url"},
							{"MY_API_HOST": "var.my_api_host"},
							{"MY_API_USER": "var.my_api_user"},
						},
					},
				},
				RestfulAPIs: []config.RestfulAPI{{Name: "my-api"}},
			},
		},
	}

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			got, err := TransformDrawIOToYAML(tc.args.yamlConfig, tc.args.resources)

			if tc.targetErr == nil {
				require.NoError(t, err)
				require.Equal(t, tc.want, got)
			} else {
				require.ErrorIs(t, err, tc.targetErr)
			}
		})
	}
}
