package utils

import (
	"errors"
	"go/format"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	tfFileName = "test.tf"
	tfContent  = []byte("tf content")

	errDummy = errors.New("dummy error")
)

func TestTerraformFormat(t *testing.T) {
	type args struct {
		folder string
	}

	tests := []struct {
		name           string
		setup          func(testing.TB) (tearDown func())
		args           args
		targetErr      error
		expectedErrMsg string
	}{
		{
			name: "happy path",
			setup: func(tb testing.TB) (tearDown func()) {
				currTerraformCommand := terraformCommand
				terraformCommand = func(folder string) error {
					require.Equal(tb, "./output", folder)

					return nil
				}

				return func() {
					terraformCommand = currTerraformCommand
				}
			},
			args: args{
				folder: "./output",
			},
		},
		{
			name: "when terraformCommand fails should return an error",
			setup: func(tb testing.TB) (tearDown func()) {
				currTerraformCommand := terraformCommand
				terraformCommand = func(folder string) error {
					require.Equal(tb, "./output", folder)

					return errDummy
				}

				return func() {
					terraformCommand = currTerraformCommand
				}
			},
			args: args{
				folder: "./output",
			},
			targetErr:      errDummy,
			expectedErrMsg: "please consider to install terraform. Terraform format fails: dummy error",
		},
	}

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			tearDown := tc.setup(t)
			defer tearDown()

			err := TerraformFormat(tc.args.folder)

			require.ErrorIs(t, err, tc.targetErr)

			if tc.expectedErrMsg != "" {
				require.Equal(t, tc.expectedErrMsg, err.Error())
			}
		})
	}
}

func TestGoFormat(t *testing.T) {
	type args struct {
		filename string
	}

	tests := []struct {
		name      string
		setup     func(tb testing.TB) (tearDown func())
		args      args
		targetErr error
	}{
		{
			name: "happy path",
			setup: func(tb testing.TB) func() {
				osReadFile = func(name string) ([]byte, error) {
					require.Equal(t, tfFileName, name)
					return tfContent, nil
				}

				formatSource = func(src []byte) ([]byte, error) {
					require.Equal(tb, tfContent, src)
					return tfContent, nil
				}

				osWriteFile = func(name string, data []byte, perm os.FileMode) error {
					require.Equal(tb, tfFileName, name)
					require.Equal(tb, tfContent, data)
					require.Equal(tb, os.ModePerm, perm)
					return nil
				}

				return func() {
					osReadFile = os.ReadFile
					formatSource = format.Source
					osWriteFile = os.WriteFile
				}
			},
			args: args{
				filename: tfFileName,
			},
		},
		{
			name: "when os.ReadFile fails should return an error",
			setup: func(tb testing.TB) func() {
				osReadFile = func(name string) ([]byte, error) {
					require.Equal(tb, tfFileName, name)
					return nil, errDummy
				}

				return func() {
					osReadFile = os.ReadFile
				}
			},
			args: args{
				filename: tfFileName,
			},
			targetErr: errDummy,
		},
		{
			name: "when format.Source fails should return an error",
			setup: func(tb testing.TB) func() {
				osReadFile = func(name string) ([]byte, error) {
					require.Equal(tb, tfFileName, name)
					return tfContent, nil
				}

				formatSource = func(src []byte) ([]byte, error) {
					require.Equal(tb, tfContent, src)
					return nil, errDummy
				}

				return func() {
					osReadFile = os.ReadFile
					formatSource = format.Source
				}
			},
			args: args{
				filename: tfFileName,
			},
			targetErr: errDummy,
		},
		{
			name: "when os.WriteFile fails should return an error",
			setup: func(tb testing.TB) func() {
				osReadFile = func(name string) ([]byte, error) {
					require.Equal(t, tfFileName, name)
					return tfContent, nil
				}

				formatSource = func(src []byte) ([]byte, error) {
					require.Equal(tb, tfContent, src)
					return tfContent, nil
				}

				osWriteFile = func(name string, data []byte, perm os.FileMode) error {
					require.Equal(tb, tfFileName, name)
					require.Equal(tb, tfContent, data)
					require.Equal(tb, os.ModePerm, perm)
					return errDummy
				}

				return func() {
					osReadFile = os.ReadFile
					formatSource = format.Source
					osWriteFile = os.WriteFile
				}
			},
			args: args{
				filename: tfFileName,
			},
			targetErr: errDummy,
		},
	}

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			tearDown := tc.setup(t)
			defer tearDown()

			err := GoFormat(tc.args.filename)

			require.ErrorIs(t, err, tc.targetErr)
		})
	}
}
