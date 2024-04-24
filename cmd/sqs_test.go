package cmd

import (
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSqs_Run(t *testing.T) {
	type args struct {
		configFile string
		output     string
	}

	tests := []struct {
		name             string
		args             args
		setup            func() (tearDown func())
		extraValidations func(testing.TB)
	}{
		{
			name: "happy path",
			args: args{
				configFile: path.Join(testdataFolder, "sqs.config.yaml"),
				output:     path.Join(testOutput),
			},
			extraValidations: func(tb testing.TB) {
				require.FileExists(tb, path.Join(testOutput, "mod/my-test-sqs.tf"))
			},
		},
	}

	// defer func() {
	// 	_ = os.RemoveAll(testOutput)
	// }()

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			if tc.setup != nil {
				tearDown := tc.setup()
				defer tearDown()
			}

			_ = sqsCmd.Flags().Set(flagConfig, tc.args.configFile)
			_ = sqsCmd.Flags().Set(flagOutput, tc.args.output)

			sqsCmd.Run(sqsCmd, []string{})

			if tc.extraValidations != nil {
				tc.extraValidations(t)
			}
		})
	}
}
