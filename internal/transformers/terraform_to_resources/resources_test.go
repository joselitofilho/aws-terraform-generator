package terraform_to_resources

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_resourceByARN(t *testing.T) {
	type args struct {
		arn string
	}

	tests := []struct {
		name string
		args args
		want resourceARN
	}{
		{
			name: "kinesis arn",
			args: args{
				arn: "arn:aws:kinesis:${var.region}:${var.account_id}:stream/ProcessedLocationEvents",
			},
			want: resourceARN{
				key:  "kinesis",
				name: "ProcessedLocationEvents",
			},
		},
		{
			name: "lambda as resource",
			args: args{
				arn: "aws_lambda_function.location_store_data_receiver_lambda.arn",
			},
			want: resourceARN{
				key:  "lambda",
				name: "locationStoreDataReceiver",
			},
		},
		{
			name: "lambda as module",
			args: args{
				arn: "module.location_store_data_receiver_lambda.lambda_invoke_arn",
			},
			want: resourceARN{
				key:  "lambda",
				name: "locationStoreDataReceiver",
			},
		},
	}

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			got := resourceByARN(tc.args.arn)

			require.Equal(t, tc.want, got)
		})
	}
}
