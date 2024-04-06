package apigateway

import (
	_ "embed"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/joselitofilho/aws-terraform-generator/internal/fmtcolor"
	"github.com/joselitofilho/aws-terraform-generator/internal/generators"
	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
	generatorerrs "github.com/joselitofilho/aws-terraform-generator/internal/generators/errors"
	"github.com/joselitofilho/aws-terraform-generator/internal/utils"
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
		return fmt.Errorf("%w: %s", generatorerrs.ErrYAMLParser, err)
	}

	apigTfTemplate := utils.MergeStringMap(map[string]string{filenameTfAPIG: string(tmplAPIGtf)},
		generators.FilterTemplatesMap(filenameTfAPIG,
			generators.CreateTemplatesMap(yamlConfig.OverrideDefaultTemplates.APIGateway)),
	)[filenameTfAPIG]

	lambdaTfTemplate := utils.MergeStringMap(map[string]string{filenameTfLambda: string(tmplLambdaTf)},
		generators.FilterTemplatesMap(
			filenameTfLambda, generators.CreateTemplatesMap(yamlConfig.OverrideDefaultTemplates.APIGateway)),
	)[filenameTfLambda]

	goTemplates := utils.MergeStringMap(defaultGoTemplateFiles, generators.FilterTemplatesMap(".go",
		generators.CreateTemplatesMap(yamlConfig.OverrideDefaultTemplates.APIGateway)))

	apigHasAlreadyGeneratedByStack := map[string]struct{}{}

	tg := generators.NewGenerator()

	for i := range yamlConfig.APIGateways {
		apiConf := yamlConfig.APIGateways[i]
		stackName := apiConf.StackName

		outputMod := path.Join(a.output, stackName, "mod")
		_ = os.MkdirAll(outputMod, os.ModePerm)

		if _, ok := apigHasAlreadyGeneratedByStack[stackName]; !ok && apiConf.APIG {
			apigHasAlreadyGeneratedByStack[stackName] = struct{}{}

			outputFile := path.Join(outputMod, filenameTfAPIG)

			data := Data{
				StackName: stackName,
				APIDomain: apiConf.APIDomain,
			}

			generators.MustGenerateFile(tg, nil, filenameTfAPIG, apigTfTemplate, outputFile, data)

			fmtcolor.White.Printf("Terraform '%s' has been generated successfully\n", filenameTfAPIG)
		}

		for j := range apiConf.Lambdas {
			buildLambdaFiles(&apiConf.Lambdas[j], apiConf.StackName, lambdaTfTemplate, outputMod, a.output,
				goTemplates)
		}
	}

	return nil
}

func buildLambdaFiles(lambdaConf *config.APIGatewayLambda, stackName, lambdaTfTemplate, outputMod, output string,
	goTemplates map[string]string,
) {
	tg := generators.NewGenerator()

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
		StackName:   stackName,
		Description: lambdaConf.Description,
		Envars:      lambdaConf.Envars,
		Verb:        lambdaConf.Verb,
		Path:        lambdaConf.Path,
		Files:       filesConf,
	}

	fileName := fmt.Sprintf("%s.tf", lambdaConf.Name)
	outputLambdaTfFile := path.Join(outputMod, fileName)

	generators.MustGenerateFile(tg, nil, fileName, lambdaTfTemplate, outputLambdaTfFile, lambdaData)

	fmtcolor.White.Printf("Terraform '%s.tf' has been generated successfully\n", fileName)

	outputLambda := path.Join(output, stackName, "lambda", lambdaConf.Name)
	_ = os.MkdirAll(outputLambda, os.ModePerm)

	generators.MustGenerateFiles(tg, goTemplates, filesConf, lambdaData, outputLambda)

	fmtcolor.White.Printf("Lambda '%s' has been generated successfully\n", lambdaData.Name)
}
