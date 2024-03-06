package tftodrawio

import (
	"fmt"
	"strings"

	"github.com/ettle/strcase"
	"github.com/joselitofilho/aws-terraform-generator/internal/drawio"
	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
	"github.com/joselitofilho/aws-terraform-generator/internal/generators/terraform"
)

var (
	envarSuffixDBHost           = "DB_HOST"
	envarSuffixKinesisStreamURL = "_KINESIS_STREAM_URL"
)

var (
	suffixKinesis = "_kinesis"
	suffixLambda  = "_lambda"
	suffixSQS     = "_sqs"
)

var (
	labelAWSLambdaEventSourceMapping = "aws_lambda_event_source_mapping"
	labelAWSKinesisStream            = "aws_kinesis_stream"
	labelAWSSQSQueue                 = "aws_sqs_queue"
)

var (
	arnKinesisKey = "kinesis"
	arnLambdaKey  = "lambda"
	arnSQSKey     = "sqs"
)

func TransformTfToDrawIO(yamlConfig *config.Config, tfConfig *terraform.Config) *drawio.ResourceCollection {
	resources := []drawio.Resource{}
	relationships := []drawio.Relationship{}

	dbResourcesByName := map[string]drawio.Resource{}
	kinesisResourcesByName := map[string]drawio.Resource{}
	lambdaResourcesByName := map[string]drawio.Resource{}
	sqsResourcesByName := map[string]drawio.Resource{}

	relationshipsMap := map[resourceARN][]resourceARN{}

	id := 1

	processTerraformModules(tfConfig.Modules, dbResourcesByName, kinesisResourcesByName, lambdaResourcesByName,
		&id, &resources, &relationships)

	processTerraformResources(tfConfig.Resources, kinesisResourcesByName, sqsResourcesByName, relationshipsMap,
		&id, &resources)

	buildRelationships(relationshipsMap, kinesisResourcesByName, lambdaResourcesByName, sqsResourcesByName,
		&relationships)

	return &drawio.ResourceCollection{Resources: resources, Relationships: relationships}
}

func buildRelationships(
	relationshipsMap map[resourceARN][]resourceARN,
	kinesisResourcesByName, lambdaResourcesByName, sqsResourcesByName map[string]drawio.Resource,
	relationships *[]drawio.Relationship,
) {
	for k, v := range relationshipsMap {
		source := getResourceByARN(k, kinesisResourcesByName, lambdaResourcesByName, sqsResourcesByName)

		for i := range v {
			target := getResourceByARN(v[i], kinesisResourcesByName, lambdaResourcesByName, sqsResourcesByName)

			*relationships = append(*relationships, drawio.Relationship{Source: source, Target: target})
		}
	}
}

func processTerraformModules(
	tfModules []*terraform.Module,
	dbResourcesByName, kinesisResourcesByName, lambdaResourcesByName map[string]drawio.Resource,
	id *int, resources *[]drawio.Resource, relationships *[]drawio.Relationship,
) {
	for _, conf := range tfModules {
		if len(conf.Labels) == 1 {
			l := conf.Labels[0]

			if strings.HasSuffix(strings.ToLower(l), suffixLambda) {
				processLambdaResource(conf,
					dbResourcesByName, kinesisResourcesByName, lambdaResourcesByName, id, resources, relationships)
			}
		}
	}
}

func processTerraformResources(
	tfResources []*terraform.Resource,
	kinesisResourcesByName, sqsResourcesByName map[string]drawio.Resource,
	relationshipsMap map[resourceARN][]resourceARN,
	id *int, resources *[]drawio.Resource,
) {
	for _, conf := range tfResources {
		if len(conf.Labels) == 2 {
			switch conf.Labels[0] {
			case labelAWSKinesisStream:
				processKinesisResource(conf, kinesisResourcesByName, id, resources)
			case labelAWSSQSQueue:
				processSQSResource(conf, sqsResourcesByName, id, resources)
			case labelAWSLambdaEventSourceMapping:
				processEventSourceMapping(conf, relationshipsMap)
			}
		}
	}
}

func processDBResourceFromEnvar(
	envar string, dbResourcesByName map[string]drawio.Resource, id *int, resources *[]drawio.Resource,
) {
	value := databaseName(envar, envarSuffixDBHost)

	if _, ok := dbResourcesByName[value]; !ok {
		dbResource := drawio.NewGenericResource(fmt.Sprintf("%d", *id), value, drawio.DatabaseType)
		*id++

		dbResourcesByName[value] = dbResource
		*resources = append(*resources, dbResource)
	}
}

func processEventSourceMapping(conf *terraform.Resource, relationshipsMap map[resourceARN][]resourceARN) {
	eventSourceARN := resourceByARN(conf.Attributes["event_source_arn"].(string))
	functionName := resourceByARN(conf.Attributes["function_name"].(string))

	relationshipsMap[eventSourceARN] = append(relationshipsMap[eventSourceARN], functionName)
}

func processKinesisResourceFromEnvar(
	envar string, kinesisResourcesByName map[string]drawio.Resource,
	id *int, resources *[]drawio.Resource,
) {
	value := databaseName(envar, envarSuffixKinesisStreamURL)

	if _, ok := kinesisResourcesByName[value]; !ok {
		kinesisResource := drawio.NewGenericResource(fmt.Sprintf("%d", *id), value, drawio.KinesisType)
		*id++

		kinesisResourcesByName[value] = kinesisResource
		*resources = append(*resources, kinesisResource)
	}
}

func processKinesisResource(
	conf *terraform.Resource, kinesisResourcesByName map[string]drawio.Resource,
	id *int, resources *[]drawio.Resource,
) {
	l := conf.Labels[1]

	if strings.HasSuffix(strings.ToLower(l), suffixKinesis) {
		value := kinesisName(l, suffixKinesis)

		resource := drawio.NewGenericResource(fmt.Sprintf("%d", *id), value, drawio.KinesisType)
		*id++

		*resources = append(*resources, resource)
		kinesisResourcesByName[value] = resource
	}
}

func processLambdaResource(
	conf *terraform.Module,
	dbResourcesByName, kinesisResourcesByName, lambdaResourcesByName map[string]drawio.Resource,
	id *int, resources *[]drawio.Resource, relationships *[]drawio.Relationship,
) {
	value := lambdaName(conf.Labels[0], suffixLambda)

	resource := drawio.NewGenericResource(fmt.Sprintf("%d", *id), value, drawio.LambdaType)
	*id++

	*resources = append(*resources, resource)
	lambdaResourcesByName[value] = resource

	for k := range conf.Attributes["lambda_function_env_vars"].(map[string]any) {
		if strings.HasSuffix(k, envarSuffixDBHost) {
			processDBResourceFromEnvar(k, dbResourcesByName, id, resources)
			*relationships = append(*relationships,
				drawio.Relationship{Source: resource, Target: dbResourcesByName[databaseName(k, envarSuffixDBHost)]})
		}

		if strings.HasSuffix(k, envarSuffixKinesisStreamURL) {
			processKinesisResourceFromEnvar(k, kinesisResourcesByName, id, resources)
			*relationships = append(*relationships, drawio.Relationship{
				Source: resource,
				Target: kinesisResourcesByName[kinesisName(k, envarSuffixKinesisStreamURL)],
			})
		}
	}
}

func processSQSResource(
	conf *terraform.Resource, sqsResourcesByName map[string]drawio.Resource,
	id *int, resources *[]drawio.Resource,
) {
	l := conf.Labels[1]

	if strings.HasSuffix(strings.ToLower(l), suffixSQS) {
		value := sqsName(l, suffixSQS)

		resource := drawio.NewGenericResource(fmt.Sprintf("%d", *id), value, drawio.SQSType)
		*id++

		*resources = append(*resources, resource)
		sqsResourcesByName[value] = resource
	}
}

///////

func getResourceByARN(
	arn resourceARN,
	kinesisResourcesByName, lambdaResourcesByName, sqsResourcesByName map[string]drawio.Resource,
) (resource drawio.Resource) {
	switch arn.key {
	case arnKinesisKey:
		resource = kinesisResourcesByName[arn.name]
	case arnLambdaKey:
		resource = lambdaResourcesByName[arn.name]
	case arnSQSKey:
		resource = sqsResourcesByName[arn.name]
	}

	return resource
}

func databaseName(str, suffix string) string {
	return strcase.ToKebab(str[:len(str)-len(suffix)])
}

func kinesisName(str, suffix string) string {
	return strcase.ToKebab(str[:len(str)-len(suffix)])
}

func lambdaName(str, suffix string) string {
	return strcase.ToCamel(str[:len(str)-len(suffix)])
}

func sqsName(str, suffix string) string {
	return strcase.ToKebab(str[:len(str)-len(suffix)])
}

////////////////////////////////////////////////////////////////////////////////

type resourceARN struct {
	key  string
	name string
}

func resourceByARN(arn string) resourceARN {
	var key, name string

	if strings.HasPrefix(arn, "arn:") {
		parts := strings.Split(arn, ":")

		key = parts[2]
		switch key {
		case arnKinesisKey:
			parts = strings.Split(arn, "/")
		}

		name = parts[len(parts)-1]
	} else {
		parts := strings.Split(arn, ".")

		key = strings.Split(parts[0], "_")[1]
		name = parts[1]
	}

	switch key {
	case arnKinesisKey:
		name = kinesisName(name, suffixKinesis)
	case arnLambdaKey:
		name = lambdaName(name, suffixLambda)
	case arnSQSKey:
		name = sqsName(name, suffixSQS)
	}

	return resourceARN{key: key, name: name}
}
