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
		want ResourceARN
	}{
		{
			name: "kinesis arn",
			args: args{
				arn: "arn:aws:kinesis:${var.region}:${var.account_id}:stream/ProcessedLocationEvents",
			},
			want: ResourceARN{
				Key:  "kinesis",
				Name: "ProcessedLocationEvents",
			},
		},
		{
			name: "lambda as resource",
			args: args{
				arn: "aws_lambda_function.location_store_data_receiver_lambda.arn",
			},
			want: ResourceARN{
				Key:   "lambda",
				Name:  "locationStoreDataReceiver",
				Label: "location_store_data_receiver_lambda",
			},
		},
		{
			name: "lambda as module",
			args: args{
				arn: "module.location_store_data_receiver_lambda.lambda_invoke_arn",
			},
			want: ResourceARN{
				Key:   "lambda",
				Name:  "locationStoreDataReceiver",
				Label: "location_store_data_receiver_lambda",
			},
		},
		{
			name: "http arn",
			args: args{
				arn: "https://sqs.eu-west-1.amazonaws.com/var.account_id/pre-pipeline-events",
			},
			want: ResourceARN{
				Key:  "sqs",
				Name: "pre-pipeline-events",
			},
		},
	}

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			got := ResourceByARN(tc.args.arn)

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

	lambdaResource := resources.NewGenericResource("id", "myLambda", resources.LambdaType)

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
				res: lambdaResource,
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
				res: lambdaResource,
				filters: config.Filters{
					"lambda": config.Filter{NotMatch: []string{"^my"}},
				},
			},
			want: false,
		},
		{
			name: "nil resource",
			fields: fields{
				yamlConfig: &config.Config{},
				tfConfig:   &terraform.Config{},
			},
			args: args{
				res: nil,
				filters: config.Filters{
					"lambda": config.Filter{NotMatch: []string{"^my"}},
				},
			},
			want: false,
		},
		{
			name: "no filter",
			fields: fields{
				yamlConfig: &config.Config{},
				tfConfig:   &terraform.Config{},
			},
			args: args{
				res:     lambdaResource,
				filters: nil,
			},
			want: true,
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
			got := strTransformFromKeyValue(tc.args.key, tc.args.value, tc.args.suffix, tc.args.f)

			require.Equal(t, tc.want, got)
		})
	}
}
