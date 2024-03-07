package draw

import (
	"fmt"
	"os"
	"path"

	"github.com/joselitofilho/aws-terraform-generator/internal/drawio"
	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
	generatorerrs "github.com/joselitofilho/aws-terraform-generator/internal/generators/errors"
	"github.com/joselitofilho/aws-terraform-generator/internal/generators/graphviz"
	"github.com/joselitofilho/aws-terraform-generator/internal/generators/terraform"
	"github.com/joselitofilho/aws-terraform-generator/internal/transformers/tftodrawio"
)

// https://awsicons.dev/
var defaultResourceImageMap = map[drawio.ResourceType]string{
	drawio.APIGatewayType: "assets/aws/api_gateway.svg",
	drawio.CronType:       "assets/aws/cron.svg",
	drawio.DatabaseType:   "assets/aws/database_dynamo_db.svg",
	drawio.EndpointType:   "assets/aws/endpoint.svg",
	drawio.GoogleBQType:   "assets/aws/google_bigquery.svg",
	drawio.KinesisType:    "assets/aws/kinesis_data_stream.svg",
	drawio.LambdaType:     "assets/aws/lambda.svg",
	drawio.RestfulAPIType: "assets/aws/restful_api.svg",
	drawio.S3Type:         "assets/aws/s3_bucket.svg",
	drawio.SNSType:        "assets/aws/sns.svg",
	drawio.SQSType:        "assets/aws/sqs.svg",
	drawio.UnknownType:    "",
}

type Draw struct {
	configFilename string
	workdir        string
	output         string
}

func NewDraw(configFilename, workdir, output string) *Draw {
	return &Draw{configFilename: configFilename, workdir: workdir, output: output}
}

func (d *Draw) Build() error {
	yamlParser := config.NewYAML(d.configFilename)

	yamlConfig, err := yamlParser.Parse()
	if err != nil {
		return fmt.Errorf("%w: %w", generatorerrs.ErrYAMLParse, err)
	}

	tfConfig, err := terraform.Parse(d.workdir)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	resources := tftodrawio.NewTransformer(yamlConfig, &tfConfig).Transform()

	dotContent, err := graphviz.Build(resources, defaultResourceImageMap)
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
