package templates

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/joselitofilho/aws-terraform-generator/internal/utils"
)

type TemplateMapValue struct {
	TemplateName string
	Template     []byte
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
