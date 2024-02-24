package sns

import (
	_ "embed"
	"fmt"
	"path/filepath"
	"strings"

	templates "github.com/joselitofilho/aws-terraform-generator/internal/generators"
	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
	"github.com/joselitofilho/aws-terraform-generator/internal/utils"
)

type Data struct {
	Name       string
	BucketName string
	Lambdas    []SNSResourceData
	SQSs       []SNSResourceData
}

type SNSResourceData struct {
	Name         string
	Events       string
	FilterPrefix string
	FilterSuffix string
}

type SNS struct {
	config string
	output string
}

func NewSNS(config, output string) *SNS {
	return &SNS{config: config, output: output}
}

func (s *SNS) Build() error {
	yamlParser := config.NewYAML(s.config)

	yamlConfig, err := yamlParser.Parse()
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	tmplName := "sns-tf-template"
	result := ""

	for i := range yamlConfig.SNSs {
		conf := yamlConfig.SNSs[i]

		data := Data{
			Name:       conf.Name,
			BucketName: conf.BucketName,
		}

		if len(conf.Files) > 0 {
			filesConf := templates.CreateFilesMap(conf.Files)

			err = templates.GenerateFiles(defaultTemplatesMap, filesConf, filepath.Dir(s.output), data)
			if err != nil {
				return fmt.Errorf("%w", err)
			}

			fmt.Printf("SNS '%s' has been generated successfully\n", conf.Name)

			continue
		}

		lambdaEvents := make([]SNSResourceData, 0, len(conf.Lambdas))

		for _, lambda := range conf.Lambdas {
			evt := SNSResourceData{
				Name:         lambda.Name,
				Events:       fmt.Sprintf(`"%s"`, strings.Join(lambda.Events, ", ")),
				FilterPrefix: lambda.FilterPrefix,
				FilterSuffix: lambda.FilterSuffix,
			}
			lambdaEvents = append(lambdaEvents, evt)
		}

		data.Lambdas = lambdaEvents

		sqsEvents := make([]SNSResourceData, 0, len(conf.SQSs))

		for _, sqs := range conf.SQSs {
			evt := SNSResourceData{
				Name:         sqs.Name,
				Events:       fmt.Sprintf(`"%s"`, strings.Join(sqs.Events, ", ")),
				FilterPrefix: sqs.FilterPrefix,
				FilterSuffix: sqs.FilterSuffix,
			}
			sqsEvents = append(sqsEvents, evt)
		}

		data.SQSs = sqsEvents

		output, err := templates.Build(data, tmplName, string(snsTFTmpl))
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		result = fmt.Sprintf("%s\n%s", result, output)
	}

	if result != "" {
		err = templates.BuildFile(Data{}, tmplName, result, s.output)
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		err = utils.TerraformFormat(s.output)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("SNS has been generated successfully")
	}

	return nil
}
