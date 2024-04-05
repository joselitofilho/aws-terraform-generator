package diagram

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

func TestDiagram_Build(t *testing.T) {
	type fields struct {
		diagramFilename string
		configFilename  string
		output          string
	}

	tests := []struct {
		name      string
		fields    fields
		targetErr error
	}{
		{
			name: "happy path",
			fields: fields{
				diagramFilename: path.Join(testdataDir, "diagram.xml"),
				configFilename:  path.Join(testdataDir, "diagram.config.yaml"),
				output:          testOutput,
			},
		},
	}

	defer func() {
		_ = os.RemoveAll(testOutput)
	}()

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			err := NewDiagram(tc.fields.diagramFilename, tc.fields.configFilename, tc.fields.output).Build()

			require.ErrorIs(t, err, tc.targetErr)
		})
	}
}
