package apigateway

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAPIGateway_Build(t *testing.T) {
	type fields struct {
		configFileName string
		output         string
	}

	tests := []struct {
		name      string
		fields    fields
		targetErr error
	}{
		// TODO: Add test cases.
	}

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			err := NewAPIGateway(tc.fields.configFileName, tc.fields.output).Build()

			require.ErrorIs(t, err, tc.targetErr)
		})
	}
}
