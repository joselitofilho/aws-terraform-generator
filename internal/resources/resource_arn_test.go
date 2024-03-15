package resources

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResourceARN_resourceByARN(t *testing.T) {
	type args struct {
		arn              string
		suggestedResType ResourceType
	}

	tests := []struct {
		name string
		args args
		want ResourceARN
	}{
		{
			name: "kinesis arn",
			args: args{
				arn:              "arn:aws:kinesis:${var.region}:${var.account_id}:stream/ProcessedLocationEvents",
				suggestedResType: KinesisType,
			},
			want: ResourceARN{
				Type:  "aws_kinesis_stream",
				Name:  "ProcessedLocationEvents",
				Label: "",
			},
		},
		{
			name: "lambda as resource",
			args: args{
				arn:              "aws_lambda_function.location_store_data_receiver_lambda.arn",
				suggestedResType: LambdaType,
			},
			want: ResourceARN{
				Type:  "aws_lambda_function",
				Name:  "",
				Label: "location_store_data_receiver_lambda",
			},
		},
		{
			name: "lambda as module",
			args: args{
				arn:              "module.location_store_data_receiver_lambda.lambda_invoke_arn",
				suggestedResType: LambdaType,
			},
			want: ResourceARN{
				Type:  "aws_lambda_function",
				Name:  "",
				Label: "location_store_data_receiver_lambda",
			},
		},
		{
			name: "http or https",
			args: args{
				arn:              "https://sqs.eu-west-1.amazonaws.com/var.account_id/pre-pipeline-events",
				suggestedResType: SQSType,
			},
			want: ResourceARN{
				Type:  "aws_sqs_queue",
				Name:  "pre-pipeline-events",
				Label: "",
			},
		},
		{
			name: "variable",
			args: args{
				arn:              "local.arn",
				suggestedResType: UnknownType,
			},
			want: ResourceARN{
				Type:  "",
				Name:  "local.arn",
				Label: "",
			},
		},
		{
			name: "any string",
			args: args{
				arn:              "Name",
				suggestedResType: UnknownType,
			},
			want: ResourceARN{
				Type:  "",
				Name:  "Name",
				Label: "",
			},
		},
		{
			name: "any string with a suggestion for the resource type",
			args: args{
				arn:              "Name",
				suggestedResType: SQSType,
			},
			want: ResourceARN{
				Type:  "aws_sqs_queue",
				Name:  "Name",
				Label: "",
			},
		},
		{
			name: "infer kinesis",
			args: args{
				arn:              "aws_kinesis_stream.my_stream.arn",
				suggestedResType: UnknownType,
			},
			want: ResourceARN{
				Type:  "aws_kinesis_stream",
				Name:  "",
				Label: "my_stream",
			},
		},
		{
			name: "infer lambda",
			args: args{
				arn:              "aws_lambda_function.location_store_data_receiver_lambda.arn",
				suggestedResType: UnknownType,
			},
			want: ResourceARN{
				Type:  "aws_lambda_function",
				Name:  "",
				Label: "location_store_data_receiver_lambda",
			},
		},
		{
			name: "infer s3 bucket",
			args: args{
				arn:              "aws_s3_bucket.my_bucket.arn",
				suggestedResType: UnknownType,
			},
			want: ResourceARN{
				Type:  "aws_s3_bucket",
				Name:  "",
				Label: "my_bucket",
			},
		},
		{
			name: "infer sns",
			args: args{
				arn:              "aws_sns_topic.my_notification.arn",
				suggestedResType: UnknownType,
			},
			want: ResourceARN{
				Type:  "aws_sns_topic",
				Name:  "",
				Label: "my_notification",
			},
		},
		{
			name: "infer sqs",
			args: args{
				arn:              "aws_sqs_queue.my_queue.arn",
				suggestedResType: UnknownType,
			},
			want: ResourceARN{
				Type:  "aws_sqs_queue",
				Name:  "",
				Label: "my_queue",
			},
		},
	}

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			got := ParseResourceARN(tc.args.arn, tc.args.suggestedResType)

			require.Equal(t, tc.want, got)
		})
	}
}

func TestResourceARN_LabelOrName(t *testing.T) {
	type fields struct {
		Type  string
		Name  string
		Label string
	}

	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "when name and label are filled then use label",
			fields: fields{
				Type:  "aws_sqs_queue",
				Name:  "my-queue",
				Label: "my_queue_sqs",
			},
			want: "my_queue_sqs",
		},
		{
			name: "when only name is filled then use name",
			fields: fields{
				Type:  "aws_sqs_queue",
				Name:  "my-queue",
				Label: "",
			},
			want: "my-queue",
		},
		{
			name: "when only label is filled then use label",
			fields: fields{
				Type:  "aws_sqs_queue",
				Name:  "",
				Label: "my_queue_sqs",
			},
			want: "my_queue_sqs",
		},
	}

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			r := &ResourceARN{
				Type:  tc.fields.Type,
				Name:  tc.fields.Name,
				Label: tc.fields.Label,
			}

			got := r.LabelOrName()

			require.Equal(t, tc.want, got)
		})
	}
}
