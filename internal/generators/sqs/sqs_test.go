package sqs

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

func TestSQS_Build(t *testing.T) {
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
			name: "default templates for multiple sqs",
			fields: fields{
				configFileName: path.Join(testdataFolder, "sqs.config.multiple.yaml"),
				output:         path.Join(testOutput, "multiple"),
			},
			extraValidations: func(tb testing.TB, output string, err error) {
				if err != nil {
					return
				}

				require.FileExists(tb, path.Join(output, "mod", "sqs.tf"))
			},
		},
		{
			name: "at least one sqs customising",
			fields: fields{
				configFileName: path.Join(testdataFolder, "sqs.config.custom.yaml"),
				output:         path.Join(testOutput, "one"),
			},
			extraValidations: func(tb testing.TB, output string, err error) {
				if err != nil {
					return
				}

				modPath := path.Join(output, "mod")
				require.FileExists(tb, path.Join(modPath, "sqs.tf"))
				require.FileExists(tb, path.Join(modPath, "target-sqs.tf"))
			},
		},
		{
			name: "all custom sqs",
			fields: fields{
				configFileName: path.Join(testdataFolder, "sqs.config.allcustom.yaml"),
				output:         path.Join(testOutput, "all"),
			},
			extraValidations: func(tb testing.TB, output string, err error) {
				if err != nil {
					return
				}

				modPath := path.Join(output, "mod")
				require.NoFileExists(tb, path.Join(modPath, "sqs.tf"))
				require.FileExists(tb, path.Join(modPath, "target-sqs.tf"))
				require.FileExists(tb, path.Join(modPath, "source-sqs.tf"))
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
			err := NewSQS(tc.fields.configFileName, tc.fields.output).Build()

			require.ErrorIs(t, err, tc.targetErr)

			if tc.extraValidations != nil {
				tc.extraValidations(t, tc.fields.output, err)
			}
		})
	}
}
