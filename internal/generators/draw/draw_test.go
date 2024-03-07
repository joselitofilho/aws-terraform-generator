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
		configFileName string
		input          string
		output         string
	}

	tests := []struct {
		name      string
		fields    fields
		targetErr error
	}{
		{
			name: "",
			fields: fields{
				configFileName: path.Join(testdataDir, "draw.config.yaml"),
				input:          "/Users/joselitofilho/dev/personal/aws-terraform-generator/output/location", // path.Join(testdataDir, "stack"),
				output:         testOutput,
			},
		},
	}

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			d := NewDraw(
				tc.fields.configFileName,
				tc.fields.input,
				tc.fields.output,
			)

			_ = os.MkdirAll(tc.fields.output, os.ModePerm)

			err := d.Build()

			require.ErrorIs(t, err, tc.targetErr)
		})
	}
}
