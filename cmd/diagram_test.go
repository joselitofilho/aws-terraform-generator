package cmd

import (
	_ "embed"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDiagram_Run(t *testing.T) {
	type args struct {
		diagram    string
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
				diagram:    path.Join(testdataFolder, "diagram.xml"),
				configFile: path.Join(testdataFolder, "diagram.config.yaml"),
				output:     path.Join(testOutput, "diagram.yaml"),
			},
			extraValidations: func(tb testing.TB) {
				require.FileExists(tb, path.Join(testOutput, "diagram.yaml"))
			},
		},
		{
			name: "diagram file does not exist",
			args: args{
				diagram:    "fileDoesNotExist.xml",
				configFile: path.Join(testdataFolder, "diagram.config.yaml"),
				output:     path.Join(testOutput, "diagram.yaml"),
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
		{
			name: "diagram config file does not exist",
			args: args{
				diagram:    path.Join(testdataFolder, "diagram.xml"),
				configFile: "fileDoesNotExist.yaml",
				output:     path.Join(testOutput, "diagram.yaml"),
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

			_ = diagramCmd.Flags().Set(flagDiagram, tc.args.diagram)
			_ = diagramCmd.Flags().Set(flagConfig, tc.args.configFile)
			_ = diagramCmd.Flags().Set(flagOutput, tc.args.output)

			diagramCmd.Run(diagramCmd, []string{})

			if tc.extraValidations != nil {
				tc.extraValidations(t)
			}
		})
	}
}
