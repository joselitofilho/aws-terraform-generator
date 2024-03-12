package yamltoresources

import (
	"testing"

	"github.com/joselitofilho/aws-terraform-generator/internal/resources"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	type args struct {
		filename string
	}

	lambdaResource := resources.NewGenericResource("4", "exampleReceiver", resources.LambdaType)

	tests := []struct {
		name      string
		args      args
		want      *resources.ResourceCollection
		targetErr error
	}{
		{
			name: "happy path",
			args: args{filename: "testdata/diagram.yaml"},
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
	}

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			got, err := Parse(tc.args.filename)

			require.ErrorIs(t, err, tc.targetErr)
			require.Equal(t, tc.want, got)
		})
	}
}
