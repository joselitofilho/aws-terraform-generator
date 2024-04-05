package resources

import (
	"testing"

	resources "github.com/diagram-code-generator/resources/pkg/resources"
	"github.com/stretchr/testify/require"
)

func TestAWSResourceFactory_CreateResource(t *testing.T) {
	type args struct {
		id    string
		value string
		style string
	}

	tests := []struct {
		name string
		args args
		want resources.Resource
	}{
		{
			name: "API Gateway Resource",
			args: args{
				id:    "APIG_ID",
				value: "myAPI",
				style: "mxgraph.aws3.api_gateway",
			},
			want: resources.NewGenericResource("APIG_ID", "myAPI", APIGatewayType.String()),
		},
		{
			name: "Cron Resource",
			args: args{
				id:    "CRON_ID",
				value: "myScheduler",
				style: "mxgraph.aws4.event_time_based",
			},
			want: resources.NewGenericResource("CRON_ID", "myScheduler", CronType.String()),
		},
		{
			name: "Database Resource",
			args: args{
				id:    "DB_ID",
				value: "myDB",
				style: "mxgraph.flowchart.database",
			},
			want: resources.NewGenericResource("DB_ID", "myDB", DatabaseType.String()),
		},
		{
			name: "Endpoint Resource",
			args: args{
				id:    "ENDPOINT_ID",
				value: "myEndpoint",
				style: "mxgraph.aws4.endpoint",
			},
			want: resources.NewGenericResource("ENDPOINT_ID", "myEndpoint", EndpointType.String()),
		},
		{
			name: "GoogleBQ Resource",
			args: args{
				id:    "GBC_ID",
				value: "myGBC",
				style: "google_bigquery",
			},
			want: resources.NewGenericResource("GBC_ID", "myGBC", GoogleBQType.String()),
		},
		{
			name: "Kinesis Resource",
			args: args{
				id:    "KINESIS_ID",
				value: "myKinesis",
				style: "mxgraph.aws3.kinesis",
			},
			want: resources.NewGenericResource("KINESIS_ID", "myKinesis", KinesisType.String()),
		},
		{
			name: "Lambda Resource",
			args: args{
				id:    "LAMBDA_ID",
				value: "myReceiver",
				style: "mxgraph.aws3.lambda",
			},
			want: resources.NewGenericResource("LAMBDA_ID", "myReceiver", LambdaType.String()),
		},
		{
			name: "Restful API Resource",
			args: args{
				id:    "RESTFULAPI_ID",
				value: "myRestAPI",
				style: "mxgraph.veeam2.restful_api",
			},
			want: resources.NewGenericResource("RESTFULAPI_ID", "myRestAPI", RestfulAPIType.String()),
		},
		{
			name: "S3 Resource",
			args: args{
				id:    "S3BUCKET_ID",
				value: "myBucket",
				style: "mxgraph.aws3.s3",
			},
			want: resources.NewGenericResource("S3BUCKET_ID", "myBucket", S3Type.String()),
		},
		{
			name: "SQS Resource",
			args: args{
				id:    "SQS_ID",
				value: "my-sqs",
				style: "mxgraph.aws3.sqs",
			},
			want: resources.NewGenericResource("SQS_ID", "my-sqs", SQSType.String()),
		},
		{
			name: "SNS Resource",
			args: args{
				id:    "SNS_ID",
				value: "my-sns",
				style: "mxgraph.aws3.sns",
			},
			want: resources.NewGenericResource("SNS_ID", "my-sns", SNSType.String()),
		},
		{
			name: "Unknown",
			args: args{
				id:    "ID",
				value: "value",
				style: "any",
			},
			want: nil,
		},
	}

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			f := &AWSResourceFactory{}
			got := f.CreateResource(tc.args.id, tc.args.value, tc.args.style)

			require.Equal(t, tc.want, got)
		})
	}
}
