package diagram

import (
	_ "embed"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	testdataDir = "testdata"
	testOutput  = "testoutput"
)

//go:embed testdata/diagram.yaml
var diagramYAML []byte

func TestDiagram_Build(t *testing.T) {
	type fields struct {
		diagramFilename string
		configFilename  string
		output          string
	}

	tests := []struct {
		name             string
		fields           fields
		extraValidations func(testing.TB)
		targetErr        error
	}{
		{
			name: "happy path",
			fields: fields{
				diagramFilename: path.Join(testdataDir, "diagram.xml"),
				configFilename:  path.Join(testdataDir, "diagram.config.yaml"),
				output:          path.Join(testOutput, "diagram.yaml"),
			},
			extraValidations: func(tb testing.TB) {
				data, err := os.ReadFile(path.Join(testOutput, "diagram.yaml"))
				require.NoError(tb, err)
				require.Equal(tb, string(diagramYAML), string(data))
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

			if tc.extraValidations != nil {
				tc.extraValidations(t)
			}
		})
	}
}
