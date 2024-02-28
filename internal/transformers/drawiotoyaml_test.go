package transformers

import (
	"testing"

	"github.com/joselitofilho/aws-terraform-generator/internal/drawio"
	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"

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
		resources  *drawio.ResourceCollection
	}

	endpointResource := drawio.NewGenericResource("id1", "https://my-domain.com", drawio.EndpointType)
	apiGatewayResource := drawio.NewGenericResource("id2", "POST /examples", drawio.APIGatewayType)
	lambdaResource := drawio.NewGenericResource("id3", "my-lambda", drawio.LambdaType)

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
				resources: &drawio.ResourceCollection{
					Resources: []drawio.Resource{apiGatewayResource},
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
				resources: &drawio.ResourceCollection{
					Resources: []drawio.Resource{
						endpointResource,
						apiGatewayResource,
						lambdaResource,
					},
					Relationships: []drawio.Relationship{
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
								Files: []config.File{
									{Name: "lambda.go"}, {Name: "main.go"},
								},
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
		resources  *drawio.ResourceCollection
	}

	database := drawio.NewGenericResource("id1", "myDatabase", drawio.DatabaseType)
	lambda := drawio.NewGenericResource("id2", "myReceiver", drawio.LambdaType)

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
				resources:  &drawio.ResourceCollection{Resources: []drawio.Resource{database}},
			},
			want: &config.Config{},
		},
		{
			name: "database receives data from a Lambda",
			args: args{
				yamlConfig: diagramConfig,
				resources: &drawio.ResourceCollection{
					Resources:     []drawio.Resource{database, lambda},
					Relationships: []drawio.Relationship{{Source: lambda, Target: database}},
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
							{"MYDATABASEDB_HOST": "var.mydatabase_db_host"},
							{"MYDATABASEDB_USER": "var.mydatabase_db_user"},
							{"MYDATABASEDB_PASSWORD_SECRET": "var.mydatabase_db_password_secret"},
						},
						Files: []config.File{{Name: "lambda.go"}, {Name: "main.go"}},
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
		resources  *drawio.ResourceCollection
	}

	kinesis := drawio.NewGenericResource("id1", "my-stream", drawio.KinesisType)

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
				resources: &drawio.ResourceCollection{
					Resources: []drawio.Resource{kinesis},
				},
			},
			want: &config.Config{
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
		resources  *drawio.ResourceCollection
	}

	lambda := drawio.NewGenericResource("id1", "myReceiver", drawio.LambdaType)
	cron := drawio.NewGenericResource("id2", "cron(0 2 * * ? *)", drawio.CronType)
	sqs := drawio.NewGenericResource("id3", "my-queue", drawio.SQSType)
	sns := drawio.NewGenericResource("id4", "my-notification", drawio.SNSType)
	s3Bucket := drawio.NewGenericResource("id5", "my-bucket", drawio.S3Type)
	kinesis := drawio.NewGenericResource("id6", "my-stream", drawio.KinesisType)

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
				resources: &drawio.ResourceCollection{
					Resources: []drawio.Resource{lambda},
				},
			},
			want: &config.Config{
				Lambdas: []config.Lambda{
					{
						Name:        "myReceiver",
						Source:      "git@",
						RoleName:    "execute_lambda",
						Description: "myReceiver lambda",
						Files:       []config.File{{Name: "lambda.go"}, {Name: "main.go"}},
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
				resources: &drawio.ResourceCollection{
					Resources: []drawio.Resource{lambda},
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
						Files:       []config.File{{Name: "lambda.go"}, {Name: "main.go"}},
					},
				},
			},
		},
		{
			name: "invoke a Lambda on a schedule using a cron expression",
			args: args{
				yamlConfig: diagramConfig,
				resources: &drawio.ResourceCollection{
					Resources:     []drawio.Resource{lambda, cron},
					Relationships: []drawio.Relationship{{Source: cron, Target: lambda}},
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
						Files:       []config.File{{Name: "lambda.go"}, {Name: "main.go"}},
					},
				},
			},
		},
		{
			name: "invoke a Lambda to receive messages from an Kinesis stream",
			args: args{
				yamlConfig: diagramConfig,
				resources: &drawio.ResourceCollection{
					Resources:     []drawio.Resource{lambda, kinesis},
					Relationships: []drawio.Relationship{{Source: kinesis, Target: lambda}},
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
						Files:           []config.File{{Name: "lambda.go"}, {Name: "main.go"}},
					},
				},
				Kinesis: []config.Kinesis{{Name: "my-stream", RetentionPeriod: "24"}},
			},
		},
		{
			name: "invoke a Lambda to receive messages from an SQS queue",
			args: args{
				yamlConfig: diagramConfig,
				resources: &drawio.ResourceCollection{
					Resources:     []drawio.Resource{lambda, sqs},
					Relationships: []drawio.Relationship{{Source: sqs, Target: lambda}},
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
						Files:       []config.File{{Name: "lambda.go"}, {Name: "main.go"}},
					},
				},
				SQSs: []config.SQS{{Name: "my-queue", MaxReceiveCount: 10}},
			},
		},
		{
			name: "invoke a Lambda to process SNS notifications",
			args: args{
				yamlConfig: diagramConfig,
				resources: &drawio.ResourceCollection{
					Resources: []drawio.Resource{lambda, sns, s3Bucket},
					Relationships: []drawio.Relationship{
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
						Files:       []config.File{{Name: "lambda.go"}, {Name: "main.go"}},
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

			if tc.targetErr == nil {
				require.NoError(t, err)
				require.Equal(t, tc.want, got)
			} else {
				require.ErrorIs(t, err, tc.targetErr)
			}
		})
	}
}

func TestTransformDrawIOToYAML_S3Bucket(t *testing.T) {
	type args struct {
		yamlConfig *config.Config
		resources  *drawio.ResourceCollection
	}

	s3Bucket := drawio.NewGenericResource("id1", "my-bucket", drawio.S3Type)
	lambda := drawio.NewGenericResource("id2", "myReceiver", drawio.LambdaType)

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
				resources: &drawio.ResourceCollection{
					Resources: []drawio.Resource{s3Bucket},
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
				resources: &drawio.ResourceCollection{
					Resources:     []drawio.Resource{s3Bucket, lambda},
					Relationships: []drawio.Relationship{{Source: lambda, Target: s3Bucket}},
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
						Files: []config.File{{Name: "lambda.go"}, {Name: "main.go"}},
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
		resources  *drawio.ResourceCollection
	}

	sqs := drawio.NewGenericResource("id1", "my-queue", drawio.SQSType)
	lambda := drawio.NewGenericResource("id2", "myReceiver", drawio.LambdaType)
	sns := drawio.NewGenericResource("id3", "my-notification", drawio.SNSType)
	s3Bucket := drawio.NewGenericResource("id4", "my-bucket", drawio.S3Type)

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
				resources: &drawio.ResourceCollection{
					Resources: []drawio.Resource{sqs},
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
				resources: &drawio.ResourceCollection{
					Resources:     []drawio.Resource{sqs, lambda},
					Relationships: []drawio.Relationship{{Source: lambda, Target: sqs}},
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
							{"MY_QUEUE_SQS_QUEUE_URL": "aws_sqs_queue.my_queue_sqs.id"},
						},
						Files: []config.File{{Name: "lambda.go"}, {Name: "main.go"}},
					},
				},
				SQSs: []config.SQS{{Name: "my-queue", MaxReceiveCount: 10}},
			},
		},
		{
			name: "enqueue a message to an SQS queue triggered by SNS notification",
			args: args{
				yamlConfig: diagramConfig,
				resources: &drawio.ResourceCollection{
					Resources: []drawio.Resource{sqs, sns, s3Bucket},
					Relationships: []drawio.Relationship{
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
		resources  *drawio.ResourceCollection
	}

	sns := drawio.NewGenericResource("id1", "my-notification", drawio.SNSType)
	s3Bucket := drawio.NewGenericResource("id1", "my-bucket", drawio.S3Type)

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
				resources:  &drawio.ResourceCollection{Resources: []drawio.Resource{sns}},
			},
			want: &config.Config{SNSs: []config.SNS{{Name: "my-notification"}}},
		},
		{
			name: "receive an SNS notification triggered by an S3 Bucket event",
			args: args{
				yamlConfig: diagramConfig,
				resources: &drawio.ResourceCollection{
					Resources:     []drawio.Resource{sns, s3Bucket},
					Relationships: []drawio.Relationship{{Source: s3Bucket, Target: sns}},
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
		resources  *drawio.ResourceCollection
	}

	restfulAPI := drawio.NewGenericResource("id1", "my-api", drawio.RestfulAPIType)
	lambda := drawio.NewGenericResource("id2", "myReceiver", drawio.LambdaType)

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
				resources:  &drawio.ResourceCollection{Resources: []drawio.Resource{restfulAPI}},
			},
			want: &config.Config{RestfulAPIs: []config.RestfulAPI{{Name: "my-api"}}},
		},
		{
			name: "restful API receives a request from a Lambda",
			args: args{
				yamlConfig: diagramConfig,
				resources: &drawio.ResourceCollection{
					Resources:     []drawio.Resource{restfulAPI, lambda},
					Relationships: []drawio.Relationship{{Source: lambda, Target: restfulAPI}},
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
						Files: []config.File{{Name: "lambda.go"}, {Name: "main.go"}},
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
