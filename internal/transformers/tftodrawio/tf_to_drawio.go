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
	envarSuffixDBHost           = "_DB_HOST"
	envarSuffixGoogleBQ         = "_BQ_PROJECT_ID"
	envarSuffixKinesisStreamURL = "_KINESIS_STREAM_URL"
	envarSuffixSQSQueueURL      = "_SQS_QUEUE_URL"
	envarSuffixRestfulAPI       = "_API_BASE_URL"
)

var (
	suffixKinesis = "_kinesis"
	suffixLambda  = "_lambda"
	suffixSQS     = "_sqs"
)

var (
	labelAWSAPIGatewayRoute          = "aws_apigatewayv2_route"
	labelAWSAPIGatewayIntegration    = "aws_apigatewayv2_integration"
	labelAWSCloudwatchEventTarget    = "aws_cloudwatch_event_target"
	labelAWSCron                     = "aws_cloudwatch_event_rule"
	labelAWSEndpoint                 = "aws_apigatewayv2_domain_name"
	labelAWSKinesisStream            = "aws_kinesis_stream"
	labelAWSLambdaEventSourceMapping = "aws_lambda_event_source_mapping"
	labelAWSSQSQueue                 = "aws_sqs_queue"
)

var (
	arnAPIGateway    = "route"
	arnCloudwatchKey = "cloudwatch"
	arnEndpoint      = "apigatewayv2"
	arnKinesisKey    = "kinesis"
	arnLambdaKey     = "lambda"
	arnSQSKey        = "sqs"
)

func TransformTfToDrawIO(yamlConfig *config.Config, tfConfig *terraform.Config) *drawio.ResourceCollection {
	resources := []drawio.Resource{}
	relationships := []drawio.Relationship{}

	apiGatewayResourcesByName := map[string]drawio.Resource{}
	cronResourcesByName := map[string]drawio.Resource{}
	dbResourcesByName := map[string]drawio.Resource{}
	endpointResourcesByName := map[string]drawio.Resource{}
	googleBQResourcesByName := map[string]drawio.Resource{}
	kinesisResourcesByName := map[string]drawio.Resource{}
	lambdaResourcesByName := map[string]drawio.Resource{}
	restfulAPIResourcesByName := map[string]drawio.Resource{}
	sqsResourcesByName := map[string]drawio.Resource{}

	endpointAPIGatewayMap := map[resourceARN]resourceARN{}
	resourceEndpointAPIGatewayMap := map[resourceARN]map[resourceARN]resourceARN{}

	relationshipsMap := map[resourceARN][]resourceARN{}

	id := 1

	processTerraformModules(tfConfig.Modules,
		dbResourcesByName, googleBQResourcesByName, kinesisResourcesByName, lambdaResourcesByName, sqsResourcesByName, restfulAPIResourcesByName,
		&id, &resources, &relationships)

	processTerraformResources(tfConfig.Resources,
		apiGatewayResourcesByName, cronResourcesByName, endpointResourcesByName, kinesisResourcesByName,
		sqsResourcesByName,
		relationshipsMap,
		endpointAPIGatewayMap, resourceEndpointAPIGatewayMap,
		&id, &resources)

	buildRelationships(relationshipsMap,
		apiGatewayResourcesByName, cronResourcesByName, endpointResourcesByName, kinesisResourcesByName,
		lambdaResourcesByName, sqsResourcesByName,
		resourceEndpointAPIGatewayMap,
		&relationships)

	return &drawio.ResourceCollection{Resources: resources, Relationships: relationships}
}

func buildRelationships(
	relationshipsMap map[resourceARN][]resourceARN,
	apiGatewayResourcesByName, cronResourcesByName, endpointResourcesByName, kinesisResourcesByName,
	lambdaResourcesByName, sqsResourcesByName map[string]drawio.Resource,
	resourceEndpointAPIGatewayMap map[resourceARN]map[resourceARN]resourceARN,
	relationships *[]drawio.Relationship,
) {
	for sourceARN, rel := range relationshipsMap {
		source := getResourceByARN(sourceARN,
			apiGatewayResourcesByName, cronResourcesByName, endpointResourcesByName, kinesisResourcesByName,
			lambdaResourcesByName, sqsResourcesByName)

		for i := range rel {
			targetARN := rel[i]

			target := getResourceByARN(targetARN,
				apiGatewayResourcesByName, cronResourcesByName, endpointResourcesByName, kinesisResourcesByName,
				lambdaResourcesByName, sqsResourcesByName)

			if endpointAPIGatewayMap, ok := resourceEndpointAPIGatewayMap[targetARN]; ok && targetARN.key != arnAPIGateway {
				updatedSource := getResourceByARN(endpointAPIGatewayMap[sourceARN],
					apiGatewayResourcesByName, cronResourcesByName, endpointResourcesByName, kinesisResourcesByName,
					lambdaResourcesByName, sqsResourcesByName)

				*relationships = append(*relationships, drawio.Relationship{Source: updatedSource, Target: target})

				continue
			}

			*relationships = append(*relationships, drawio.Relationship{Source: source, Target: target})
		}
	}
}

func processTerraformModules(
	tfModules []*terraform.Module,
	dbResourcesByName, googleBQResourcesByName, kinesisResourcesByName, lambdaResourcesByName, sqsResourcesByName,
	restfulAPIResourcesByName map[string]drawio.Resource,
	id *int, resources *[]drawio.Resource, relationships *[]drawio.Relationship,
) {
	for _, conf := range tfModules {
		if len(conf.Labels) == 1 {
			l := conf.Labels[0]

			if strings.HasSuffix(strings.ToLower(l), suffixLambda) {
				processLambdaModule(conf,
					dbResourcesByName, googleBQResourcesByName, kinesisResourcesByName, lambdaResourcesByName, sqsResourcesByName,
					restfulAPIResourcesByName,
					id, resources, relationships)
			}
		}
	}
}

func processTerraformResources(
	tfResources []*terraform.Resource,
	apiGatewayResourcesByName, cronResourcesByName, endpointResourcesByName, kinesisResourcesByName,
	sqsResourcesByName map[string]drawio.Resource,
	relationshipsMap map[resourceARN][]resourceARN,
	endpointAPIGatewayMap map[resourceARN]resourceARN,
	resourceEndpointAPIGatewayMap map[resourceARN]map[resourceARN]resourceARN,
	id *int, resources *[]drawio.Resource,
) {
	for _, conf := range tfResources {
		if len(conf.Labels) == 2 {
			switch conf.Labels[0] {
			case labelAWSAPIGatewayRoute:
				processAPIGatewayRoute(conf, apiGatewayResourcesByName, id, resources, relationshipsMap,
					endpointAPIGatewayMap)
			case labelAWSAPIGatewayIntegration:
				processAPIGatewayIntegration(conf,
					endpointAPIGatewayMap, resourceEndpointAPIGatewayMap, relationshipsMap)
			case labelAWSCloudwatchEventTarget:
				processCloudwatchEventTarget(conf, relationshipsMap)
			case labelAWSCron:
				processCronResource(conf, cronResourcesByName, id, resources)
			case labelAWSEndpoint:
				processEndpointResource(conf, endpointResourcesByName, id, resources)
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

func processAPIGatewayRoute(
	conf *terraform.Resource, resourcesByName map[string]drawio.Resource,
	id *int, resources *[]drawio.Resource, relationshipsMap map[resourceARN][]resourceARN,
	endpointAPIGatewayMap map[resourceARN]resourceARN,
) {
	value := conf.Attributes["route_key"].(string)

	resource := drawio.NewGenericResource(fmt.Sprintf("%d", *id), value, drawio.APIGatewayType)
	*id++

	*resources = append(*resources, resource)
	resourcesByName[value] = resource

	apiIDARN := resourceByARN(conf.Attributes["api_id"].(string))
	routeKeyARN := resourceARN{key: strings.Split(conf.Labels[1], "_")[1], name: conf.Attributes["route_key"].(string)}

	relationshipsMap[apiIDARN] = append(relationshipsMap[apiIDARN], routeKeyARN)

	endpointAPIGatewayMap[apiIDARN] = routeKeyARN
}

func processAPIGatewayIntegration(
	conf *terraform.Resource,
	endpointAPIGatewayMap map[resourceARN]resourceARN,
	resourceEndpointAPIGatewayMap map[resourceARN]map[resourceARN]resourceARN,
	relationshipsMap map[resourceARN][]resourceARN,
) {
	apiIDARN := resourceByARN(conf.Attributes["api_id"].(string))
	integrationURIARN := resourceByARN(conf.Attributes["integration_uri"].(string))

	relationshipsMap[apiIDARN] = append(relationshipsMap[apiIDARN], integrationURIARN)

	resourceEndpointAPIGatewayMap[integrationURIARN] = map[resourceARN]resourceARN{
		apiIDARN: endpointAPIGatewayMap[apiIDARN],
	}
}

func processCloudwatchEventTarget(conf *terraform.Resource, relationshipsMap map[resourceARN][]resourceARN) {
	ruleARN := resourceByARN(conf.Attributes["rule"].(string))
	arn := resourceByARN(conf.Attributes["arn"].(string))

	relationshipsMap[ruleARN] = append(relationshipsMap[ruleARN], arn)
}

func processCronResource(
	conf *terraform.Resource, resourcesByName map[string]drawio.Resource,
	id *int, resources *[]drawio.Resource,
) {
	value := conf.Attributes["schedule_expression"].(string)

	resource := drawio.NewGenericResource(fmt.Sprintf("%d", *id), value, drawio.CronType)
	*id++

	*resources = append(*resources, resource)
	resourcesByName[conf.Labels[1]] = resource
}

func processDBResourceFromEnvar(
	envar string, resourcesByName map[string]drawio.Resource, id *int, resources *[]drawio.Resource,
) *drawio.Resource {
	value := databaseName(envar, envarSuffixDBHost)

	resource, ok := resourcesByName[value]

	if !ok {
		resource = drawio.NewGenericResource(fmt.Sprintf("%d", *id), value, drawio.DatabaseType)
		*id++

		resourcesByName[value] = resource
		*resources = append(*resources, resource)
	}

	return &resource
}

func processEndpointResource(
	conf *terraform.Resource, resourcesByName map[string]drawio.Resource,
	id *int, resources *[]drawio.Resource,
) {
	value := conf.Attributes["domain_name"].(string)

	resource := drawio.NewGenericResource(fmt.Sprintf("%d", *id), value, drawio.EndpointType)
	*id++

	*resources = append(*resources, resource)
	resourcesByName[conf.Labels[1]] = resource
}

func processEventSourceMapping(conf *terraform.Resource, relationshipsMap map[resourceARN][]resourceARN) {
	eventSourceARN := resourceByARN(conf.Attributes["event_source_arn"].(string))
	functionName := resourceByARN(conf.Attributes["function_name"].(string))

	relationshipsMap[eventSourceARN] = append(relationshipsMap[eventSourceARN], functionName)
}

func processGoogleBQResourceFromEnvar(
	envar string, resourcesByName map[string]drawio.Resource, id *int, resources *[]drawio.Resource,
) *drawio.Resource {
	value := googleBQName(envar, envarSuffixGoogleBQ)

	resource, ok := resourcesByName[value]

	if !ok {
		resource = drawio.NewGenericResource(fmt.Sprintf("%d", *id), value, drawio.GoogleBQType)
		*id++

		resourcesByName[value] = resource
		*resources = append(*resources, resource)
	}

	return &resource
}

func processKinesisResourceFromEnvar(
	envar string, resourcesByName map[string]drawio.Resource,
	id *int, resources *[]drawio.Resource,
) *drawio.Resource {
	value := kinesisName(envar, envarSuffixKinesisStreamURL)

	resource, ok := resourcesByName[value]

	if !ok {
		resource = drawio.NewGenericResource(fmt.Sprintf("%d", *id), value, drawio.KinesisType)
		*id++

		resourcesByName[value] = resource
		*resources = append(*resources, resource)
	}

	return &resource
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

func processLambdaModule(conf *terraform.Module,
	dbResourcesByName, googleBQResourcesByName, kinesisResourcesByName, lambdaResourcesByName, sqsResourcesByName,
	restfulAPIResourcesByName map[string]drawio.Resource,
	id *int, resources *[]drawio.Resource, relationships *[]drawio.Relationship,
) {
	value := lambdaName(conf.Labels[0], suffixLambda)

	resource := drawio.NewGenericResource(fmt.Sprintf("%d", *id), value, drawio.LambdaType)
	*id++

	*resources = append(*resources, resource)
	lambdaResourcesByName[value] = resource

	for k := range conf.Attributes["lambda_function_env_vars"].(map[string]any) {
		if strings.HasSuffix(k, envarSuffixDBHost) {
			target := processDBResourceFromEnvar(k, dbResourcesByName, id, resources)
			*relationships = append(*relationships,
				drawio.Relationship{Source: resource, Target: *target})
		}

		if strings.HasSuffix(k, envarSuffixGoogleBQ) {
			target := processGoogleBQResourceFromEnvar(k, googleBQResourcesByName, id, resources)
			*relationships = append(*relationships,
				drawio.Relationship{Source: resource, Target: *target})
		}

		if strings.HasSuffix(k, envarSuffixKinesisStreamURL) {
			target := processKinesisResourceFromEnvar(k, kinesisResourcesByName, id, resources)
			*relationships = append(*relationships,
				drawio.Relationship{Source: resource, Target: *target})
		}

		if strings.HasSuffix(k, envarSuffixSQSQueueURL) {
			target := processSQSResourceFromEnvar(k, sqsResourcesByName, id, resources)
			*relationships = append(*relationships,
				drawio.Relationship{Source: resource, Target: *target})
		}

		if strings.HasSuffix(k, envarSuffixRestfulAPI) {
			target := processRestfulAPIResourceFromEnvar(k, restfulAPIResourcesByName, id, resources)
			*relationships = append(*relationships,
				drawio.Relationship{Source: resource, Target: *target})
		}
	}
}

func processSQSResourceFromEnvar(
	envar string, resourcesByName map[string]drawio.Resource, id *int, resources *[]drawio.Resource,
) *drawio.Resource {
	value := sqsName(envar, envarSuffixSQSQueueURL)

	resource, ok := resourcesByName[value]

	if !ok {
		resource = drawio.NewGenericResource(fmt.Sprintf("%d", *id), value, drawio.SQSType)
		*id++

		resourcesByName[value] = resource
		*resources = append(*resources, resource)
	}

	return &resource
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

func processRestfulAPIResourceFromEnvar(
	envar string, resourcesByName map[string]drawio.Resource, id *int, resources *[]drawio.Resource,
) *drawio.Resource {
	value := restfulAPIName(envar, envarSuffixRestfulAPI)

	resource, ok := resourcesByName[value]

	if !ok {
		resource = drawio.NewGenericResource(fmt.Sprintf("%d", *id), value, drawio.RestfulAPIType)
		*id++

		resourcesByName[value] = resource
		*resources = append(*resources, resource)
	}

	return &resource
}

///////

func getResourceByARN(
	arn resourceARN,
	apiGatewayResourcesByName, cronResourcesByName, endpointResourcesByName, kinesisResourcesByName,
	lambdaResourcesByName, sqsResourcesByName map[string]drawio.Resource,
) (resource drawio.Resource) {
	switch arn.key {
	case arnAPIGateway:
		resource = apiGatewayResourcesByName[arn.name]
	case arnCloudwatchKey:
		resource = cronResourcesByName[arn.name]
	case arnEndpoint:
		resource = endpointResourcesByName[arn.name]
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

func googleBQName(str, suffix string) string {
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

func restfulAPIName(str, suffix string) string {
	return strcase.ToCamel(str[:len(str)-len(suffix)])
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

		keyParts := strings.Split(parts[0], "_")

		key = keyParts[1]
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
