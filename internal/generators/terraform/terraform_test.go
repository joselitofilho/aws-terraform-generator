package terraform

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseTerraformFiles(t *testing.T) {
	type args struct {
		directory string
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
				directory: "testdata/apigateway",
			},
			want: Config{},
		},
		{
			name: "loctionEventProcessor lambda",
			args: args{
				directory: "testdata/lambda",
			},
			want: Config{},
		},
		{
			name: "sqs.tf",
			args: args{
				directory: "testdata/sqs",
			},
			want: Config{},
		},
	}

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			_, err := Parse(tc.args.directory)

			require.ErrorIs(t, err, tc.targetErr)
			// require.Equal(t, tc.want, got)
		})
	}
}
