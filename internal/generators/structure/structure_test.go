package structure

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

func TestStructure_Build(t *testing.T) {
	type fields struct {
		configFileName string
		output         string
	}

	happypathPath := path.Join(testOutput, "happypath")

	tests := []struct {
		name             string
		fields           fields
		extraValidations func(testing.TB, error)
		targetErr        error
	}{
		{
			name: "happy path",
			fields: fields{
				configFileName: path.Join(testdataFolder, "structure.config.yaml"),
				output:         happypathPath,
			},
			extraValidations: func(tb testing.TB, err error) {
				if err != nil {
					return
				}

				teststackPath := path.Join(happypathPath, "teststack")

				devPath := path.Join(teststackPath, "dev")
				require.FileExists(tb, path.Join(devPath, "main.tf"))
				require.FileExists(tb, path.Join(devPath, "terragrunt.hcl"))
				require.FileExists(tb, path.Join(devPath, "vars.tf"))

				uatPath := path.Join(teststackPath, "uat")
				require.FileExists(tb, path.Join(uatPath, "main.tf"))
				require.FileExists(tb, path.Join(uatPath, "terragrunt.hcl"))
				require.FileExists(tb, path.Join(uatPath, "vars.tf"))

				prdPath := path.Join(teststackPath, "prd")
				require.FileExists(tb, path.Join(prdPath, "main.tf"))
				require.FileExists(tb, path.Join(prdPath, "terragrunt.hcl"))
				require.FileExists(tb, path.Join(prdPath, "vars.tf"))

				modPath := path.Join(teststackPath, "mod")
				require.FileExists(tb, path.Join(modPath, "main.tf"))
				require.FileExists(tb, path.Join(modPath, "vars.tf"))

				require.DirExists(tb, teststackPath, "lambda")
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
			err := NewStructure(tc.fields.configFileName, tc.fields.output).Build()

			require.ErrorIs(t, err, tc.targetErr)

			if tc.extraValidations != nil {
				tc.extraValidations(t, err)
			}
		})
	}
}
