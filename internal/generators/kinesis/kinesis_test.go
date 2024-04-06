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
		extraValidations func(testing.TB, string, error)
		targetErr        error
	}{
		{
			name: "default templates for multiple kinesis",
			fields: fields{
				configFileName: path.Join(testdataFolder, "kinesis.config.multiple.yaml"),
				output:         path.Join(testOutput, "multiple"),
			},
			extraValidations: func(tb testing.TB, output string, err error) {
				if err != nil {
					return
				}

				require.FileExists(tb, path.Join(output, "mod", "kinesis.tf"))
			},
		},
		{
			name: "override default template for multiple kinesis",
			fields: fields{
				configFileName: path.Join(testdataFolder, "kinesis.config.override.default.tmpls.yaml"),
				output:         path.Join(testOutput, "override"),
			},
			extraValidations: func(tb testing.TB, output string, err error) {
				if err != nil {
					return
				}

				require.FileExists(tb, path.Join(output, "mod", "kinesis.tf"))
			},
		},
		{
			name: "at least one kinesis customising",
			fields: fields{
				configFileName: path.Join(testdataFolder, "kinesis.config.custom.yaml"),
				output:         path.Join(testOutput, "one"),
			},
			extraValidations: func(tb testing.TB, output string, err error) {
				if err != nil {
					return
				}

				modPath := path.Join(output, "mod")
				require.FileExists(tb, path.Join(modPath, "kinesis.tf"))
				require.FileExists(tb, path.Join(modPath, "myKinesis.tf"))
			},
		},
		{
			name: "all custom kinesis",
			fields: fields{
				configFileName: path.Join(testdataFolder, "kinesis.config.allcustom.yaml"),
				output:         path.Join(testOutput, "all"),
			},
			extraValidations: func(tb testing.TB, output string, err error) {
				if err != nil {
					return
				}

				modPath := path.Join(output, "mod")
				require.NoFileExists(tb, path.Join(modPath, "kinesis.tf"))
				require.FileExists(tb, path.Join(modPath, "myKinesis.tf"))
				require.FileExists(tb, path.Join(modPath, "myAnotherKinesis.tf"))
			},
		},
		{
			name: "when yaml parser fails should return an error",
			fields: fields{
				configFileName: "",
				output:         "",
			},
			targetErr: generatorserrs.ErrYAMLParser,
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
				tc.extraValidations(t, tc.fields.output, err)
			}
		})
	}
}
