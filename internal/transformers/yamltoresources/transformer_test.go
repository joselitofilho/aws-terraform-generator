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

var (
	apiGateway       = resources.NewGenericResource("1", "POST /v1/examples", resources.APIGatewayType)
	lambdaAPIGateway = resources.NewGenericResource("2", "exampleAPIReceiver", resources.LambdaType)
	lambdaResource   = resources.NewGenericResource("5", "exampleReceiver", resources.LambdaType)

	wantResourceCollection = &resources.ResourceCollection{
		Resources: []resources.Resource{
			apiGateway,
			lambdaAPIGateway,
			resources.NewGenericResource("3", "mystack-api.domain-${var.environment}.com",
				resources.EndpointType),
			resources.NewGenericResource("4", "myKinesis", resources.KinesisType),
			lambdaResource,
			resources.NewGenericResource("9", "MyAPI", resources.RestfulAPIType),
			resources.NewGenericResource("10", "my-bucket", resources.S3Type),
			resources.NewGenericResource("11", "target", resources.SQSType),
			resources.NewGenericResource("12", "source", resources.SQSType),
			resources.NewGenericResource("13", "example", resources.SNSType),
		},
		Relationships: []resources.Relationship{
			{
				Source: apiGateway,
				Target: lambdaAPIGateway,
			},
			{
				Source: resources.NewGenericResource("6", "cron(0 1 * * ? *)", resources.CronType),
				Target: lambdaResource,
			},
			{
				Source: resources.NewGenericResource("7", "aws_kinesis_stream.mykinesis_kinesis.arn",
					resources.KinesisType),
				Target: lambdaResource,
			},
			{
				Source: resources.NewGenericResource("8", "aws_sqs_queue.source_sqs.arn", resources.SQSType),
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
