package generators

import (
	_ "embed"
	"os"
	"path"
	"testing"

	templategenerators "github.com/diagram-code-generator/template/pkg/generators"
	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
	"github.com/stretchr/testify/require"
)

func TestMustGenerateFile(t *testing.T) {
	type args struct {
		tg           *templategenerators.TemplateGenerator
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
		extraValidations func(testing.TB, string)
	}{
		{
			name: "successful go file generation and formatting",
			args: args{
				tg:           NewGenerator(),
				templatesMap: map[string]string{"test.go": "type  My{{.Name}}Struct    struct   {}"},
				fileName:     "test.go",
				outputFile:   path.Join(testOutput, "output.go"),
				data:         struct{ Name string }{Name: "World"},
			},
			extraValidations: func(tb testing.TB, outputFile string) {
				data, err := os.ReadFile(outputFile)
				require.NoError(tb, err)
				require.Equal(tb, "type MyWorldStruct struct{}", string(data))
			},
		},
		{
			name: "successful go file generation using extra functions",
			args: args{
				tg: NewGenerator(),
				templatesMap: map[string]string{"lambda.go": "{{getFileByName $.Files \"lambda.go\"}} " +
					"{{ range getFileImports $.Files \"lambda.go\" }}\"{{ . }}\"{{end}}"},
				fileName:   "lambda.go",
				outputFile: path.Join(testOutput, "lambda.go"),
				data: struct{ Files map[string]config.File }{
					Files: map[string]config.File{"lambda.go": {
						Imports: []string{"context"},
						Tmpl:    "tmpl",
					}}},
			},
			extraValidations: func(tb testing.TB, outputFile string) {
				data, err := os.ReadFile(outputFile)
				require.NoError(tb, err)
				require.Equal(tb, `{ tmpl [context]} "context"`, string(data))
			},
		},
		{
			name: "when file ext is not supported should log a message and the file will not be generated",
			args: args{
				tg:           NewGenerator(),
				templatesMap: map[string]string{"test.txt": "Hello, {{.Name}}!"},
				fileName:     "test.txt",
				outputFile:   path.Join(testOutput, "output.txt"),
				data:         struct{ Name string }{Name: "World"},
			},
			extraValidations: func(tb testing.TB, outputFile string) {
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
			MustGenerateFile(
				tc.args.tg, tc.args.templatesMap, tc.args.fileName, tc.args.fileTmpl, tc.args.outputFile, tc.args.data)

			tc.extraValidations(t, tc.args.outputFile)
		})
	}
}

func TestMustGenerateFiles(t *testing.T) {
	type args struct {
		tg                  *templategenerators.TemplateGenerator
		defaultTemplatesMap map[string]string
		filesMap            map[string]File
		data                any
		output              string
	}

	testOutput := "./testoutput"
	_ = os.MkdirAll(testOutput, os.ModePerm)

	tests := []struct {
		name string
		args args
	}{
		{
			name: "generate single file",
			args: args{
				tg: NewGenerator(),
				defaultTemplatesMap: map[string]string{
					"template.txt": "Hello, {{.Name}}!",
				},
				filesMap: map[string]File{"test.go": {Tmpl: "type  My{{.Name}}Struct    struct   {}"}},
				data:     struct{ Name string }{"World"},
				output:   testOutput,
			},
		},
	}

	defer func() {
		_ = os.RemoveAll(testOutput)
	}()

	for _, tc := range tests {
		t.Run(tc.name, func(_ *testing.T) {
			MustGenerateFiles(tc.args.tg, tc.args.defaultTemplatesMap, tc.args.filesMap, tc.args.data, tc.args.output)
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
