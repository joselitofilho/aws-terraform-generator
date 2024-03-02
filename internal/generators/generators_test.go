package generators

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuild(t *testing.T) {
	type args struct {
		data            any
		templateName    string
		templateContent string
	}

	tests := []struct {
		name      string
		args      args
		want      string
		targetErr error
	}{
		// TODO: Add test cases.
	}

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			got, err := Build(tc.args.data, tc.args.templateName, tc.args.templateContent)

			require.ErrorIs(t, err, tc.targetErr)
			require.Equal(t, tc.want, got)
		})
	}
}
