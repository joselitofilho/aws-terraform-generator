package apigateway

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/ettle/strcase"

	templates "github.com/joselitofilho/aws-terraform-generator/internal/generators"
	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
)

type APIGateway struct {
	config string
	output string
}

func NewAPIGateway(config, output string) *APIGateway {
	return &APIGateway{config: config, output: output}
}

func (a *APIGateway) Build() error {
	yamlParser := config.NewYAML(a.config)

	yamlConfig, err := yamlParser.Parse()
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	defaultTmplsMap := map[string]string{
		"apig.tf":   string(apigTFTmpl),
		"lambda.tf": string(lambdaTFTmpl),
	}

	for i := range yamlConfig.APIGateways {
		apiConf := yamlConfig.APIGateways[i]

		output := fmt.Sprintf("%s/%s/mod", a.output, apiConf.StackName)
		_ = os.MkdirAll(output, os.ModePerm)

		if apiConf.APIG {
			fileName := "apig.tf"
			tmpl := string(apigTFTmpl)
			outputFile := fmt.Sprintf("%s/%s", output, fileName)

			data := Data{
				StackName: apiConf.StackName,
				APIDomain: apiConf.APIDomain,
			}

			err = templates.GenerateFile(defaultTmplsMap, fileName, tmpl, outputFile, data)
			if err != nil {
				return fmt.Errorf("%w", err)
			}

			fmt.Printf("Terraform '%s' has been generated successfully\n", fileName)
		}

		for j := range apiConf.Lambdas {
			lambdaConf := apiConf.Lambdas[j]

			envars := map[string]string{}
			for i := range lambdaConf.Envars {
				for key, value := range lambdaConf.Envars[i] {
					envars[key] = value
				}
			}

			filesConf := templates.CreateFilesMap(lambdaConf.Files)

			lambdaData := LambdaData{
				ModuleLambdaSource: lambdaConf.Source,
				StackName:          apiConf.StackName,
				Name:               lambdaConf.Name,
				NameSnakeCase:      strcase.ToSnake(lambdaConf.Name),
				NamePascalCase:     strcase.ToPascal(lambdaConf.Name),
				Description:        lambdaConf.Description,
				Envars:             envars,
				Verb:               lambdaConf.Verb,
				Path:               lambdaConf.Path,
				Files:              filesConf,
			}

			fileName := fmt.Sprintf("%s.tf", lambdaConf.Name)
			tmpl := string(lambdaTFTmpl)
			outputFile := fmt.Sprintf("%s/%s", output, fileName)

			err = templates.GenerateFile(defaultTmplsMap, fileName, tmpl, outputFile, lambdaData)
			if err != nil {
				return fmt.Errorf("%w", err)
			}

			fmt.Printf("Terraform '%s' has been generated successfully\n", fileName)

			output := fmt.Sprintf("%s/%s/lambda/%s", a.output, apiConf.StackName, lambdaConf.Name)
			_ = os.MkdirAll(output, os.ModePerm)

			err = templates.GenerateFiles(defaultTemplatesMap, filesConf, output, lambdaData)
			if err != nil {
				return fmt.Errorf("%w", err)
			}

			fmt.Printf("Lambda '%s' has been generated successfully\n", lambdaData.Name)
		}
	}

	return nil
}