package generators

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/ettle/strcase"

	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
	"github.com/joselitofilho/aws-terraform-generator/internal/utils"
)

var ErrUnsupportedFileType = errors.New("unsupported file type")

type TemplateMapValue struct {
	TemplateName string
	Template     []byte
}

func Build(data any, templateName, templateContent string) (string, error) {
	tmpl, err := buildAndParseTemplate(templateName, templateContent)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	// Execute the template with the data and capture the output in a buffer
	var output bytes.Buffer

	err = tmpl.Execute(&output, data)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	return output.String(), nil
}

func BuildFile(data any, templateName, templateContent, path string) error {
	tmpl, err := buildAndParseTemplate(templateName, templateContent)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	// Execute the template with the data and write the output to a file
	output, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	defer output.Close()

	err = tmpl.Execute(output, data)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

func CreateFilesMap(files []config.File) map[string]File {
	filesConf := map[string]File{}
	for i := range files {
		filesConf[files[i].Name] = File{
			Tmpl:    files[i].Tmpl,
			Imports: files[i].Imports,
		}
	}
	return filesConf
}

func GenerateFile(
	defaultTemplatesMap map[string]string, fileName, fileTmpl, outputFile string, data any,
) error {
	var (
		tmpl     string
		tmplName = fmt.Sprintf("%s-template", strings.ReplaceAll(fileName, ".", "-"))
	)

	if fileTmpl == "" {
		tmpl = defaultTemplatesMap[fileName]
	} else {
		tmpl = fileTmpl
	}

	err := BuildFile(data, tmplName, tmpl, outputFile)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	err = formatFileBasedOnExt(fileName, outputFile)
	if err != nil && !errors.Is(err, ErrUnsupportedFileType) {
		fmt.Println(err)
		err = nil
	}

	return err
}

func GenerateFiles(templatesMap map[string]TemplateMapValue, filesMap map[string]File, output string, data any) error {
	for fileName, file := range filesMap {
		var (
			tmplName string
			tmpl     string
		)

		if file.Tmpl == "" {
			if tmplMaplValue, hasValue := templatesMap[fileName]; hasValue {
				tmplName = tmplMaplValue.TemplateName
				tmpl = string(tmplMaplValue.Template)
			}
		} else {
			tmplName = fmt.Sprintf("%s-template", strings.ReplaceAll(fileName, ".", "-"))
			tmpl = file.Tmpl
		}

		outputFile := fmt.Sprintf("%s/%s", output, fileName)

		err := BuildFile(data, tmplName, tmpl, outputFile)
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		err = formatFileBasedOnExt(fileName, outputFile)
		if err != nil && !errors.Is(err, ErrUnsupportedFileType) {
			fmt.Println(err)
		}
	}

	return nil
}

// buildAndParseTemplate Create a new template and parse the template content
func buildAndParseTemplate(name, content string) (*template.Template, error) {
	tmpl, err := template.New(name).
		Funcs(template.FuncMap{
			"getFileByName":  func(files map[string]File, name string) File { return files[name] },
			"getFileImports": func(files map[string]File, name string) []string { return files[name].Imports },
			"ToCamel":        func(s string) string { return strcase.ToCamel(s) },
			"ToKebab":        func(s string) string { return strcase.ToKebab(s) },
			"ToLower":        func(s string) string { return strings.ToLower(s) },
			"ToPascal":       func(s string) string { return strcase.ToPascal(s) },
			"ToSpace":        func(s string) string { return strings.ReplaceAll(strcase.ToKebab(s), "-", " ") },
			"ToSnake":        func(s string) string { return strcase.ToSnake(s) },
			"ToUpper":        func(s string) string { return strings.ToUpper(s) },
		}).
		Parse(content)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return tmpl, nil
}

func formatFileBasedOnExt(fileName string, outputFile string) error {
	var err error

	ext := path.Ext(fileName)

	switch ext {
	case ".go":
		err = utils.GoFormat(outputFile)
	case ".tf":
		err = utils.TerraformFormat(outputFile)
	default:
		err = nil
	}

	return err
}
