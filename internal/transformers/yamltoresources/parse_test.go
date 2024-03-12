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

	tests := []struct {
		name      string
		args      args
		want      *resources.ResourceCollection
		targetErr error
	}{
		{
			name: "happy path",
			args: args{filename: "testdata/diagram.yaml"},
			want: wantResourceCollection,
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
