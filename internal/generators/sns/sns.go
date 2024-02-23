package sns

import (
	_ "embed"
	"fmt"
	"strings"

	templates "github.com/joselitofilho/aws-terraform-generator/internal/generators"
	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
	"github.com/joselitofilho/aws-terraform-generator/internal/utils"
)

//go:embed tmpls/sns.tf.tmpl
var snsTFTmpl []byte

type Data struct {
	Name       string
	BucketName string
	Lambdas    []LambdaSNSData
}

type LambdaSNSData struct {
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

		lambdaEvents := make([]LambdaSNSData, 0, len(conf.Lambdas))

		for _, lambda := range conf.Lambdas {
			levt := LambdaSNSData{
				Name:         lambda.Name,
				Events:       fmt.Sprintf(`"%s"`, strings.Join(lambda.Events, ", ")),
				FilterPrefix: lambda.FilterPrefix,
				FilterSuffix: lambda.FilterSuffix,
			}
			lambdaEvents = append(lambdaEvents, levt)
		}

		data.Lambdas = lambdaEvents

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
