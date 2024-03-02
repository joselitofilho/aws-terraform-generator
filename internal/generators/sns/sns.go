package sns

import (
	_ "embed"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/joselitofilho/aws-terraform-generator/internal/generators"
	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
	generatorserrs "github.com/joselitofilho/aws-terraform-generator/internal/generators/errors"
	"github.com/joselitofilho/aws-terraform-generator/internal/utils"
)

type Data struct {
	Name       string
	BucketName string
	Lambdas    []ResourceData
	SQSs       []ResourceData
}

type ResourceData struct {
	Name         string
	Events       string
	FilterPrefix string
	FilterSuffix string
}

type SNS struct {
	configFileName string
	output         string
}

func NewSNS(configFileName, output string) *SNS {
	return &SNS{configFileName: configFileName, output: output}
}

func (s *SNS) Build() error {
	yamlParser := config.NewYAML(s.configFileName)

	yamlConfig, err := yamlParser.Parse()
	if err != nil {
		return fmt.Errorf("%w: %w", generatorserrs.ErrYAMLParse, err)
	}

	modPath := path.Join(s.output, "mod")
	_ = os.MkdirAll(modPath, os.ModePerm)

	tmplName := "sns-tf-template"
	result := make([]string, 0, len(yamlConfig.SNSs))

	for i := range yamlConfig.SNSs {
		conf := yamlConfig.SNSs[i]

		data := Data{
			Name:       conf.Name,
			BucketName: conf.BucketName,
		}

		lambdaEvents := make([]ResourceData, 0, len(conf.Lambdas))

		for _, lambda := range conf.Lambdas {
			lambdaEvents = append(lambdaEvents, ResourceData{
				Name:         lambda.Name,
				Events:       fmt.Sprintf("%q", strings.Join(lambda.Events, ", ")),
				FilterPrefix: lambda.FilterPrefix,
				FilterSuffix: lambda.FilterSuffix,
			})
		}

		data.Lambdas = lambdaEvents

		sqsEvents := make([]ResourceData, 0, len(conf.SQSs))

		for _, sqs := range conf.SQSs {
			evt := ResourceData{
				Name:         sqs.Name,
				Events:       fmt.Sprintf("%q", strings.Join(sqs.Events, ", ")),
				FilterPrefix: sqs.FilterPrefix,
				FilterSuffix: sqs.FilterSuffix,
			}
			sqsEvents = append(sqsEvents, evt)
		}

		data.SQSs = sqsEvents

		if len(conf.Files) > 0 {
			filesConf := generators.CreateFilesMap(conf.Files)

			err = generators.GenerateFiles(nil, filesConf, data, modPath)
			if err != nil {
				return fmt.Errorf("%w", err)
			}

			fmt.Printf("SNS '%s' has been generated successfully\n", conf.Name)
		}

		output, err := generators.Build(data, tmplName, string(snsTFTmpl))
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		result = append(result, output)
	}

	if len(result) > 0 {
		outputFile := path.Join(modPath, "sns.tf")

		if err := generators.BuildFile(Data{}, tmplName, strings.Join(result, "\n"), outputFile); err != nil {
			return fmt.Errorf("%w", err)
		}

		if err := utils.TerraformFormat(outputFile); err != nil {
			fmt.Println(err)
		}

		fmt.Println("SNS has been generated successfully")
	}

	return nil
}
