package templates

import (
	_ "embed"
	"errors"
	"fmt"
	"strings"

	"github.com/joselitofilho/aws-terraform-generator/internal/utils"
)

var ErrUnsupportedFileType = errors.New("unsupported file type")

type TemplateMapValue struct {
	TemplateName string
	Template     []byte
}

func GenerateFiles(
	defaultTemplatesMap map[string]string, fileName, fileTmpl string, data any, outputFile string,
) error {
	var (
		tmpl     string
		tmplName = fmt.Sprintf("%s-template", strings.ReplaceAll(fileName, ".", ""))
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

	ext := strings.Split(fileName, ".")[1]

	switch ext {
	case "go":
		err = utils.GoFormat(outputFile)
	case "tf":
		err = utils.TerraformFormat(outputFile)
	case "hcl":
		err = nil
	default:
		err = ErrUnsupportedFileType
	}

	if err != nil && !errors.Is(err, ErrUnsupportedFileType) {
		fmt.Println(err)
		err = nil
	}

	return err
}

func GenerateGoFiles(
	defaultTemplatesMap map[string]TemplateMapValue, output string, codeConf map[string]Code, data any,
) error {
	for tmplContext, tmplMapValue := range defaultTemplatesMap {
		if _, ok := codeConf[tmplContext]; !ok {
			codeConf[tmplContext] = Code{Tmpl: string(tmplMapValue.Template)}
		}
	}

	for tmplContext, tmplCode := range codeConf {
		var (
			tmplName string
			tmpl     string
		)

		tmplMapValue, ok := defaultTemplatesMap[tmplContext]
		if ok {
			tmplName = tmplMapValue.TemplateName
			tmpl = string(tmplMapValue.Template)

			if len(tmplCode.Tmpl) > 0 {
				tmpl = tmplCode.Tmpl
			}
		} else {
			tmplName = fmt.Sprintf("%s-go-template", tmplContext)
			tmpl = tmplCode.Tmpl
		}

		outputFile := fmt.Sprintf("%s/%s.go", output, strings.ToLower(tmplContext))

		err := BuildFile(data, tmplName, tmpl, outputFile)
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		err = utils.GoFormat(outputFile)
		if err != nil {
			fmt.Println(err)
		}
	}

	return nil
}
