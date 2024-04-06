package lambda

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

func TestLambda_Build(t *testing.T) {
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
			name: "happy path",
			fields: fields{
				configFileName: path.Join(testdataFolder, "lambda.config.yaml"),
				output:         path.Join(testOutput, "happypath", "teststack"),
			},
			extraValidations: func(tb testing.TB, output string, err error) {
				if err != nil {
					return
				}

				modPath := path.Join(output, "mod")
				require.FileExists(tb, path.Join(modPath, "exampleReceiver.tf"))

				lambdaPath := path.Join(output, "lambda", "exampleReceiver")
				require.FileExists(tb, path.Join(lambdaPath, "lambda.go"))
				require.FileExists(tb, path.Join(lambdaPath, "main.go"))
			},
		},
		{
			name: "override default template for multiple lambda",
			fields: fields{
				configFileName: path.Join(testdataFolder, "lambda.config.override.default.tmpls.yaml"),
				output:         path.Join(testOutput, "override", "teststack"),
			},
			extraValidations: func(tb testing.TB, output string, err error) {
				if err != nil {
					return
				}

				lambdaTf := path.Join(output, "mod", "exampleReceiver.tf")
				require.FileExists(tb, lambdaTf)

				lambdaTfData, err := os.ReadFile(lambdaTf)
				require.NoError(t, err)
				require.Equal(t, string(lambdaTfData), `resource "aws_lambda_function" "example_receiver_lambda" {}`)

				lambdaPath := path.Join(output, "lambda", "exampleReceiver")

				lambdaGo := path.Join(lambdaPath, "lambda.go")
				require.FileExists(tb, lambdaGo)

				lambdaGoData, err := os.ReadFile(lambdaGo)
				require.NoError(t, err)
				require.Equal(t, string(lambdaGoData), "package main\n")

				mainGo := path.Join(lambdaPath, "main.go")
				require.FileExists(tb, mainGo)

				mainGoData, err := os.ReadFile(mainGo)
				require.NoError(t, err)
				require.Equal(t, string(mainGoData), "package main\n")
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
			err := NewLambda(tc.fields.configFileName, tc.fields.output).Build()

			require.ErrorIs(t, err, tc.targetErr)

			if tc.extraValidations != nil {
				tc.extraValidations(t, tc.fields.output, err)
			}
		})
	}
}
