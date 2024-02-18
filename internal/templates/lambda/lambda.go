package lambda

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/ettle/strcase"

	"github.com/joselitofilho/aws-terraform-generator/internal/templates"
	"github.com/joselitofilho/aws-terraform-generator/internal/templates/config"
	"github.com/joselitofilho/aws-terraform-generator/internal/utils"
)

type Lambda struct {
	input  string
	output string
}

func NewLambda(input, output string) *Lambda {
	return &Lambda{input: input, output: output}
}

func (l *Lambda) Build() error {
	yamlParser := config.NewYAML(l.input)

	yamlConfig, err := yamlParser.Parse()
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	lambdaConf := yamlConfig.Lambdas[0] // TODO: Loop

	envars := map[string]string{}
	for i := range lambdaConf.Envars {
		envars = lambdaConf.Envars[i]
	}

	sqsTriggers := make([]SQSTrigger, len(lambdaConf.SQSTriggers))
	for i := range lambdaConf.SQSTriggers {
		sqsTriggers[i] = SQSTrigger{
			SourceARN: lambdaConf.SQSTriggers[i].SourceARN,
		}
	}

	crons := make([]Cron, len(lambdaConf.Cron))
	for i := range lambdaConf.Cron {
		crons[i] = Cron{
			ScheduleExpression: lambdaConf.Cron[i].ScheduleExpression,
			IsEnabled:          lambdaConf.Cron[i].IsEnabled,
		}
	}

	codeConf := map[string]Code{}
	for i := range lambdaConf.Code {
		codeConf[lambdaConf.Code[i].Key] = Code{
			Tmpl:    lambdaConf.Code[i].Tmpl,
			Imports: lambdaConf.Code[i].Imports,
		}
	}

	data := Data{
		Name:           lambdaConf.Name,
		NameSnakeCase:  strcase.ToSnake(lambdaConf.Name),
		NamePascalCase: strcase.ToPascal(lambdaConf.Name),
		Description:    lambdaConf.Description,
		Envars:         envars,
		SQSTriggers:    sqsTriggers,
		Crons:          crons,
		Code:           codeConf,
	}

	tmplName := "lambda-tf-template"

	output := fmt.Sprintf("%s/mod", l.output)
	_ = os.MkdirAll(output, os.ModePerm)

	outputFile := fmt.Sprintf("%s/%s.tf", output, lambdaConf.Name)

	tmpl := string(lambdaTFTmpl)

	err = templates.BuildFile(data, tmplName, tmpl, outputFile)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	fmt.Println("Terraform has been generated successfully")

	err = utils.TerraformFormat(output)
	if err != nil {
		fmt.Println(err)
	}

	output = fmt.Sprintf("%s/lambda/%s", l.output, lambdaConf.Name)
	_ = os.MkdirAll(output, os.ModePerm)

	err = generateGoFiles(output, codeConf, data)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	fmt.Println("Lambda has been generated successfully")

	return nil
}
