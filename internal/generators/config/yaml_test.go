package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestYAML_Parse(t *testing.T) {
	type fields struct {
		fileName string
	}

	tests := []struct {
		name      string
		fields    fields
		want      *Config
		targetErr error
	}{
		// TODO: Add test cases.
	}

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			got, err := NewYAML(tc.fields.fileName).Parse()

			if tc.targetErr == nil {
				require.NoError(t, err)
				require.Equal(t, tc.want, got)
			} else {
				require.ErrorIs(t, err, tc.targetErr)
			}
		})
	}
}
