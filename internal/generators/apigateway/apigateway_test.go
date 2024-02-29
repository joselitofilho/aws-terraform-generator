package apigateway

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

func TestAPIGateway_Build(t *testing.T) {
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
				configFileName: path.Join(testdataFolder, "apigateway.config.yaml"),
				output:         happypathPath,
			},
			extraValidations: func(tb testing.TB, err error) {
				if err != nil {
					return
				}

				modPath := path.Join(happypathPath, "teststack", "mod")
				require.FileExists(tb, path.Join(modPath, "apig.tf"))
				require.FileExists(tb, path.Join(modPath, "exampleAPIReceiver.tf"))

				lambdaPath := path.Join(happypathPath, "teststack", "lambda", "exampleAPIReceiver")
				require.FileExists(tb, path.Join(lambdaPath, "lambda.go"))
				require.FileExists(tb, path.Join(lambdaPath, "main.go"))
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
			err := NewAPIGateway(tc.fields.configFileName, tc.fields.output).Build()

			require.ErrorIs(t, err, tc.targetErr)

			if tc.extraValidations != nil {
				tc.extraValidations(t, err)
			}
		})
	}
}
