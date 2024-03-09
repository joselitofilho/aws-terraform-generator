package terraformtoresources

import (
	"testing"

	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
	"github.com/joselitofilho/aws-terraform-generator/internal/generators/terraform"
	"github.com/joselitofilho/aws-terraform-generator/internal/resources"
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

func TestTransformer_hasResourceMatched(t *testing.T) {
	type fields struct {
		yamlConfig *config.Config
		tfConfig   *terraform.Config
	}

	type args struct {
		res     resources.Resource
		filters config.Filters
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "match",
			fields: fields{
				yamlConfig: &config.Config{},
				tfConfig:   &terraform.Config{},
			},
			args: args{
				res: resources.NewGenericResource("id", "myLambda", resources.LambdaType),
				filters: config.Filters{
					"lambda": config.Filter{Match: []string{"^my"}},
				},
			},
			want: true,
		},
		{
			name: "not match",
			fields: fields{
				yamlConfig: &config.Config{},
				tfConfig:   &terraform.Config{},
			},
			args: args{
				res: resources.NewGenericResource("id", "myLambda", resources.LambdaType),
				filters: config.Filters{
					"lambda": config.Filter{NotMatch: []string{"^my"}},
				},
			},
			want: false,
		},
	}

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			tr := NewTransformer(tc.fields.yamlConfig, tc.fields.tfConfig)

			got := tr.hasResourceMatched(tc.args.res, tc.args.filters)

			require.Equal(t, tc.want, got)
		})
	}
}

func Test_strTransformFromEnvar(t *testing.T) {
	type args struct {
		key    string
		value  string
		suffix string
		f      func(s string) string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "key is equal to suffix and value contains the client and environment vars",
			args: args{
				key:    "MY_ENVAR",
				value:  "var.client-var.environment-value",
				suffix: "MY_ENVAR",
				f: func(s string) string {
					require.Equal(t, s, "value")
					return s
				},
			},
			want: "value",
		},
		{
			name: "key is equal to suffix and value contains the kinesis arn",
			args: args{
				key:    "MY_ENVAR",
				value:  "aws_kinesis_stream.myValue.arn",
				suffix: "MY_ENVAR",
				f: func(s string) string {
					require.Equal(t, s, "MyValue")
					return s
				},
			},
			want: "MyValue",
		},
		{
			name: "key is equal to suffix and value contains the lambda arn",
			args: args{
				key:    "MY_ENVAR",
				value:  "aws_lambda_function.myValue.arn",
				suffix: "MY_ENVAR",
				f: func(s string) string {
					require.Equal(t, s, "myValue")
					return s
				},
			},
			want: "myValue",
		},
		{
			name: "key is equal to suffix and value contains the s3 bucket arn",
			args: args{
				key:    "MY_ENVAR",
				value:  "aws_s3_bucket.myValue.arn",
				suffix: "MY_ENVAR",
				f: func(s string) string {
					require.Equal(t, s, "myValue")
					return s
				},
			},
			want: "myValue",
		},
		{
			name: "key is equal to suffix and value contains the sqs arn",
			args: args{
				key:    "MY_ENVAR",
				value:  "aws_sqs_queue.myValue.arn",
				suffix: "MY_ENVAR",
				f: func(s string) string {
					require.Equal(t, s, "my-value")
					return s
				},
			},
			want: "my-value",
		},
	}

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			got := strTransformFromEnvar(tc.args.key, tc.args.value, tc.args.suffix, tc.args.f)

			require.Equal(t, tc.want, got)
		})
	}
}
