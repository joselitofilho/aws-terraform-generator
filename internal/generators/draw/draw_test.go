package draw

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	testdataDir = "testdata"
	testOutput  = "testoutput"
)

func TestDraw_Build(t *testing.T) {
	type fields struct {
		workdirs       []string
		files          []string
		configFileName string
		output         string
	}

	tests := []struct {
		name      string
		fields    fields
		targetErr error
	}{
		{
			name: "happy path",
			fields: fields{
				workdirs:       []string{path.Join(testdataDir, "mystack")},
				configFileName: path.Join(testdataDir, "draw.config.yaml"),
				output:         testOutput,
			},
		},
	}

	defer func() {
		_ = os.RemoveAll(testOutput)
	}()

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			d := NewDraw(
				tc.fields.workdirs,
				tc.fields.files,
				tc.fields.configFileName,
				tc.fields.output,
			)

			_ = os.MkdirAll(tc.fields.output, os.ModePerm)

			err := d.Build()

			require.ErrorIs(t, err, tc.targetErr)
			require.FileExists(t, path.Join(testOutput, "diagram.dot"))
		})
	}
}
