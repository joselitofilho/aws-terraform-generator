package kinesis

import (
	_ "embed"
	"os"
	"path"
	"testing"

	generatorserrs "github.com/joselitofilho/aws-terraform-generator/internal/generators/errors"

	"github.com/stretchr/testify/require"
)

var (
	testdataFolder = "../testdata"
	testOutput     = "./testoutput"
)

func TestKinesis_Build(t *testing.T) {
	type fields struct {
		configFileName string
		output         string
	}

	tests := []struct {
		name             string
		fields           fields
		extraValidations func(testing.TB, error)
		targetErr        error
	}{
		{
			name: "one kinesis with extra file",
			fields: fields{
				configFileName: path.Join(testdataFolder, "kinesis.config.yaml"),
				output:         path.Join(testOutput, "one"),
			},
			extraValidations: func(tb testing.TB, err error) {
				if err != nil {
					return
				}

				modPath := path.Join(testOutput, "one", "mod")
				require.FileExists(tb, path.Join(modPath, "kinesis.tf"))
				require.FileExists(tb, path.Join(modPath, "custom.tf"))
			},
		},
		{
			name: "multiple kinesis",
			fields: fields{
				configFileName: path.Join(testdataFolder, "kinesis.config.multiple.yaml"),
				output:         path.Join(testOutput, "multiple"),
			},
			extraValidations: func(tb testing.TB, err error) {
				if err != nil {
					return
				}

				require.FileExists(tb, path.Join(testOutput, "multiple", "mod", "kinesis.tf"))
			},
		},
		{
			name: "when yaml parser fails should return an error",
			fields: fields{
				configFileName: "",
				output:         "",
			},
			targetErr: generatorserrs.ErrYAMLParse,
		},
	}

	defer func() {
		_ = os.RemoveAll(testOutput)
	}()

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			err := NewKinesis(tc.fields.configFileName, tc.fields.output).Build()

			require.ErrorIs(t, err, tc.targetErr)

			if tc.extraValidations != nil {
				tc.extraValidations(t, err)
			}
		})
	}
}