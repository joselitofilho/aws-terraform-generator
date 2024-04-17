package cmd

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSns_Run(t *testing.T) {
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
				configFile: path.Join(testdataFolder, "sns.config.yaml"),
				output:     path.Join(testOutput),
			},
			extraValidations: func(tb testing.TB) {
				require.FileExists(tb, path.Join(testOutput, "mod/my-test-sns.tf"))
			},
		},
		{
			name: "sns config file does not exist",
			args: args{
				configFile: "fileDoesNotExist.yaml",
				output:     path.Join(testOutput, "mod/my-test-sns.tf"),
			},
			setup: func() (tearDown func()) {
				osExit = func(code int) {
					require.Equal(t, 1, code)
				}

				return func() {
					osExit = os.Exit
				}
			},
		},
	}

	defer func() {
		_ = os.RemoveAll(testOutput)
	}()

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			if tc.setup != nil {
				tearDown := tc.setup()
				defer tearDown()
			}

			_ = snsCmd.Flags().Set(flagConfig, tc.args.configFile)
			_ = snsCmd.Flags().Set(flagOutput, tc.args.output)

			snsCmd.Run(snsCmd, []string{})

			if tc.extraValidations != nil {
				tc.extraValidations(t)
			}
		})
	}
}
