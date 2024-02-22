package templates

import (
	"bytes"
	"fmt"
	"os"
	"text/template"

	"github.com/joselitofilho/aws-terraform-generator/internal/templates/config"
)

func Build(data any, templateName, templateContent string) (string, error) {
	// Create a new template and parse the template content
	tmpl, err := template.New(templateName).Parse(string(templateContent))
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
	// Create a new template and parse the template content
	tmpl, err := template.New(templateName).Parse(string(templateContent))
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	tmpl.Funcs(template.FuncMap{
		"getFileByName": func(files map[string]File, name string) File {
			return files[name]
		},
		"getFileImports": func(files map[string]File, name string) []string {
			return files[name].Imports
		},
	})

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
