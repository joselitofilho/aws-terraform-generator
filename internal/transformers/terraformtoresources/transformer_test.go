package terraformtoresources

import (
	"testing"

	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
	"github.com/joselitofilho/aws-terraform-generator/internal/generators/terraform"
	"github.com/joselitofilho/aws-terraform-generator/internal/resources"
	"github.com/stretchr/testify/require"
)

func TestTransformer_Transform(t *testing.T) {
	type fields struct {
		yamlConfig *config.Config
		tfConfig   *terraform.Config
	}

	lambdaResource := resources.NewGenericResource("1", "exampleReceiver", resources.LambdaType)

	tests := []struct {
		name   string
		fields fields
		want   *resources.ResourceCollection
	}{

		{
			name: "empty",
			fields: fields{
				yamlConfig: &config.Config{},
				tfConfig:   &terraform.Config{},
			},
			want: &resources.ResourceCollection{
				Resources:     []resources.Resource{},
				Relationships: []resources.Relationship{},
			},
		},
		{
			name: "API Gateway route",
			fields: fields{
				yamlConfig: &config.Config{},
				tfConfig: &terraform.Config{
					Resources: []*terraform.Resource{
						{
							Type:   "aws_apigatewayv2_route",
							Name:   "apigw_route_example_api_receiver",
							Labels: []string{"aws_apigatewayv2_route", "apigw_route_example_api_receiver"},
							Attributes: map[string]any{
								"api_id":    "aws_apigatewayv2_api.mystack_api.id",
								"route_key": "POST /v1/examples",
								"target":    "integrations/${aws_apigatewayv2_integration.example_api_receiver.id}",
							},
						},
					},
				},
			},
			want: &resources.ResourceCollection{
				Resources: []resources.Resource{
					resources.NewGenericResource("1", "POST /v1/examples", resources.APIGatewayType)},
				Relationships: []resources.Relationship{},
			},
		},
		{
			name: "API Gateway integration",
			fields: fields{
				yamlConfig: &config.Config{},
				tfConfig: &terraform.Config{
					Resources: []*terraform.Resource{
						{
							Type:   "aws_apigatewayv2_integration",
							Name:   "example_api_receiver",
							Labels: []string{"aws_apigatewayv2_integration", "example_api_receiver"},
							Attributes: map[string]any{
								"api_id":          "aws_apigatewayv2_api.mystack_api.id",
								"integration_uri": "aws_lambda_function.example_api_receiver_lambda.invoke_arn",
							},
						},
					},
				},
			},
			want: &resources.ResourceCollection{
				Resources:     []resources.Resource{},
				Relationships: []resources.Relationship{},
			},
		},
		{
			name: "cloudwatch event rune",
			fields: fields{
				yamlConfig: &config.Config{},
				tfConfig: &terraform.Config{
					Resources: []*terraform.Resource{
						{
							Type:   "aws_cloudwatch_event_rule",
							Name:   "example_receiver_cron_rule",
							Labels: []string{"aws_cloudwatch_event_rule", "example_receiver_cron_rule"},
							Attributes: map[string]any{
								"schedule_expression": "cron()",
							},
						},
					},
				},
			},
			want: &resources.ResourceCollection{
				Resources:     []resources.Resource{resources.NewGenericResource("1", "cron()", resources.CronType)},
				Relationships: []resources.Relationship{},
			},
		},
		{
			name: "cron",
			fields: fields{
				yamlConfig: &config.Config{},
				tfConfig: &terraform.Config{
					Resources: []*terraform.Resource{
						{
							Type:   "aws_cloudwatch_event_target",
							Name:   "example_receiver_cron_target",
							Labels: []string{"aws_cloudwatch_event_target", "example_receiver_cron_target"},
							Attributes: map[string]any{
								"rule": "aws_cloudwatch_event_rule.example_receiver_cron.name",
								"arn":  "aws_lambda_function.example_receiver_lambda.arn",
							},
						},
					},
				},
			},
			want: &resources.ResourceCollection{
				Resources:     []resources.Resource{},
				Relationships: []resources.Relationship{},
			},
		},
		{
			name: "endpoint",
			fields: fields{
				yamlConfig: &config.Config{},
				tfConfig: &terraform.Config{
					Resources: []*terraform.Resource{
						{
							Type:   "aws_apigatewayv2_domain_name",
							Name:   "my_restful_api",
							Labels: []string{"aws_apigatewayv2_domain_name", "my_restful_api"},
							Attributes: map[string]any{
								"domain_name": "local.api_domain",
							},
						},
					},
				},
			},
			want: &resources.ResourceCollection{
				Resources:     []resources.Resource{resources.NewGenericResource("1", "local.api_domain", resources.EndpointType)},
				Relationships: []resources.Relationship{},
			},
		},
		{
			name: "kinesis",
			fields: fields{
				yamlConfig: &config.Config{},
				tfConfig: &terraform.Config{
					Resources: []*terraform.Resource{
						{
							Type:   "aws_kinesis_stream",
							Name:   "my_stream",
							Labels: []string{"aws_kinesis_stream", "my_stream"},
							Attributes: map[string]any{
								"name": "MyStream",
							},
						},
					},
				},
			},
			want: &resources.ResourceCollection{
				Resources: []resources.Resource{
					resources.NewGenericResource("1", "MyStream", resources.KinesisType)},
				Relationships: []resources.Relationship{},
			},
		},
		{
			name: "lambda event source mapping with kinesis",
			fields: fields{
				yamlConfig: &config.Config{},
				tfConfig: &terraform.Config{
					Resources: []*terraform.Resource{
						{
							Type:   "aws_lambda_event_source_mapping",
							Name:   "example_checker_lambda_sqs_trigger",
							Labels: []string{"aws_lambda_event_source_mapping", "example_checker_lambda_sqs_trigger"},
							Attributes: map[string]any{
								"event_source_arn": "aws_kinesis_stream.my_stream_kinesis.arn",
								"function_name":    "aws_lambda_function.example_checker_lambda.arn",
							},
						},
					},
				},
			},
			want: &resources.ResourceCollection{
				Resources: []resources.Resource{
					resources.NewGenericResource("1", "MyStream", resources.KinesisType)},
				Relationships: []resources.Relationship{},
			},
		},
		// {
		// 	name: "lambda event source mapping with sqs", // TODO: Implement
		// },
		{
			name: "lambda as resource",
			fields: fields{
				yamlConfig: &config.Config{},
				tfConfig: &terraform.Config{
					Resources: []*terraform.Resource{
						{
							Type:   "aws_lambda_function",
							Name:   "example_receiver_lambda",
							Labels: []string{"aws_lambda_function", "example_receiver_lambda"},
							Attributes: map[string]any{
								"function_name": "example_receiver_lambda",
							},
						},
					},
				},
			},
			want: &resources.ResourceCollection{
				Resources:     []resources.Resource{lambdaResource},
				Relationships: []resources.Relationship{},
			},
		},
		{
			name: "lambda as module",
			fields: fields{
				yamlConfig: &config.Config{},
				tfConfig: &terraform.Config{
					Modules: []*terraform.Module{{Labels: []string{"example_receiver_lambda"}}},
				},
			},
			want: &resources.ResourceCollection{
				Resources:     []resources.Resource{lambdaResource},
				Relationships: []resources.Relationship{},
			},
		},
		{
			name: "s3 bucket",
			fields: fields{
				yamlConfig: &config.Config{},
				tfConfig: &terraform.Config{
					Resources: []*terraform.Resource{
						{
							Type:   "aws_s3_bucket",
							Name:   "my_bucket",
							Labels: []string{"aws_s3_bucket", "my_bucket"},
							Attributes: map[string]any{
								"bucket": "var.client-var.environment-my-bucket",
							},
						},
					},
				},
			},
			want: &resources.ResourceCollection{
				Resources: []resources.Resource{
					resources.NewGenericResource("1", "my-bucket", resources.S3Type)},
				Relationships: []resources.Relationship{},
			},
		},
		{
			name: "sqs",
			fields: fields{
				yamlConfig: &config.Config{},
				tfConfig: &terraform.Config{
					Resources: []*terraform.Resource{
						{
							Type:   "aws_sqs_queue",
							Name:   "my_sqs",
							Labels: []string{"aws_sqs_queue", "my_sqs"},
							Attributes: map[string]any{
								"name": "my-queue",
							},
						},
					},
				},
			},
			want: &resources.ResourceCollection{
				Resources: []resources.Resource{
					resources.NewGenericResource("1", "my-queue", resources.SQSType)},
				Relationships: []resources.Relationship{},
			},
		},
	}

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			tr := NewTransformer(
				tc.fields.yamlConfig,
				tc.fields.tfConfig,
			)

			got := tr.Transform()

			require.Equal(t, tc.want, got)
		})
	}
}

func TestTransformer_TransformFromLambdaToDatabase(t *testing.T) {
	type fields struct {
		yamlConfig *config.Config
		tfConfig   *terraform.Config
	}

	lambdaResource := resources.NewGenericResource("1", "myReceiver", resources.LambdaType)
	dbResource := resources.NewGenericResource("2", "doc", resources.DatabaseType)

	tests := []struct {
		name   string
		fields fields
		want   *resources.ResourceCollection
	}{
		{
			name: "lambda as resource",
			fields: fields{
				yamlConfig: &config.Config{},
				tfConfig: &terraform.Config{
					Resources: []*terraform.Resource{
						{
							Type:   "aws_lambda_function",
							Name:   "my_receiver_lambda",
							Labels: []string{"aws_lambda_function", "my_receiver_lambda"},
							Attributes: map[string]any{
								"function_name": "my_receiver_lambda",
								"environment": map[string]map[string]any{
									"variables": {
										"DOCDB_HOST": "var.doc_db_host",
									},
								},
							},
						},
					},
				},
			},
			want: &resources.ResourceCollection{
				Resources:     []resources.Resource{lambdaResource, dbResource},
				Relationships: []resources.Relationship{{Source: lambdaResource, Target: dbResource}},
			},
		},
		{
			name: "lambda as module with resources from envar",
			fields: fields{
				yamlConfig: &config.Config{},
				tfConfig: &terraform.Config{
					Modules: []*terraform.Module{
						{
							Labels: []string{"my_receiver_lambda"},
							Attributes: map[string]any{
								"function_name": "my_receiver_lambda",
								"lambda_function_env_vars": map[string]any{
									"DOCDB_HOST": "var.doc_db_host",
								},
							},
						},
					},
				},
			},
			want: &resources.ResourceCollection{
				Resources:     []resources.Resource{lambdaResource, dbResource},
				Relationships: []resources.Relationship{{Source: lambdaResource, Target: dbResource}},
			},
		},
	}

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			tr := NewTransformer(
				tc.fields.yamlConfig,
				tc.fields.tfConfig,
			)

			got := tr.Transform()

			require.Equal(t, tc.want, got)
		})
	}
}

func TestTransformer_TransformFromLambdaToGoogleBQ(t *testing.T) {
	type fields struct {
		yamlConfig *config.Config
		tfConfig   *terraform.Config
	}

	lambdaResource := resources.NewGenericResource("1", "myReceiver", resources.LambdaType)
	bqResource := resources.NewGenericResource("2", "google", resources.GoogleBQType)

	tests := []struct {
		name   string
		fields fields
		want   *resources.ResourceCollection
	}{
		{
			name: "lambda as resource",
			fields: fields{
				yamlConfig: &config.Config{},
				tfConfig: &terraform.Config{
					Resources: []*terraform.Resource{
						{
							Type:   "aws_lambda_function",
							Name:   "my_receiver_lambda",
							Labels: []string{"aws_lambda_function", "my_receiver_lambda"},
							Attributes: map[string]any{
								"function_name": "my_receiver_lambda",
								"environment": map[string]map[string]any{
									"variables": {
										"GOOGLE_BQ_PROJECT_ID": "google",
									},
								},
							},
						},
					},
				},
			},
			want: &resources.ResourceCollection{
				Resources:     []resources.Resource{lambdaResource, bqResource},
				Relationships: []resources.Relationship{{Source: lambdaResource, Target: bqResource}},
			},
		},
		{
			name: "lambda as module with resources from envar",
			fields: fields{
				yamlConfig: &config.Config{},
				tfConfig: &terraform.Config{
					Modules: []*terraform.Module{
						{
							Labels: []string{"my_receiver_lambda"},
							Attributes: map[string]any{
								"function_name": "my_receiver_lambda",
								"lambda_function_env_vars": map[string]any{
									"BQ_PROJECT_ID": "google",
								},
							},
						},
					},
				},
			},
			want: &resources.ResourceCollection{
				Resources:     []resources.Resource{lambdaResource, bqResource},
				Relationships: []resources.Relationship{{Source: lambdaResource, Target: bqResource}},
			},
		},
	}

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			tr := NewTransformer(
				tc.fields.yamlConfig,
				tc.fields.tfConfig,
			)

			got := tr.Transform()

			require.Equal(t, tc.want, got)
		})
	}
}

func TestTransformer_TransformFromLambdaToKinesis(t *testing.T) {
	type fields struct {
		yamlConfig *config.Config
		tfConfig   *terraform.Config
	}

	lambdaResource := resources.NewGenericResource("1", "myReceiver", resources.LambdaType)
	kinesisResource := resources.NewGenericResource("2", "MyStream", resources.KinesisType)

	tests := []struct {
		name   string
		fields fields
		want   *resources.ResourceCollection
	}{
		{
			name: "lambda as resource",
			fields: fields{
				yamlConfig: &config.Config{},
				tfConfig: &terraform.Config{
					Resources: []*terraform.Resource{
						{
							Type:   "aws_lambda_function",
							Name:   "my_receiver_lambda",
							Labels: []string{"aws_lambda_function", "my_receiver_lambda"},
							Attributes: map[string]any{
								"function_name": "my_receiver_lambda",
								"environment": map[string]map[string]any{
									"variables": {
										"MY_STREAM_KINESIS_STREAM_URL": "aws_kinesis_stream.my_stream_kinesis.name",
									},
								},
							},
						},
					},
				},
			},
			want: &resources.ResourceCollection{
				Resources:     []resources.Resource{lambdaResource, kinesisResource},
				Relationships: []resources.Relationship{{Source: lambdaResource, Target: kinesisResource}},
			},
		},
		{
			name: "lambda as module with resources from envar",
			fields: fields{
				yamlConfig: &config.Config{},
				tfConfig: &terraform.Config{
					Modules: []*terraform.Module{
						{
							Labels: []string{"my_receiver_lambda"},
							Attributes: map[string]any{
								"function_name": "my_receiver_lambda",
								"lambda_function_env_vars": map[string]any{
									"KINESIS_STREAM_URL": "aws_kinesis_stream.my_stream_kinesis.name",
								},
							},
						},
					},
				},
			},
			want: &resources.ResourceCollection{
				Resources:     []resources.Resource{lambdaResource, kinesisResource},
				Relationships: []resources.Relationship{{Source: lambdaResource, Target: kinesisResource}},
			},
		},
	}

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			tr := NewTransformer(
				tc.fields.yamlConfig,
				tc.fields.tfConfig,
			)

			got := tr.Transform()

			require.Equal(t, tc.want, got)
		})
	}
}

func TestTransformer_TransformFromLambdaToRestfulAPI(t *testing.T) {
	type fields struct {
		yamlConfig *config.Config
		tfConfig   *terraform.Config
	}

	lambdaResource := resources.NewGenericResource("1", "myReceiver", resources.LambdaType)
	restfulAPIResource := resources.NewGenericResource("2", "myRestful", resources.RestfulAPIType)

	tests := []struct {
		name   string
		fields fields
		want   *resources.ResourceCollection
	}{
		{
			name: "lambda as resource",
			fields: fields{
				yamlConfig: &config.Config{},
				tfConfig: &terraform.Config{
					Resources: []*terraform.Resource{
						{
							Type:   "aws_lambda_function",
							Name:   "my_receiver_lambda",
							Labels: []string{"aws_lambda_function", "my_receiver_lambda"},
							Attributes: map[string]any{
								"function_name": "my_receiver_lambda",
								"environment": map[string]map[string]any{
									"variables": {
										"MY_RESTFUL_API_BASE_URL": "myRestful",
									},
								},
							},
						},
					},
				},
			},
			want: &resources.ResourceCollection{
				Resources:     []resources.Resource{lambdaResource, restfulAPIResource},
				Relationships: []resources.Relationship{{Source: lambdaResource, Target: restfulAPIResource}},
			},
		},
		{
			name: "lambda as module with resources from envar",
			fields: fields{
				yamlConfig: &config.Config{},
				tfConfig: &terraform.Config{
					Modules: []*terraform.Module{
						{
							Labels: []string{"my_receiver_lambda"},
							Attributes: map[string]any{
								"function_name": "my_receiver_lambda",
								"lambda_function_env_vars": map[string]any{
									"API_BASE_URL": "myRestful",
								},
							},
						},
					},
				},
			},
			want: &resources.ResourceCollection{
				Resources:     []resources.Resource{lambdaResource, restfulAPIResource},
				Relationships: []resources.Relationship{{Source: lambdaResource, Target: restfulAPIResource}},
			},
		},
	}

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			tr := NewTransformer(
				tc.fields.yamlConfig,
				tc.fields.tfConfig,
			)

			got := tr.Transform()

			require.Equal(t, tc.want, got)
		})
	}
}

func TestTransformer_TransformFromLambdaToS3Bucket(t *testing.T) {
	type fields struct {
		yamlConfig *config.Config
		tfConfig   *terraform.Config
	}

	lambdaResource := resources.NewGenericResource("1", "myReceiver", resources.LambdaType)
	s3BucketResource := resources.NewGenericResource("2", "my-bucket", resources.S3Type)

	tests := []struct {
		name   string
		fields fields
		want   *resources.ResourceCollection
	}{
		{
			name: "lambda as resource",
			fields: fields{
				yamlConfig: &config.Config{},
				tfConfig: &terraform.Config{
					Resources: []*terraform.Resource{
						{
							Type:   "aws_lambda_function",
							Name:   "my_receiver_lambda",
							Labels: []string{"aws_lambda_function", "my_receiver_lambda"},
							Attributes: map[string]any{
								"function_name": "my_receiver_lambda",
								"environment": map[string]map[string]any{
									"variables": {
										"MY_BUCKET_S3_BUCKET": "my-bucket",
									},
								},
							},
						},
					},
				},
			},
			want: &resources.ResourceCollection{
				Resources:     []resources.Resource{lambdaResource, s3BucketResource},
				Relationships: []resources.Relationship{{Source: lambdaResource, Target: s3BucketResource}},
			},
		},
		{
			name: "lambda as module with resources from envar",
			fields: fields{
				yamlConfig: &config.Config{},
				tfConfig: &terraform.Config{
					Modules: []*terraform.Module{
						{
							Labels: []string{"my_receiver_lambda"},
							Attributes: map[string]any{
								"function_name": "my_receiver_lambda",
								"lambda_function_env_vars": map[string]any{
									"S3_BUCKET": "my-bucket",
								},
							},
						},
					},
				},
			},
			want: &resources.ResourceCollection{
				Resources:     []resources.Resource{lambdaResource, s3BucketResource},
				Relationships: []resources.Relationship{{Source: lambdaResource, Target: s3BucketResource}},
			},
		},
	}

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			tr := NewTransformer(
				tc.fields.yamlConfig,
				tc.fields.tfConfig,
			)

			got := tr.Transform()

			require.Equal(t, tc.want, got)
		})
	}
}

func TestTransformer_TransformFromLambdaToSQS(t *testing.T) {
	type fields struct {
		yamlConfig *config.Config
		tfConfig   *terraform.Config
	}

	lambdaResource := resources.NewGenericResource("1", "myReceiver", resources.LambdaType)
	sqsResource := resources.NewGenericResource("2", "my-queue", resources.SQSType)

	tests := []struct {
		name   string
		fields fields
		want   *resources.ResourceCollection
	}{
		{
			name: "lambda as resource",
			fields: fields{
				yamlConfig: &config.Config{},
				tfConfig: &terraform.Config{
					Resources: []*terraform.Resource{
						{
							Type:   "aws_lambda_function",
							Name:   "my_receiver_lambda",
							Labels: []string{"aws_lambda_function", "my_receiver_lambda"},
							Attributes: map[string]any{
								"function_name": "my_receiver_lambda",
								"environment": map[string]map[string]any{
									"variables": {
										"MY_QUEUE_SQS_QUEUE_URL": "my-queue",
									},
								},
							},
						},
					},
				},
			},
			want: &resources.ResourceCollection{
				Resources:     []resources.Resource{lambdaResource, sqsResource},
				Relationships: []resources.Relationship{{Source: lambdaResource, Target: sqsResource}},
			},
		},
		{
			name: "lambda as module with resources from envar",
			fields: fields{
				yamlConfig: &config.Config{},
				tfConfig: &terraform.Config{
					Modules: []*terraform.Module{
						{
							Labels: []string{"my_receiver_lambda"},
							Attributes: map[string]any{
								"function_name": "my_receiver_lambda",
								"lambda_function_env_vars": map[string]any{
									"SQS_QUEUE_URL": "my-queue",
								},
							},
						},
					},
				},
			},
			want: &resources.ResourceCollection{
				Resources:     []resources.Resource{lambdaResource, sqsResource},
				Relationships: []resources.Relationship{{Source: lambdaResource, Target: sqsResource}},
			},
		},
	}

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			tr := NewTransformer(
				tc.fields.yamlConfig,
				tc.fields.tfConfig,
			)

			got := tr.Transform()

			require.Equal(t, tc.want, got)
		})
	}
}
