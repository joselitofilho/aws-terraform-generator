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

func TestTransformDrawIOToYAML_Lambda(t *testing.T) {
	type args struct {
		yamlConfig *config.Config
		resources  *drawio.ResourceCollection
	}

	tests := []struct {
		name      string
		args      args
		want      *config.Config
		targetErr error
	}{
		{
			name: "happy path - lambda as module",
			args: args{
				yamlConfig: diagramConfig,
				resources: &drawio.ResourceCollection{
					Resources: []drawio.Resource{
						drawio.NewGenericResource("id1", "myReceiver", drawio.LambdaType),
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
						Files: []config.File{
							{Name: "lambda.go"}, {Name: "main.go"},
						},
					},
				},
			},
		},
		{
			name: "happy path - lambda as resource",
			args: args{
				yamlConfig: &config.Config{
					Diagram: config.Diagram{
						Lambda: config.DriagramLambda{
							Source:   "../artifacts",
							RoleName: "execute_lambda",
							Runtime:  "go1.x",
						},
					},
				},
				resources: &drawio.ResourceCollection{
					Resources: []drawio.Resource{
						drawio.NewGenericResource("id1", "myReceiver", drawio.LambdaType),
					},
				},
			},
			want: &config.Config{
				Lambdas: []config.Lambda{
					{
						Name:        "myReceiver",
						Source:      "../artifacts",
						RoleName:    "execute_lambda",
						Runtime:     "go1.x",
						Description: "myReceiver lambda",
						Files: []config.File{
							{Name: "lambda.go"}, {Name: "main.go"},
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
			name: "happy path - api gateway",
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
		{
			name: "when api gateway doen't have any endpoint or lambda associated should return a valid config",
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

	tests := []struct {
		name      string
		args      args
		want      *config.Config
		targetErr error
	}{
		{
			name: "happy path - s3 bucket",
			args: args{
				yamlConfig: diagramConfig,
				resources: &drawio.ResourceCollection{
					Resources: []drawio.Resource{
						drawio.NewGenericResource("id1", "my-bucket", drawio.S3Type),
					},
				},
			},
			want: &config.Config{
				Buckets: []config.S3{
					{
						Name:           "my-bucket",
						ExpirationDays: 90,
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

func TestTransformDrawIOToYAML_SQS(t *testing.T) {
	type args struct {
		yamlConfig *config.Config
		resources  *drawio.ResourceCollection
	}

	tests := []struct {
		name      string
		args      args
		want      *config.Config
		targetErr error
	}{
		{
			name: "happy path - sqs",
			args: args{
				yamlConfig: diagramConfig,
				resources: &drawio.ResourceCollection{
					Resources: []drawio.Resource{
						drawio.NewGenericResource("id1", "my-queue", drawio.SQSType),
					},
				},
			},
			want: &config.Config{
				SQSs: []config.SQS{
					{
						Name:            "my-queue",
						MaxReceiveCount: 10,
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

func TestTransformDrawIOToYAML_SNS(t *testing.T) {
	type args struct {
		yamlConfig *config.Config
		resources  *drawio.ResourceCollection
	}

	tests := []struct {
		name      string
		args      args
		want      *config.Config
		targetErr error
	}{
		{
			name: "happy path - sns",
			args: args{
				yamlConfig: diagramConfig,
				resources: &drawio.ResourceCollection{
					Resources: []drawio.Resource{
						drawio.NewGenericResource("id1", "my-sns", drawio.SNSType),
					},
				},
			},
			want: &config.Config{SNSs: []config.SNS{{Name: "my-sns"}}},
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

	tests := []struct {
		name      string
		args      args
		want      *config.Config
		targetErr error
	}{
		{
			name: "happy path - restful api",
			args: args{
				yamlConfig: diagramConfig,
				resources: &drawio.ResourceCollection{
					Resources: []drawio.Resource{
						drawio.NewGenericResource("id1", "my-api", drawio.RestfulAPIType),
					},
				},
			},
			want: &config.Config{RestfulAPIs: []config.RestfulAPI{{Name: "my-api"}}},
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
