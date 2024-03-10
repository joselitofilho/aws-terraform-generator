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
			name: "from lambda to kinesis",
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
										"KINESIS_STREAM_URL": "aws_kinesis_stream.my_stream_kinesis.name",
									},
								},
							},
						},
					},
				},
			},
			want: &resources.ResourceCollection{
				Resources: []resources.Resource{lambdaResource, kinesisResource},
				Relationships: []resources.Relationship{{
					Source: lambdaResource, Target: kinesisResource,
				}},
			},
		},
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
