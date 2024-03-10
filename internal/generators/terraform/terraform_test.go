package terraform

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseTerraformFiles(t *testing.T) {
	type args struct {
		directories []string
		files       []string
	}

	tests := []struct {
		name      string
		args      args
		want      Config
		targetErr error
	}{
		{
			name: "apig.tf",
			args: args{
				directories: []string{"testdata/apigateway"},
			},
			want: Config{},
		},
		{
			name: "lambda as module",
			args: args{
				files: []string{"testdata/lambda/lambdaAsModule.tf"},
			},
			want: Config{},
		},
		{
			name: "lambda as resource",
			args: args{
				files: []string{"testdata/lambda/lambdaAsResource.tf"},
			},
			want: Config{},
		},
		{
			name: "sqs.tf",
			args: args{
				directories: []string{"testdata/sqs"},
			},
			want: Config{},
		},
	}

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			_, err := Parse(tc.args.directories, tc.args.files)

			require.ErrorIs(t, err, tc.targetErr)
		})
	}
}
