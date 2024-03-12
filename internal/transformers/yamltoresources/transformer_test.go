package yamltoresources

import (
	_ "embed"
	"testing"

	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
	"github.com/joselitofilho/aws-terraform-generator/internal/resources"
	"gopkg.in/yaml.v3"

	"github.com/stretchr/testify/require"
)

//go:embed testdata/diagram.yaml
var diagramData []byte

func TestTransformer_Transform(t *testing.T) {
	type fields struct {
		yamlConfig *config.Config
	}

	var diagramYAML *config.Config

	err := yaml.Unmarshal(diagramData, &diagramYAML)
	require.NoError(t, err)

	lambdaResource := resources.NewGenericResource("4", "exampleReceiver", resources.LambdaType)

	tests := []struct {
		name      string
		fields    fields
		want      *resources.ResourceCollection
		targetErr error
	}{
		{
			name:   "happy path",
			fields: fields{yamlConfig: diagramYAML},
			want: &resources.ResourceCollection{
				Resources: []resources.Resource{
					resources.NewGenericResource("1", "POST /v1/examples", resources.APIGatewayType),
					resources.NewGenericResource("2", "mystack-api.domain-${var.environment}.com",
						resources.EndpointType),
					resources.NewGenericResource("3", "myKinesis", resources.KinesisType),
					lambdaResource,
					resources.NewGenericResource("8", "MyAPI", resources.RestfulAPIType),
					resources.NewGenericResource("9", "my-bucket", resources.S3Type),
					resources.NewGenericResource("10", "target", resources.SQSType),
					resources.NewGenericResource("11", "source", resources.SQSType),
					resources.NewGenericResource("12", "example", resources.SNSType),
				},
				Relationships: []resources.Relationship{
					{
						Source: resources.NewGenericResource("5", "cron(0 1 * * ? *)", resources.CronType),
						Target: lambdaResource,
					},
					{
						Source: resources.NewGenericResource("6", "aws_kinesis_stream.mykinesis_kinesis.arn",
							resources.KinesisType),
						Target: lambdaResource,
					},
					{
						Source: resources.NewGenericResource("7", "aws_sqs_queue.source_sqs.arn", resources.SQSType),
						Target: lambdaResource,
					},
				},
			},
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
