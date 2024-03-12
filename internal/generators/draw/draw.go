package draw

import (
	"fmt"
	"os"
	"path"

	"github.com/joselitofilho/aws-terraform-generator/internal/fmtcolor"
	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
	generatorerrs "github.com/joselitofilho/aws-terraform-generator/internal/generators/errors"
	"github.com/joselitofilho/aws-terraform-generator/internal/generators/graphviz"
	"github.com/joselitofilho/aws-terraform-generator/internal/generators/terraform"
	"github.com/joselitofilho/aws-terraform-generator/internal/resources"
	"github.com/joselitofilho/aws-terraform-generator/internal/transformers/resourcestoyaml"
	"github.com/joselitofilho/aws-terraform-generator/internal/transformers/terraformtoresources"
	"gopkg.in/yaml.v3"
)

// https://awsicons.dev/
var DefaultResourceImageMap = config.Images{
	resources.APIGatewayType: "assets/diagram/api_gateway.svg",
	resources.CronType:       "assets/diagram/cron.svg",
	resources.DatabaseType:   "assets/diagram/database_dynamo_db.svg",
	resources.EndpointType:   "assets/diagram/endpoint.svg",
	resources.GoogleBQType:   "assets/diagram/google_bigquery.svg",
	resources.KinesisType:    "assets/diagram/kinesis_data_stream.svg",
	resources.LambdaType:     "assets/diagram/lambda.svg",
	resources.RestfulAPIType: "assets/diagram/restful_api.svg",
	resources.S3Type:         "assets/diagram/s3_bucket.svg",
	resources.SNSType:        "assets/diagram/sns.svg",
	resources.SQSType:        "assets/diagram/sqs.svg",
	resources.UnknownType:    "",
}

type Draw struct {
	workdirs       []string
	files          []string
	configFilename string
	output         string
}

func NewDraw(workdirs, files []string, configFilename, output string) *Draw {
	return &Draw{workdirs: workdirs, files: files, configFilename: configFilename, output: output}
}

func (d *Draw) Build() error {
	yamlParser := config.NewYAML(d.configFilename)

	yamlConfig, err := yamlParser.Parse()
	if err != nil {
		return fmt.Errorf("%w: %w", generatorerrs.ErrYAMLParse, err)
	}

	tfConfig, err := terraform.Parse(d.workdirs, d.files)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	resc := terraformtoresources.NewTransformer(yamlConfig, &tfConfig).Transform()

	_ = os.Mkdir(d.output, os.ModePerm)

	diagramConfig, err := resourcestoyaml.NewTransformer(yamlConfig, resc).Transform()
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	yamlFilename := "diagram"
	if yamlConfig.Draw.Name != "" {
		yamlFilename = yamlConfig.Draw.Name
	}

	yamlFilename += ".yaml"

	yamlData, err := yaml.Marshal(diagramConfig)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	yamlfile, err := os.Create(path.Join(d.output, yamlFilename))
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	defer yamlfile.Close()

	if _, err = yamlfile.WriteString(string(yamlData)); err != nil {
		return fmt.Errorf("%w", err)
	}

	fmtcolor.White.Println("The diagram yaml file has been generated successfully.")

	dotConfig := graphviz.Config{Orientation: yamlConfig.Draw.Orientation}

	resourceImageMap := mergeImages(DefaultResourceImageMap, yamlConfig.Draw.Images)

	dotContent := graphviz.Build(resc, resourceImageMap, dotConfig)

	dotFilename := "diagram"
	if yamlConfig.Draw.Name != "" {
		dotFilename = yamlConfig.Draw.Name
	}

	dotFilename += ".dot"

	dotfile, err := os.Create(path.Join(d.output, dotFilename))
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	defer dotfile.Close()

	if _, err := dotfile.WriteString(dotContent); err != nil {
		return fmt.Errorf("%w", err)
	}

	fmtcolor.White.Println("The graphviz dot file has been generated successfully.")

	return nil
}

func mergeImages(defaultImages, configImages config.Images) config.Images {
	result := defaultImages

	for k, v := range configImages {
		result[k] = v
	}

	return result
}
