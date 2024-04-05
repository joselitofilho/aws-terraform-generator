package cmd

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDiagram_Run(t *testing.T) {
	var (
		testdataFolder = "./testdata"
		testOutput     = "./testoutput"
	)

	type args struct {
		diagram    string
		configFile string
		output     string
	}

	tests := []struct {
		name  string
		args  args
		setup func() (tearDown func())
	}{
		{
			name: "happy path",
			args: args{
				diagram:    path.Join(testdataFolder, "diagram.xml"),
				configFile: path.Join(testdataFolder, "diagram.config.yaml"),
				output:     path.Join(testOutput, "diagram.yaml"),
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tearDown := tt.setup()
				defer tearDown()
			}

			_ = diagramCmd.Flags().Set(flagDiagram, tt.args.diagram)
			_ = diagramCmd.Flags().Set(flagConfig, tt.args.configFile)
			_ = diagramCmd.Flags().Set(flagOutput, tt.args.output)

			diagramCmd.Run(diagramCmd, []string{})
		})
	}
}
