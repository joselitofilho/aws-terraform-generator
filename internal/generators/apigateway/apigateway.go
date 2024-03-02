package apigateway

import (
	_ "embed"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/joselitofilho/aws-terraform-generator/internal/generators"
	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
	generatorerrs "github.com/joselitofilho/aws-terraform-generator/internal/generators/errors"
)

type APIGateway struct {
	configFileName string
	output         string
}

func NewAPIGateway(configFileName, output string) *APIGateway {
	return &APIGateway{configFileName: configFileName, output: output}
}

func (a *APIGateway) Build() error {
	yamlParser := config.NewYAML(a.configFileName)

	yamlConfig, err := yamlParser.Parse()
	if err != nil {
		return fmt.Errorf("%w: %w", generatorerrs.ErrYAMLParse, err)
	}

	for i := range yamlConfig.APIGateways {
		apiConf := yamlConfig.APIGateways[i]

		outputMod := path.Join(a.output, apiConf.StackName, "mod")
		_ = os.MkdirAll(outputMod, os.ModePerm)

		if apiConf.APIG {
			outputFile := path.Join(outputMod, filenameTfAPIG)

			data := Data{
				StackName: apiConf.StackName,
				APIDomain: apiConf.APIDomain,
			}

			err = generators.GenerateFile(defaultTfTemplateFiles, filenameTfAPIG, "", outputFile, data)
			if err != nil {
				return fmt.Errorf("%w", err)
			}

			fmt.Printf("Terraform '%s' has been generated successfully\n", filenameTfAPIG)
		}

		for j := range apiConf.Lambdas {
			lambdaConf := apiConf.Lambdas[j]

			envars := map[string]string{}

			for i := range lambdaConf.Envars {
				for key, value := range lambdaConf.Envars[i] {
					envars[key] = value
				}
			}

			filesConf := generators.CreateFilesMap(lambdaConf.Files)

			asModule := strings.Contains(lambdaConf.Source, "git@")

			roleName := lambdaConf.RoleName
			if roleName == "" {
				roleName = "iam_for_lambda"
			}

			lambdaData := LambdaData{
				Name:        lambdaConf.Name,
				AsModule:    asModule,
				Source:      lambdaConf.Source,
				RoleName:    roleName,
				Runtime:     lambdaConf.Runtime,
				StackName:   apiConf.StackName,
				Description: lambdaConf.Description,
				Envars:      envars,
				Verb:        lambdaConf.Verb,
				Path:        lambdaConf.Path,
				Files:       filesConf,
			}

			fileName := fmt.Sprintf("%s.tf", lambdaConf.Name)

			err = generators.GenerateFiles(map[string]string{fileName: string(lambdaTFTmpl)}, nil, lambdaData, outputMod)
			if err != nil {
				return fmt.Errorf("%w", err)
			}

			fmt.Printf("Terraform '%s.tf' has been generated successfully\n", fileName)

			outputLambda := path.Join(a.output, apiConf.StackName, "lambda", lambdaConf.Name)
			_ = os.MkdirAll(outputLambda, os.ModePerm)

			err = generators.GenerateFiles(defaultGoTemplateFiles, filesConf, lambdaData, outputLambda)
			if err != nil {
				return fmt.Errorf("%w", err)
			}

			fmt.Printf("Lambda '%s' has been generated successfully\n", lambdaData.Name)
		}
	}

	return nil
}
