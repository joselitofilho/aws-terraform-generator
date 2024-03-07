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

type Transformer struct {
	yamlConfig *config.Config
	tfConfig   *terraform.Config

	resources     []drawio.Resource
	relationships []drawio.Relationship

	apiGatewayResourcesByName map[string]drawio.Resource
	cronResourcesByName       map[string]drawio.Resource
	dbResourcesByName         map[string]drawio.Resource
	endpointResourcesByName   map[string]drawio.Resource
	googleBQResourcesByName   map[string]drawio.Resource
	kinesisResourcesByName    map[string]drawio.Resource
	lambdaResourcesByName     map[string]drawio.Resource
	restfulAPIResourcesByName map[string]drawio.Resource
	sqsResourcesByName        map[string]drawio.Resource

	endpointAPIGatewayMap         map[resourceARN]resourceARN
	resourceEndpointAPIGatewayMap map[resourceARN]map[resourceARN]resourceARN

	relationshipsMap map[resourceARN][]resourceARN

	id int
}

func NewTransformer(yamlConfig *config.Config, tfConfig *terraform.Config) *Transformer {
	return &Transformer{
		yamlConfig: yamlConfig,
		tfConfig:   tfConfig,

		resources:     []drawio.Resource{},
		relationships: []drawio.Relationship{},

		apiGatewayResourcesByName: map[string]drawio.Resource{},
		cronResourcesByName:       map[string]drawio.Resource{},
		dbResourcesByName:         map[string]drawio.Resource{},
		endpointResourcesByName:   map[string]drawio.Resource{},
		googleBQResourcesByName:   map[string]drawio.Resource{},
		kinesisResourcesByName:    map[string]drawio.Resource{},
		lambdaResourcesByName:     map[string]drawio.Resource{},
		restfulAPIResourcesByName: map[string]drawio.Resource{},
		sqsResourcesByName:        map[string]drawio.Resource{},

		endpointAPIGatewayMap:         map[resourceARN]resourceARN{},
		resourceEndpointAPIGatewayMap: map[resourceARN]map[resourceARN]resourceARN{},

		relationshipsMap: map[resourceARN][]resourceARN{},

		id: 1,
	}
}

func (t *Transformer) Transform() *drawio.ResourceCollection {
	t.processTerraformModules()

	t.processTerraformResources()

	t.buildRelationships()

	return &drawio.ResourceCollection{Resources: t.resources, Relationships: t.relationships}
}

func (t *Transformer) buildRelationships() {
	for sourceARN, rel := range t.relationshipsMap {
		source := t.getResourceByARN(sourceARN)

		for i := range rel {
			targetARN := rel[i]

			target := t.getResourceByARN(targetARN)

			if endpointAPIGatewayMap, ok := t.resourceEndpointAPIGatewayMap[targetARN]; ok &&
				targetARN.key != arnAPIGateway {
				updatedSource := t.getResourceByARN(endpointAPIGatewayMap[sourceARN])

				t.relationships = append(t.relationships, drawio.Relationship{Source: updatedSource, Target: target})

				continue
			}

			t.relationships = append(t.relationships, drawio.Relationship{Source: source, Target: target})
		}
	}
}

func (t *Transformer) processTerraformModules() {
	for _, tfModule := range t.tfConfig.Modules {
		if len(tfModule.Labels) == 1 {
			l := tfModule.Labels[0]

			if strings.HasSuffix(strings.ToLower(l), suffixLambda) {
				t.processLambdaModule(tfModule)
			}
		}
	}
}

func (t *Transformer) processTerraformResources() {
	for _, tfResourceConf := range t.tfConfig.Resources {
		if len(tfResourceConf.Labels) == 2 {
			switch tfResourceConf.Labels[0] {
			case labelAWSAPIGatewayRoute:
				t.processAPIGatewayRoute(tfResourceConf, t.apiGatewayResourcesByName)
			case labelAWSAPIGatewayIntegration:
				t.processAPIGatewayIntegration(tfResourceConf, t.endpointAPIGatewayMap)
			case labelAWSCloudwatchEventTarget:
				t.processCloudwatchEventTarget(tfResourceConf, t.relationshipsMap)
			case labelAWSCron:
				t.processCronResource(tfResourceConf, t.cronResourcesByName)
			case labelAWSEndpoint:
				t.processEndpointResource(tfResourceConf, t.endpointResourcesByName)
			case labelAWSKinesisStream:
				t.processKinesisResource(tfResourceConf, t.kinesisResourcesByName)
			case labelAWSSQSQueue:
				t.processSQSResource(tfResourceConf, t.sqsResourcesByName)
			case labelAWSLambdaEventSourceMapping:
				t.processEventSourceMapping(tfResourceConf, t.relationshipsMap)
			}
		}
	}
}

func (t *Transformer) processAPIGatewayRoute(conf *terraform.Resource, resourcesByName map[string]drawio.Resource) {
	value := conf.Attributes["route_key"].(string)

	resource := drawio.NewGenericResource(fmt.Sprintf("%d", t.id), value, drawio.APIGatewayType)
	t.id++

	t.resources = append(t.resources, resource)
	resourcesByName[value] = resource

	apiIDARN := resourceByARN(conf.Attributes["api_id"].(string))
	routeKeyARN := resourceARN{key: strings.Split(conf.Labels[1], "_")[1], name: conf.Attributes["route_key"].(string)}

	t.relationshipsMap[apiIDARN] = append(t.relationshipsMap[apiIDARN], routeKeyARN)

	t.endpointAPIGatewayMap[apiIDARN] = routeKeyARN
}

func (t *Transformer) processAPIGatewayIntegration(
	conf *terraform.Resource, endpointAPIGatewayMap map[resourceARN]resourceARN,
) {
	apiIDARN := resourceByARN(conf.Attributes["api_id"].(string))
	integrationURIARN := resourceByARN(conf.Attributes["integration_uri"].(string))

	t.relationshipsMap[apiIDARN] = append(t.relationshipsMap[apiIDARN], integrationURIARN)

	t.resourceEndpointAPIGatewayMap[integrationURIARN] = map[resourceARN]resourceARN{
		apiIDARN: endpointAPIGatewayMap[apiIDARN],
	}
}

func (t *Transformer) processCloudwatchEventTarget(conf *terraform.Resource, relationshipsMap map[resourceARN][]resourceARN) {
	ruleARN := resourceByARN(conf.Attributes["rule"].(string))
	arn := resourceByARN(conf.Attributes["arn"].(string))

	relationshipsMap[ruleARN] = append(relationshipsMap[ruleARN], arn)
}

func (t *Transformer) processCronResource(
	conf *terraform.Resource, resourcesByName map[string]drawio.Resource,
) {
	value := conf.Attributes["schedule_expression"].(string)

	resource := drawio.NewGenericResource(fmt.Sprintf("%d", t.id), value, drawio.CronType)
	t.id++

	t.resources = append(t.resources, resource)
	resourcesByName[conf.Labels[1]] = resource
}

func (t *Transformer) processDBResourceFromEnvar(
	envar string, resourcesByName map[string]drawio.Resource,
) *drawio.Resource {
	value := databaseName(envar, envarSuffixDBHost)

	resource, ok := resourcesByName[value]

	if !ok {
		resource = drawio.NewGenericResource(fmt.Sprintf("%d", t.id), value, drawio.DatabaseType)
		t.id++

		resourcesByName[value] = resource
		t.resources = append(t.resources, resource)
	}

	return &resource
}

func (t *Transformer) processEndpointResource(
	conf *terraform.Resource, resourcesByName map[string]drawio.Resource,
) {
	value := replaceVars(conf.Attributes["domain_name"].(string), t.tfConfig.Locals)

	resource := drawio.NewGenericResource(fmt.Sprintf("%d", t.id), value, drawio.EndpointType)
	t.id++

	t.resources = append(t.resources, resource)
	resourcesByName[conf.Labels[1]] = resource
}

func (t *Transformer) processEventSourceMapping(conf *terraform.Resource, relationshipsMap map[resourceARN][]resourceARN) {
	eventSourceARN := resourceByARN(conf.Attributes["event_source_arn"].(string))
	functionName := resourceByARN(conf.Attributes["function_name"].(string))

	relationshipsMap[eventSourceARN] = append(relationshipsMap[eventSourceARN], functionName)
}

func (t *Transformer) processGoogleBQResourceFromEnvar(
	envar string, resourcesByName map[string]drawio.Resource,
) *drawio.Resource {
	value := googleBQName(envar, envarSuffixGoogleBQ)

	resource, ok := resourcesByName[value]

	if !ok {
		resource = drawio.NewGenericResource(fmt.Sprintf("%d", t.id), value, drawio.GoogleBQType)
		t.id++

		resourcesByName[value] = resource
		t.resources = append(t.resources, resource)
	}

	return &resource
}

func (t *Transformer) processKinesisResourceFromEnvar(
	envar string, resourcesByName map[string]drawio.Resource,
) *drawio.Resource {
	value := kinesisName(envar, envarSuffixKinesisStreamURL)

	resource, ok := resourcesByName[value]

	if !ok {
		resource = drawio.NewGenericResource(fmt.Sprintf("%d", t.id), value, drawio.KinesisType)
		t.id++

		resourcesByName[value] = resource
		t.resources = append(t.resources, resource)
	}

	return &resource
}

func (t *Transformer) processKinesisResource(
	conf *terraform.Resource, kinesisResourcesByName map[string]drawio.Resource,
) {
	l := conf.Labels[1]

	if strings.HasSuffix(strings.ToLower(l), suffixKinesis) {
		value := kinesisName(l, suffixKinesis)

		resource := drawio.NewGenericResource(fmt.Sprintf("%d", t.id), value, drawio.KinesisType)
		t.id++

		t.resources = append(t.resources, resource)
		kinesisResourcesByName[value] = resource
	}
}

func (t *Transformer) processLambdaModule(conf *terraform.Module) {
	value := lambdaName(conf.Labels[0], suffixLambda)

	resource := drawio.NewGenericResource(fmt.Sprintf("%d", t.id), value, drawio.LambdaType)
	t.id++

	t.resources = append(t.resources, resource)
	t.lambdaResourcesByName[value] = resource

	for k := range conf.Attributes["lambda_function_env_vars"].(map[string]any) {
		if strings.HasSuffix(k, envarSuffixDBHost) {
			target := t.processDBResourceFromEnvar(k, t.dbResourcesByName)
			t.relationships = append(t.relationships,
				drawio.Relationship{Source: resource, Target: *target})
		}

		if strings.HasSuffix(k, envarSuffixGoogleBQ) {
			target := t.processGoogleBQResourceFromEnvar(k, t.googleBQResourcesByName)
			t.relationships = append(t.relationships,
				drawio.Relationship{Source: resource, Target: *target})
		}

		if strings.HasSuffix(k, envarSuffixKinesisStreamURL) {
			target := t.processKinesisResourceFromEnvar(k, t.kinesisResourcesByName)
			t.relationships = append(t.relationships,
				drawio.Relationship{Source: resource, Target: *target})
		}

		if strings.HasSuffix(k, envarSuffixSQSQueueURL) {
			target := t.processSQSResourceFromEnvar(k, t.sqsResourcesByName)
			t.relationships = append(t.relationships,
				drawio.Relationship{Source: resource, Target: *target})
		}

		if strings.HasSuffix(k, envarSuffixRestfulAPI) {
			target := t.processRestfulAPIResourceFromEnvar(k, t.restfulAPIResourcesByName)
			t.relationships = append(t.relationships,
				drawio.Relationship{Source: resource, Target: *target})
		}
	}
}

func (t *Transformer) processSQSResourceFromEnvar(
	envar string, resourcesByName map[string]drawio.Resource,
) *drawio.Resource {
	value := sqsName(envar, envarSuffixSQSQueueURL)

	resource, ok := resourcesByName[value]

	if !ok {
		resource = drawio.NewGenericResource(fmt.Sprintf("%d", t.id), value, drawio.SQSType)
		t.id++

		resourcesByName[value] = resource
		t.resources = append(t.resources, resource)
	}

	return &resource
}

func (t *Transformer) processSQSResource(conf *terraform.Resource, sqsResourcesByName map[string]drawio.Resource) {
	l := conf.Labels[1]

	if strings.HasSuffix(strings.ToLower(l), suffixSQS) {
		value := sqsName(l, suffixSQS)

		resource := drawio.NewGenericResource(fmt.Sprintf("%d", t.id), value, drawio.SQSType)
		t.id++

		t.resources = append(t.resources, resource)
		sqsResourcesByName[value] = resource
	}
}

func (t *Transformer) processRestfulAPIResourceFromEnvar(
	envar string, resourcesByName map[string]drawio.Resource,
) *drawio.Resource {
	value := restfulAPIName(envar, envarSuffixRestfulAPI)

	resource, ok := resourcesByName[value]

	if !ok {
		resource = drawio.NewGenericResource(fmt.Sprintf("%d", t.id), value, drawio.RestfulAPIType)
		t.id++

		resourcesByName[value] = resource
		t.resources = append(t.resources, resource)
	}

	return &resource
}

func (t *Transformer) getResourceByARN(arn resourceARN) (resource drawio.Resource) {
	switch arn.key {
	case arnAPIGateway:
		resource = t.apiGatewayResourcesByName[arn.name]
	case arnCloudwatchKey:
		resource = t.cronResourcesByName[arn.name]
	case arnEndpoint:
		resource = t.endpointResourcesByName[arn.name]
	case arnKinesisKey:
		resource = t.kinesisResourcesByName[arn.name]
	case arnLambdaKey:
		resource = t.lambdaResourcesByName[arn.name]
	case arnSQSKey:
		resource = t.sqsResourcesByName[arn.name]
	}

	return resource
}

////////////////////////////////////////////////////////////////////////////////

func replaceVars(str string, tfLocals []*terraform.Local) string {
	result := str

	for i := range tfLocals {
		for k, v := range tfLocals[i].Attributes {
			result = strings.ReplaceAll(result, k, v.(string))
		}
	}

	return result
}

////////////////////////////////////////////////////////////////////////////////

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
