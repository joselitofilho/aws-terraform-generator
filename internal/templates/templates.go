package templates

import (
	"bytes"
	"fmt"
	"os"
	"text/template"
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
