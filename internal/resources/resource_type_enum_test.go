package resources

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResourceType_String(t *testing.T) {
	tests := []struct {
		name string
		rt   ResourceType
		want string
	}{
		{name: "APIGateway", rt: APIGatewayType, want: "APIGateway"},
		{name: "Cron", rt: CronType, want: "Cron"},
		{name: "Database", rt: DatabaseType, want: "Database"},
		{name: "Endpoint", rt: EndpointType, want: "Endpoint"},
		{name: "GoogleBQ", rt: GoogleBQType, want: "GoogleBQ"},
		{name: "Kinesis", rt: KinesisType, want: "Kinesis"},
		{name: "Lambda", rt: LambdaType, want: "Lambda"},
		{name: "RestfulAPI", rt: RestfulAPIType, want: "RestfulAPI"},
		{name: "S3", rt: S3Type, want: "S3"},
		{name: "SQS", rt: SQSType, want: "SQS"},
		{name: "SNS", rt: SNSType, want: "SNS"},
		{name: "Unknown", rt: "", want: "Unknown"},
	}

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			got := tc.rt.String()

			require.Equal(t, tc.want, got)
		})
	}
}

func TestParseResourceType(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output ResourceType
	}{
		{name: "Parse APIGateway", input: "APIGateway", output: APIGatewayType},
		{name: "Parse Cron", input: "Cron", output: CronType},
		{name: "Parse Database", input: "Database", output: DatabaseType},
		{name: "Parse Endpoint", input: "Endpoint", output: EndpointType},
		{name: "Parse GoogleBQ", input: "GoogleBQ", output: GoogleBQType},
		{name: "Parse Kinesis", input: "Kinesis", output: KinesisType},
		{name: "Parse Lambda", input: "Lambda", output: LambdaType},
		{name: "Parse RestfulAPI", input: "RestfulAPI", output: RestfulAPIType},
		{name: "Parse S3", input: "S3", output: S3Type},
		{name: "Parse SQS", input: "SQS", output: SQSType},
		{name: "Parse SNS", input: "SNS", output: SNSType},
		{name: "Parse Unknown", input: "Unknown", output: UnknownType},
		{name: "Parse lowercase", input: "sqs", output: SQSType},
		{name: "Parse uppercase", input: "SNS", output: SNSType},
		{name: "Parse mixed case", input: "ApIgAtEwAy", output: APIGatewayType},
	}

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			result := ParseResourceType(tc.input)

			require.Equal(t, tc.output, result)
		})
	}
}
