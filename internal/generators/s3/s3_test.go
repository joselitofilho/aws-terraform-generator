package s3

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

func TestS3_Build(t *testing.T) {
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
			name: "one bucket with extra files",
			fields: fields{
				configFileName: path.Join(testdataFolder, "s3.config.yaml"),
				output:         path.Join(testOutput, "one"),
			},
			extraValidations: func(tb testing.TB, err error) {
				if err != nil {
					return
				}

				modPath := path.Join(testOutput, "one", "mod")
				require.FileExists(tb, path.Join(modPath, "my-bucket-s3.tf"))
				require.FileExists(tb, path.Join(modPath, "s3.tf"))
			},
		},
		{
			name: "multiple buckets",
			fields: fields{
				configFileName: path.Join(testdataFolder, "s3.config.multiple.yaml"),
				output:         path.Join(testOutput, "multiple"),
			},
			extraValidations: func(tb testing.TB, err error) {
				if err != nil {
					return
				}

				require.FileExists(tb, path.Join(testOutput, "multiple", "mod", "s3.tf"))
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
			err := NewS3(tc.fields.configFileName, tc.fields.output).Build()

			require.ErrorIs(t, err, tc.targetErr)

			if tc.extraValidations != nil {
				tc.extraValidations(t, err)
			}
		})
	}
}
