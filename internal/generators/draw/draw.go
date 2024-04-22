package draw

import (
	"fmt"
	"os"
	"path"

	"gopkg.in/yaml.v3"

	"github.com/diagram-code-generator/resources/pkg/parser/graphviz/dot"
	hcl "github.com/joselitofilho/hcl-parser-go/pkg/parser/hcl"

	"github.com/joselitofilho/aws-terraform-generator/internal/fmtcolor"
	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
	generatorerrs "github.com/joselitofilho/aws-terraform-generator/internal/generators/errors"
	awsresources "github.com/joselitofilho/aws-terraform-generator/internal/resources"
	"github.com/joselitofilho/aws-terraform-generator/internal/transformers/resourcestoyaml"
	"github.com/joselitofilho/aws-terraform-generator/internal/transformers/terraformtoresources"
)

// DefaultResourceImageMap defines the default resource images. Images from here: https://awsicons.dev/
var DefaultResourceImageMap = config.Images{
	awsresources.APIGatewayType: "assets/diagram/api_gateway.svg",
	awsresources.CronType:       "assets/diagram/cron.svg",
	awsresources.DatabaseType:   "assets/diagram/database_dynamo_db.svg",
	awsresources.EndpointType:   "assets/diagram/endpoint.svg",
	awsresources.GoogleBQType:   "assets/diagram/google_bigquery.svg",
	awsresources.KinesisType:    "assets/diagram/kinesis_data_stream.svg",
	awsresources.LambdaType:     "assets/diagram/lambda.svg",
	awsresources.RestfulAPIType: "assets/diagram/restful_api.svg",
	awsresources.S3Type:         "assets/diagram/s3_bucket.svg",
	awsresources.SNSType:        "assets/diagram/sns.svg",
	awsresources.SQSType:        "assets/diagram/sqs.svg",
	awsresources.UnknownType:    "",
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
		return fmt.Errorf("%w: %w", generatorerrs.ErrYAMLParser, err)
	}

	tfConfig, err := hcl.Parse(d.workdirs, d.files)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	resc := terraformtoresources.NewTransformer(yamlConfig, tfConfig).Transform()

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

	nodeAttrs := make(map[string]any)
	for k, v := range dot.DefaultNodeAttrs {
		nodeAttrs[k] = v
	}

	resourceImageMap := mergeImages(DefaultResourceImageMap, yamlConfig.Draw.Images)
	dotConfig := &dot.Config{
		Direction:        yamlConfig.Draw.Direction,
		Splines:          yamlConfig.Draw.Splines,
		NodeAttrs:        nodeAttrs,
		ResourceImageMap: resourceImageMap.ToStringMap(),
	}

	dotContent := dot.NewDotDiagram(dotConfig).Build(resc)

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
