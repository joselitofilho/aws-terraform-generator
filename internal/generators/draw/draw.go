package draw

import (
	"fmt"
	"os"
	"path"

	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
	generatorerrs "github.com/joselitofilho/aws-terraform-generator/internal/generators/errors"
	"github.com/joselitofilho/aws-terraform-generator/internal/generators/graphviz"
	"github.com/joselitofilho/aws-terraform-generator/internal/generators/terraform"
	"github.com/joselitofilho/aws-terraform-generator/internal/resources"
	"github.com/joselitofilho/aws-terraform-generator/internal/transformers/tftodrawio"
)

// https://awsicons.dev/
var defaultResourceImageMap = map[resources.ResourceType]string{
	resources.APIGatewayType: "assets/aws/api_gateway.svg",
	resources.CronType:       "assets/aws/cron.svg",
	resources.DatabaseType:   "assets/aws/database_dynamo_db.svg",
	resources.EndpointType:   "assets/aws/endpoint.svg",
	resources.GoogleBQType:   "assets/aws/google_bigquery.svg",
	resources.KinesisType:    "assets/aws/kinesis_data_stream.svg",
	resources.LambdaType:     "assets/aws/lambda.svg",
	resources.RestfulAPIType: "assets/aws/restful_api.svg",
	resources.S3Type:         "assets/aws/s3_bucket.svg",
	resources.SNSType:        "assets/aws/sns.svg",
	resources.SQSType:        "assets/aws/sqs.svg",
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

	resc := tftodrawio.NewTransformer(yamlConfig, &tfConfig).Transform()

	dotConfig := graphviz.Config{Orientation: yamlConfig.Draw.Orientation}

	dotContent, err := graphviz.Build(resc, defaultResourceImageMap, dotConfig)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	_ = os.Mkdir(d.output, os.ModePerm)

	dotfile, err := os.Create(path.Join(d.output, "diagram.dot"))
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	defer dotfile.Close()

	if _, err := dotfile.WriteString(dotContent); err != nil {
		return fmt.Errorf("%w", err)
	}

	fmt.Println("The graphviz dot file has been generated successfully.")

	return nil
}
