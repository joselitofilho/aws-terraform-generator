package apigateway

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/ettle/strcase"

	"github.com/joselitofilho/aws-terraform-generator/internal/templates"
	"github.com/joselitofilho/aws-terraform-generator/internal/templates/config"
	"github.com/joselitofilho/aws-terraform-generator/internal/utils"
)

type APIGateway struct {
	input  string
	output string
}

func NewAPIGateway(input, output string) *APIGateway {
	return &APIGateway{input: input, output: output}
}

func (a *APIGateway) Build() error {
	yamlParser := config.NewYAML(a.input)

	yamlConfig, err := yamlParser.Parse()
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	for i := range yamlConfig.APIGateways {
		apiConf := yamlConfig.APIGateways[i]

		if apiConf.APIG {
			data := Data{
				StackName: apiConf.StackName,
				APIDomain: apiConf.APIDomain,
			}

			output := fmt.Sprintf("%s/%s/mod", a.output, apiConf.StackName)
			_ = os.MkdirAll(output, os.ModePerm)

			outputFile := fmt.Sprintf("%s/apig.tf", output)

			tmplName := "apig-tf-template"
			tmpl := string(apigTFTmpl)

			err = templates.BuildFile(data, tmplName, tmpl, outputFile)
			if err != nil {
				return fmt.Errorf("%w", err)
			}

			err = utils.TerraformFormat(output)
			if err != nil {
				fmt.Println(err)
			}

			fmt.Println("Terraform has been generated successfully")
		}

		for j := range apiConf.Lambdas {
			lambdaConf := apiConf.Lambdas[j]

			envars := map[string]string{}
			for i := range lambdaConf.Envars {
				envars = lambdaConf.Envars[i]
			}

			codeConf := map[string]templates.Code{}
			for i := range lambdaConf.Code {
				codeConf[lambdaConf.Code[i].Key] = templates.Code{
					Tmpl:    lambdaConf.Code[i].Tmpl,
					Imports: lambdaConf.Code[i].Imports,
				}
			}

			lambdaData := LambdaData{
				Name:           lambdaConf.Name,
				NameSnakeCase:  strcase.ToSnake(lambdaConf.Name),
				NamePascalCase: strcase.ToPascal(lambdaConf.Name),
				Description:    lambdaConf.Description,
				Envars:         envars,
				Verb:           lambdaConf.Verb,
				Path:           lambdaConf.Path,
				Code:           codeConf,
			}

			output := fmt.Sprintf("%s/%s/lambda/%s", a.output, apiConf.StackName, lambdaConf.Name)
			_ = os.MkdirAll(output, os.ModePerm)

			err = templates.GenerateGoFiles(defaultTemplatesMap, output, codeConf, lambdaData)
			if err != nil {
				return fmt.Errorf("%w", err)
			}

			fmt.Printf("Lambda '%s' has been generated successfully\n", lambdaData.Name)
		}
	}

	return nil
}
