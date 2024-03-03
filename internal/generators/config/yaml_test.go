package config

import (
	"errors"
	"os"
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/stretchr/testify/require"
)

var (
	testdataFolder = "../testdata"

	errDummy = errors.New("dummy error")
)

func TestYAML_Parse(t *testing.T) {
	type fields struct {
		fileName string
	}

	tests := []struct {
		setup     func(testing.TB) func(testing.TB)
		name      string
		fields    fields
		want      *Config
		targetErr error
	}{
		{
			setup:  func(_ testing.TB) func(testing.TB) { return func(_ testing.TB) {} },
			name:   "API Gateway",
			fields: fields{fileName: testdataFolder + "/apigateway.config.yaml"},
			want: &Config{APIGateways: []APIGateway{{
				StackName: "teststack",
				APIDomain: "teststack-api.domain-${var.environment}.com",
				APIG:      true,
				Lambdas: []APIGatewayLambda{{
					Name:        "exampleAPIReceiver",
					Source:      "git@github.com:username/terraform-aws-lambda?ref=reference",
					RoleName:    "execute_lambda",
					Runtime:     "go1.x",
					Description: "Trigger the example API receiver via API Gateway",
					Verb:        "POST",
					Path:        "/v1/examples",
					Envars:      []map[string]string{{"MYVAR": "MYVAR_VALUE"}},
					Files: []File{
						{
							Name:    "lambda.go",
							Imports: []string{"context", "github.com/logging"},
							Tmpl:    "package main\n\nimport (\n  {{ range getFileImports $.Files \"lambda.go\" }}\"{{ . }}\"\n  {{end}}\n)",
						},
						{Name: "main.go", Tmpl: "package main"},
					},
				}},
			}}},
		},
		{
			setup:  func(_ testing.TB) func(testing.TB) { return func(_ testing.TB) {} },
			name:   "Diagram",
			fields: fields{fileName: testdataFolder + "/diagram.config.yaml"},
			want: &Config{Diagram: Diagram{
				StackName: "teststack",
				Lambda:    DriagramLambda{Source: "git@", RoleName: "execute_lambda", Runtime: "go1.x"},
			}},
		},
		{
			setup:  func(_ testing.TB) func(testing.TB) { return func(_ testing.TB) {} },
			name:   "Lambda",
			fields: fields{fileName: testdataFolder + "/lambda.config.yaml"},
			want: &Config{Lambdas: []Lambda{{
				Name:            "exampleReceiver",
				Source:          "git@github.com:username/terraform-aws-lambda?ref=reference",
				RoleName:        "execute_lambda",
				Runtime:         "go1.x",
				Description:     "Trigger on schedule and initiate the execution of example receiver",
				Envars:          []map[string]string{{"SQS_QUEUE_URL": "aws_sqs_queue.target_sqs.name"}},
				KinesisTriggers: []KinesisTrigger{{SourceARN: "aws_kinesis_stream.source_mykinesis_kinesis.arn"}},
				SQSTriggers:     []SQSTrigger{{SourceARN: "aws_sqs_queue.source_sqs.arn"}},
				Crons:           []Cron{{ScheduleExpression: "cron(0 1 * * ? *)", IsEnabled: "var.trigger_enabled"}},
				Files: []File{{
					Name:    "lambda.go",
					Imports: []string{"github.com/mylogging/logging"},
					Tmpl:    `package main`,
				}},
			}}},
		},
		{
			setup:  func(_ testing.TB) func(testing.TB) { return func(_ testing.TB) {} },
			name:   "RestfulAPI",
			fields: fields{fileName: testdataFolder + "/restfulapi.config.yaml"},
			want:   &Config{RestfulAPIs: []RestfulAPI{{Name: "MyAPI"}}},
		},
		{
			setup:  func(_ testing.TB) func(testing.TB) { return func(_ testing.TB) {} },
			name:   "S3",
			fields: fields{fileName: testdataFolder + "/s3.config.yaml"},
			want: &Config{Buckets: []S3{{Name: "my-bucket", ExpirationDays: 90, Files: []File{{
				Name: "my-bucket-s3.tf", Tmpl: `resource "aws_s3_bucket" "{{ToSnake $.Name}}_bucket" {}`,
			}}}}},
		},
		{
			setup:  func(_ testing.TB) func(testing.TB) { return func(_ testing.TB) {} },
			name:   "SNS",
			fields: fields{fileName: testdataFolder + "/sns.config.yaml"},
			want: &Config{SNSs: []SNS{{
				Name:       "example",
				BucketName: "my-bucket",
				Lambdas: []SNSResource{{
					Name:         "exampleReceiver",
					Events:       []string{"s3:ObjectCreated:*"},
					FilterPrefix: "my_prefix",
					FilterSuffix: ".txt",
				}},
				SQSs: []SNSResource{{
					Name:         "target",
					Events:       []string{"s3:ObjectCreated:*"},
					FilterPrefix: "my_prefix",
					FilterSuffix: ".txt",
				}},
				Files: []File{{
					Name: "example-sns.tf",
					Tmpl: `resource "aws_s3_bucket_notification" "s3_bucket_notification_{{ToSnake $.Name}}" {}`,
				}},
			}}},
		},
		{
			setup:  func(_ testing.TB) func(testing.TB) { return func(_ testing.TB) {} },
			name:   "SQS",
			fields: fields{fileName: testdataFolder + "/sqs.config.yaml"},
			want: &Config{SQSs: []SQS{
				{
					Name:            "target",
					MaxReceiveCount: 15,
					Files: []File{{
						Name: "target-sqs.tf", Tmpl: `resource "aws_sqs_queue" "{{ToSnake $.Name}}_sqs" {}`,
					}}},
				{Name: "source", MaxReceiveCount: 10},
			}},
		},
		{
			setup:  func(_ testing.TB) func(testing.TB) { return func(_ testing.TB) {} },
			name:   "Structure",
			fields: fields{fileName: testdataFolder + "/structure.config.yaml"},
			want: &Config{Structure: Structure{
				Stacks: []Stack{{
					Name: "teststack",
					Folders: []Folder{
						{
							Name:  "dev",
							Files: []File{{Name: "main.tf"}, {Name: "terragrunt.hcl"}, {Name: "vars.tf"}},
						},
						{
							Name:  "uat",
							Files: []File{{Name: "main.tf"}, {Name: "terragrunt.hcl"}, {Name: "vars.tf"}},
						},
						{
							Name:  "prd",
							Files: []File{{Name: "main.tf"}, {Name: "terragrunt.hcl"}, {Name: "vars.tf"}},
						},
						{
							Name: "mod",
							Files: []File{
								{Name: "main.tf", Tmpl: `locals { stack_name = "{{$.StackName}}-${var.environment}" }`},
								{Name: "vars.tf"},
							},
						},
						{
							Name: "lambda",
						},
					},
				}},
				DefaultTemplates: []FilenameTemplateMap{{"main.tf": `# Terraform backend and required providers configuration`}},
			}},
		},
		{
			setup: func(_ testing.TB) func(testing.TB) {
				osReadFile = func(name string) ([]byte, error) {
					require.Empty(t, name)
					return nil, errDummy
				}

				return func(_ testing.TB) {
					osReadFile = os.ReadFile
				}
			},
			name:      "Empty File",
			fields:    fields{fileName: ""},
			targetErr: errDummy,
		},
		{
			setup: func(_ testing.TB) func(testing.TB) {
				yamlUnmarshal = func(_ []byte, _ any) error {
					return errDummy
				}
				return func(_ testing.TB) {
					yamlUnmarshal = yaml.Unmarshal
				}
			},
			name:      "Invalid YAML Syntax",
			fields:    fields{fileName: testdataFolder + "/invalid_sintax.yaml"},
			targetErr: errDummy,
		},
	}

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			tearDown := tc.setup(t)
			defer tearDown(t)

			got, err := NewYAML(tc.fields.fileName).Parse()

			require.ErrorIs(t, err, tc.targetErr)
			require.Equal(t, tc.want, got)
		})
	}
}
