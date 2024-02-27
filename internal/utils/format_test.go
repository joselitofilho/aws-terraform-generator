package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTerraformFormat(t *testing.T) {
	type args struct {
		folder string
	}

	tests := []struct {
		name      string
		args      args
		targetErr error
	}{
		// TODO: Add test cases.
	}

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			err := TerraformFormat(tc.args.folder)

			require.ErrorIs(t, err, tc.targetErr)
		})
	}
}

func TestGoFormat(t *testing.T) {
	type args struct {
		fileName string
	}

	tests := []struct {
		name      string
		args      args
		targetErr error
	}{
		// TODO: Add test cases.
	}

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			err := GoFormat(tc.args.fileName)

			require.ErrorIs(t, err, tc.targetErr)
		})
	}
}
