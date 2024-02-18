package lambda

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/ettle/strcase"

	"github.com/joselitofilho/aws-terraform-generator/internal/inputs/yaml"
	"github.com/joselitofilho/aws-terraform-generator/internal/templates"
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
	yamlParser := yaml.NewYAML(l.input)

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

	data := Data{
		Name:           lambdaConf.Name,
		NameSnakeCase:  strcase.ToSnake(lambdaConf.Name),
		NamePascalCase: strcase.ToPascal(lambdaConf.Name),
		Description:    lambdaConf.Description,
		Envars:         envars,
		SQSTriggers:    sqsTriggers,
		Crons:          crons,
		Code:           Code{Lambda: LambdaCode{Imports: lambdaConf.Code.Lambda.Imports}},
	}

	tmplName := "lambda-tf-template"

	output := fmt.Sprintf("%s/mod", l.output)
	_ = os.MkdirAll(output, os.ModePerm)

	outputFile := fmt.Sprintf("%s/%s.tf", output, lambdaConf.Name)

	err = templates.BuildFile(data, tmplName, string(lambdaTFTmpl), outputFile)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	fmt.Println("Terraform has been generated successfully")

	err = utils.TerraformFormat(output)
	if err != nil {
		fmt.Println(err)
	}

	tmplName = "lambda-go-template"

	output = fmt.Sprintf("%s/lambda/%s", l.output, lambdaConf.Name)
	_ = os.MkdirAll(output, os.ModePerm)

	outputFile = fmt.Sprintf("%s/main.go", output)

	err = templates.BuildFile(data, tmplName, string(mainGoTmpl), outputFile)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	err = utils.GoFormat(outputFile)
	if err != nil {
		fmt.Println(err)
	}

	outputFile = fmt.Sprintf("%s/lambda.go", output)

	err = templates.BuildFile(data, tmplName, string(lambdaGoTmpl), outputFile)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	err = utils.GoFormat(outputFile)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Lambda has been generated successfully")

	return nil
}
