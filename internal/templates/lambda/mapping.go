package lambda

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/joselitofilho/aws-terraform-generator/internal/templates"
	"github.com/joselitofilho/aws-terraform-generator/internal/utils"
)

type templateMapValue struct {
	templateName string
	template     []byte
}

var (
	//go:embed tmpls/lambda.tf.tmpl
	lambdaTFTmpl []byte

	//go:embed tmpls/lambda.go.tmpl
	lambdaGoTmpl []byte

	//go:embed tmpls/main.go.tmpl
	mainGoTmpl []byte
)

func generateGoFiles(output string, codeConf map[string]Code, data any) error {
	defaultTemplatesMap := map[string]templateMapValue{
		"main":   {templateName: "main-go-template", template: mainGoTmpl},
		"lambda": {templateName: "lambda-go-template", template: lambdaGoTmpl},
	}

	for tmplContext, tmplMapValue := range defaultTemplatesMap {
		if _, ok := codeConf[tmplContext]; !ok {
			codeConf[tmplContext] = Code{Tmpl: string(tmplMapValue.template)}
		}
	}

	for tmplContext, tmplCode := range codeConf {
		var (
			tmplName string
			tmpl     string
		)

		tmplMapValue, ok := defaultTemplatesMap[tmplContext]
		if ok {
			tmplName = tmplMapValue.templateName
			tmpl = string(tmplMapValue.template)

			if len(tmplCode.Tmpl) > 0 {
				tmpl = tmplCode.Tmpl
			}
		} else {
			tmplName = fmt.Sprintf("%s-go-template", tmplContext)
			tmpl = tmplCode.Tmpl
		}

		outputFile := fmt.Sprintf("%s/%s.go", output, strings.ToLower(tmplContext))

		err := templates.BuildFile(data, tmplName, tmpl, outputFile)
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
