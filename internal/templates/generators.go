package templates

import (
	_ "embed"
	"errors"
	"fmt"
	"path"
	"strings"

	"github.com/joselitofilho/aws-terraform-generator/internal/utils"
)

var ErrUnsupportedFileType = errors.New("unsupported file type")

type TemplateMapValue struct {
	TemplateName string
	Template     []byte
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
