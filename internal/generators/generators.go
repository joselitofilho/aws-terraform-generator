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

func Build(data any, templateName, templateContent string) (string, error) {
	tmpl, err := buildAndParseTemplate(templateName, templateContent)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	// Execute the template with the data and capture the output in a buffer.
	var output bytes.Buffer

	err = tmpl.Execute(&output, data)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	return output.String(), nil
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

func CreateTemplatesMap(filenameTemplatesList []config.FilenameTemplateMap) map[string]string {
	templatesMap := map[string]string{}

	for i := range filenameTemplatesList {
		for filename, tmpl := range filenameTemplatesList[i] {
			templatesMap[filename] = tmpl
		}
	}

	return templatesMap
}

func FilterTemplatesMap(filter string, templatesMap map[string]string) map[string]string {
	filtred := map[string]string{}

	for filename, tmpl := range templatesMap {
		if strings.Contains(filename, filter) {
			filtred[filename] = tmpl
		}
	}

	return filtred
}

func GenerateFile(templatesMap map[string]string, fileName, fileTmpl, outputFile string, data any) error {
	var (
		tmpl     string
		tmplName = fmt.Sprintf("%s-template", strings.ReplaceAll(fileName, ".", "-"))
	)

	if fileTmpl == "" {
		tmpl = templatesMap[fileName]
	} else {
		tmpl = fileTmpl
	}

	err := buildFile(data, tmplName, tmpl, outputFile)
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

func GenerateFiles(templatesMap map[string]string, filesMap map[string]File, data any, output string) error {
	mergedTemplates := map[string]string{}

	for filename, tmpl := range templatesMap {
		mergedTemplates[filename] = tmpl
	}

	for filename, file := range filesMap {
		mergedTemplates[filename] = file.Tmpl
	}

	for filename, fileTmpl := range mergedTemplates {
		tmplName := fmt.Sprintf("%s-template", strings.ReplaceAll(filename, ".", "-"))

		outputFile := path.Join(output, filename)

		err := buildFile(data, tmplName, fileTmpl, outputFile)
		if err != nil {
			// TODO: Append error
			fmt.Println(err)
		}

		err = formatFileBasedOnExt(filename, outputFile)
		if err != nil && !errors.Is(err, ErrUnsupportedFileType) {
			// TODO: Append error
			fmt.Println(err)
		}
	}

	// TODO: Return errors
	return nil
}

// buildAndParseTemplate Create a new template and parse the template content.
func buildAndParseTemplate(name, content string) (*template.Template, error) {
	tmpl, err := template.New(name).
		Funcs(template.FuncMap{
			"getFileByName":  func(files map[string]File, name string) File { return files[name] },
			"getFileImports": func(files map[string]File, name string) []string { return files[name].Imports },
			"ToCamel":        strcase.ToCamel,
			"ToKebab":        strcase.ToKebab,
			"ToLower":        strings.ToLower,
			"ToPascal":       strcase.ToPascal,
			"ToSpace":        func(s string) string { return strings.ReplaceAll(strcase.ToKebab(s), "-", " ") },
			"ToSnake":        strcase.ToSnake,
			"ToUpper":        strings.ToUpper,
		}).
		Parse(content)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return tmpl, nil
}

func buildFile(data any, templateName, templateContent, outputPath string) error {
	tmpl, err := buildAndParseTemplate(templateName, templateContent)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	// Execute the template with the data and write the output to a file.
	output, err := os.Create(outputPath)
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

func formatFileBasedOnExt(fileName, outputFile string) error {
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
