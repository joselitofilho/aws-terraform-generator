package generators

import (
	_ "embed"
	"os"
	"path"
	"testing"

	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
	"github.com/stretchr/testify/require"
)

func TestBuild(t *testing.T) {
	type args struct {
		data            any
		templateName    string
		templateContent string
	}

	tests := []struct {
		name           string
		args           args
		want           string
		expectedErrMsg string
	}{
		{
			name: "valid case",
			args: args{
				data:            map[string]any{"Name": "John", "Age": 30},
				templateName:    "example",
				templateContent: "Name: {{.Name}}, Age: {{.Age}}",
			},
			want: "Name: John, Age: 30",
		},
		{
			name: "getFileByName func",
			args: args{
				data: map[string]any{"Files": map[string]File{"lambda.go": {
					Imports: []string{"context"},
					Tmpl:    "package main",
				}}},
				templateName:    "example",
				templateContent: `{{getFileByName $.Files "lambda.go"}}`,
			},
			want: "{package main [context]}",
		},
		{
			name: "getFileImports func",
			args: args{
				data: map[string]any{"Files": map[string]File{"lambda.go": {
					Imports: []string{"context"},
				}}},
				templateName:    "example",
				templateContent: `{{ range getFileImports $.Files "lambda.go" }}"{{ . }}"{{end}}`,
			},
			want: `"context"`,
		},
		{
			name: "ToCamel func",
			args: args{
				data:            map[string]any{"Name": "my-name"},
				templateName:    "example",
				templateContent: "Name: {{ToCamel .Name}}",
			},
			want: "Name: myName",
		},
		{
			name: "ToKebab func",
			args: args{
				data:            map[string]any{"Name": "myName"},
				templateName:    "example",
				templateContent: "Name: {{ToKebab .Name}}",
			},
			want: "Name: my-name",
		},
		{
			name: "ToLower func",
			args: args{
				data:            map[string]any{"Name": "MY-NAME"},
				templateName:    "example",
				templateContent: "Name: {{ToLower .Name}}",
			},
			want: "Name: my-name",
		},
		{
			name: "ToPascal func",
			args: args{
				data:            map[string]any{"Name": "my-name"},
				templateName:    "example",
				templateContent: "Name: {{ToPascal .Name}}",
			},
			want: "Name: MyName",
		},
		{
			name: "ToSpace func",
			args: args{
				data:            map[string]any{"Name": "my-name"},
				templateName:    "example",
				templateContent: "Name: {{ToSpace .Name}}",
			},
			want: "Name: my name",
		},
		{
			name: "ToSnake func",
			args: args{
				data:            map[string]any{"Name": "my-name"},
				templateName:    "example",
				templateContent: "Name: {{ToSnake .Name}}",
			},
			want: "Name: my_name",
		},
		{
			name: "ToUpper func",
			args: args{
				data:            map[string]any{"Name": "my-name"},
				templateName:    "example",
				templateContent: "Name: {{ToUpper .Name}}",
			},
			want: "Name: MY-NAME",
		},
		{
			name: "missing field",
			args: args{
				data:            map[string]any{"Name": "John", "Age": 30},
				templateName:    "invalid",
				templateContent: "{{ .MissingField }}",
			},
			want: "<no value>",
		},
		{
			name: "invalid template",
			args: args{
				data:            map[string]any{"Name": "John", "Age": 30},
				templateName:    "invalid",
				templateContent: "{{ InvalidFunction .Name }}",
			},
			want:           "",
			expectedErrMsg: "template: invalid:1: function \"InvalidFunction\" not defined",
		},
	}

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			got, err := Build(tc.args.data, tc.args.templateName, tc.args.templateContent)

			if tc.expectedErrMsg == "" {
				require.NoError(t, err)
				require.Equal(t, tc.want, got)
			} else {
				require.Error(t, err)
				require.Equal(t, tc.expectedErrMsg, err.Error())
			}
		})
	}
}

func TestCreateFilesMap(t *testing.T) {
	type args struct {
		files []config.File
	}

	tests := []struct {
		name string
		args args
		want map[string]File
	}{
		{
			name: "empty input",
			args: args{
				files: []config.File{},
			},
			want: map[string]File{},
		},
		{
			name: "single file",
			args: args{
				files: []config.File{
					{
						Name:    "example.txt",
						Tmpl:    "This is an example file.",
						Imports: []string{},
					},
				},
			},
			want: map[string]File{
				"example.txt": {
					Tmpl:    "This is an example file.",
					Imports: []string{},
				},
			},
		},
		{
			name: "multiple files",
			args: args{
				files: []config.File{
					{
						Name:    "file1.txt",
						Tmpl:    "Template for file 1",
						Imports: []string{"package1", "package2"},
					},
					{
						Name:    "file2.txt",
						Tmpl:    "Template for file 2",
						Imports: []string{"package3"},
					},
				},
			},
			want: map[string]File{
				"file1.txt": {
					Tmpl:    "Template for file 1",
					Imports: []string{"package1", "package2"},
				},
				"file2.txt": {
					Tmpl:    "Template for file 2",
					Imports: []string{"package3"},
				},
			},
		},
	}

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			got := CreateFilesMap(tc.args.files)

			require.Equal(t, tc.want, got)
		})
	}
}

func TestCreateTemplatesMap(t *testing.T) {
	type args struct {
		filenameTemplatesList []config.FilenameTemplateMap
	}

	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "empty List",
			args: args{
				filenameTemplatesList: []config.FilenameTemplateMap{},
			},
			want: map[string]string{},
		},
		{
			name: "single element",
			args: args{
				filenameTemplatesList: []config.FilenameTemplateMap{
					{"file1.txt": "template1"},
				},
			},
			want: map[string]string{"file1.txt": "template1"},
		},
		{
			name: "multiple elements",
			args: args{
				filenameTemplatesList: []config.FilenameTemplateMap{
					{"file1.txt": "template1"},
					{"file2.txt": "template2"},
					{"file3.txt": "template3"},
				},
			},
			want: map[string]string{
				"file1.txt": "template1",
				"file2.txt": "template2",
				"file3.txt": "template3",
			},
		},
		{
			name: "duplicate keys",
			args: args{
				filenameTemplatesList: []config.FilenameTemplateMap{
					{"file1.txt": "template1"},
					{"file1.txt": "template2"},
				},
			},
			want: map[string]string{"file1.txt": "template2"},
		},
	}

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			got := CreateTemplatesMap(tc.args.filenameTemplatesList)

			require.Equal(t, tc.want, got)
		})
	}
}

func TestFilterTemplatesMap(t *testing.T) {
	type args struct {
		filter       string
		templatesMap map[string]string
	}

	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "empty Map",
			args: args{
				filter:       "test",
				templatesMap: map[string]string{},
			},
			want: map[string]string{},
		},
		{
			name: "no matches",
			args: args{
				filter: "xyz",
				templatesMap: map[string]string{
					"file1.txt": "template1",
					"file2.txt": "template2",
					"file3.txt": "template3",
				},
			},
			want: map[string]string{},
		},
		{
			name: "single match",
			args: args{
				filter: "2",
				templatesMap: map[string]string{
					"file1.txt": "template1",
					"file2.txt": "template2",
					"file3.txt": "template3",
				},
			},
			want: map[string]string{"file2.txt": "template2"},
		},
		{
			name: "multiple matches",
			args: args{
				filter: "file",
				templatesMap: map[string]string{
					"file1.txt": "template1",
					"file2.txt": "template2",
					"file3.txt": "template3",
					"other.txt": "othertemplate",
				},
			},
			want: map[string]string{
				"file1.txt": "template1",
				"file2.txt": "template2",
				"file3.txt": "template3",
			},
		},
	}

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			got := FilterTemplatesMap(tc.args.filter, tc.args.templatesMap)

			require.Equal(t, tc.want, got)
		})
	}
}

func TestGenerateFile(t *testing.T) {
	type args struct {
		templatesMap map[string]string
		fileName     string
		fileTmpl     string
		outputFile   string
		data         any
	}

	testOutput := "./testoutput"
	_ = os.MkdirAll(testOutput, os.ModePerm)

	tests := []struct {
		name             string
		args             args
		extraValidations func(testing.TB, string, error)
		targetErr        error
	}{
		{
			name: "successful go file generation and formatting",
			args: args{
				templatesMap: map[string]string{"test.go": "type  My{{.Name}}Struct    struct   {}"},
				fileName:     "test.go",
				outputFile:   path.Join(testOutput, "output.go"),
				data:         struct{ Name string }{Name: "World"},
			},
			extraValidations: func(tb testing.TB, outputFile string, err error) {
				if err != nil {
					return
				}

				data, err := os.ReadFile(outputFile)
				require.NoError(tb, err)
				require.Equal(tb, "type MyWorldStruct struct{}", string(data))
			},
		},
		{
			name: "successful tf file generation and formatting",
			args: args{
				templatesMap: map[string]string{"test.tf": `resource    "aws_s3_bucket"    "{{.Name}}_bucket"  {}`},
				fileName:     "test.tf",
				outputFile:   path.Join(testOutput, "output.tf"),
				data:         struct{ Name string }{Name: "world"},
			},
			extraValidations: func(tb testing.TB, outputFile string, err error) {
				if err != nil {
					return
				}

				data, err := os.ReadFile(outputFile)
				require.NoError(tb, err)
				require.Equal(tb, `resource "aws_s3_bucket" "world_bucket" {}`, string(data))
			},
		},
		{
			name: "successful file generation without formatting for an unsuported ext",
			args: args{
				templatesMap: map[string]string{"test.txt": "Hello, {{.Name}}!"},
				fileName:     "test.txt",
				outputFile:   path.Join(testOutput, "output.txt"),
				data:         struct{ Name string }{Name: "World"},
			},
			extraValidations: func(tb testing.TB, outputFile string, err error) {
				if err != nil {
					return
				}

				data, err := os.ReadFile(outputFile)
				require.NoError(tb, err)
				require.Equal(tb, "Hello, World!", string(data))
			},
		},
	}

	defer func() {
		_ = os.RemoveAll(testOutput)
	}()

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			err := GenerateFile(tc.args.templatesMap, tc.args.fileName, tc.args.fileTmpl, tc.args.outputFile, tc.args.data)

			require.ErrorIs(t, err, tc.targetErr)

			if tc.extraValidations != nil {
				tc.extraValidations(t, tc.args.outputFile, err)
			}
		})
	}
}

func TestGenerateFiles(t *testing.T) {
	type args struct {
		templatesMap map[string]string
		filesMap     map[string]File
		data         any
		output       string
	}

	testOutput := "./testoutput"
	_ = os.MkdirAll(testOutput, os.ModePerm)

	tests := []struct {
		name      string
		args      args
		targetErr error
	}{
		{
			name: "generate single file",
			args: args{
				templatesMap: map[string]string{
					"template.txt": "Hello, {{.Name}}!",
				},
				filesMap: nil,
				data:     struct{ Name string }{"John"},
				output:   testOutput,
			},
			targetErr: nil,
		},
	}

	defer func() {
		_ = os.RemoveAll(testOutput)
	}()

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			err := GenerateFiles(tc.args.templatesMap, tc.args.filesMap, tc.args.data, tc.args.output)

			require.ErrorIs(t, err, tc.targetErr)
		})
	}
}
