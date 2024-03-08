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
	envarSuffixGoogleBQ         = "BQ_PROJECT_ID"
	envarSuffixKinesisStreamURL = "KINESIS_STREAM_URL"
	envarSuffixS3BucketURL      = "S3_BUCKET"
	envarSuffixSQSQueueURL      = "SQS_QUEUE_URL"
	envarSuffixRestfulAPI       = "API_BASE_URL"
)

var (
	suffixKinesis  = "_kinesis"
	suffixLambda   = "_lambda"
	suffixS3Bucket = "_bucket"
	suffixSQS      = "_sqs"
)

var (
	labelAWSAPIGatewayRoute          = "aws_apigatewayv2_route"
	labelAWSAPIGatewayIntegration    = "aws_apigatewayv2_integration"
	labelAWSCloudwatchEventTarget    = "aws_cloudwatch_event_target"
	labelAWSCron                     = "aws_cloudwatch_event_rule"
	labelAWSEndpoint                 = "aws_apigatewayv2_domain_name"
	labelAWSKinesisStream            = "aws_kinesis_stream"
	labelAWSLambdaFunction           = "aws_lambda_function"
	labelAWSLambdaEventSourceMapping = "aws_lambda_event_source_mapping"
	labelAWSS3Bucket                 = "aws_s3_bucket"
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
	s3BucketResourcesByName   map[string]drawio.Resource
	sqsResourcesByName        map[string]drawio.Resource

	endpointAPIGatewayMap   map[resourceARN][]resourceARN
	resourceAPIGIntegration map[resourceARN]resourceARN

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
		s3BucketResourcesByName:   map[string]drawio.Resource{},
		sqsResourcesByName:        map[string]drawio.Resource{},

		endpointAPIGatewayMap:   map[resourceARN][]resourceARN{},
		resourceAPIGIntegration: map[resourceARN]resourceARN{},

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

			if integration, ok := t.resourceAPIGIntegration[targetARN]; ok {
				for _, apig := range t.endpointAPIGatewayMap[integration] {
					updatedSource := t.getResourceByARN(apig)

					t.relationships = append(t.relationships, drawio.Relationship{Source: updatedSource, Target: target})
				}

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
				t.processAPIGatewayIntegration(tfResourceConf)
			case labelAWSCloudwatchEventTarget:
				t.processCloudwatchEventTarget(tfResourceConf)
			case labelAWSCron:
				t.processCronResource(tfResourceConf, t.cronResourcesByName)
			case labelAWSEndpoint:
				t.processEndpointResource(tfResourceConf, t.endpointResourcesByName)
			case labelAWSKinesisStream:
				t.processKinesisResource(tfResourceConf, t.kinesisResourcesByName)
			case labelAWSS3Bucket:
				t.processS3BucketResource(tfResourceConf, t.s3BucketResourcesByName)
			case labelAWSSQSQueue:
				t.processSQSResource(tfResourceConf, t.sqsResourcesByName)
			case labelAWSLambdaEventSourceMapping:
				t.processEventSourceMapping(tfResourceConf)
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
	routeKeyARN := resourceARN{key: strings.Split(conf.Labels[0], "_")[2], name: conf.Attributes["route_key"].(string)}
	target := resourceARN{key: strings.Split(conf.Attributes["target"].(string), ".")[1]}

	t.relationshipsMap[apiIDARN] = append(t.relationshipsMap[apiIDARN], routeKeyARN)

	t.endpointAPIGatewayMap[target] = append(t.endpointAPIGatewayMap[target], routeKeyARN)
}

func (t *Transformer) processAPIGatewayIntegration(conf *terraform.Resource) {
	label := resourceARN{key: conf.Labels[1]}
	apiIDARN := resourceByARN(conf.Attributes["api_id"].(string))
	integrationURIARN := resourceByARN(conf.Attributes["integration_uri"].(string))

	t.relationshipsMap[apiIDARN] = append(t.relationshipsMap[apiIDARN], integrationURIARN)

	t.resourceAPIGIntegration[integrationURIARN] = label
}

func (t *Transformer) processCloudwatchEventTarget(conf *terraform.Resource) {
	ruleARN := resourceByARN(conf.Attributes["rule"].(string))
	arn := resourceByARN(conf.Attributes["arn"].(string))

	t.relationshipsMap[ruleARN] = append(t.relationshipsMap[ruleARN], arn)
}

func (t *Transformer) processCronResource(
	conf *terraform.Resource, resourcesByName map[string]drawio.Resource,
) {
	var value string

	switch v := conf.Attributes["schedule_expression"].(type) {
	case string:
		value = v
	default:
		value = "schedule_expression"
	}

	resource := drawio.NewGenericResource(fmt.Sprintf("%d", t.id), value, drawio.CronType)
	t.id++

	t.resources = append(t.resources, resource)
	resourcesByName[conf.Labels[1]] = resource
}

func (t *Transformer) processDBResourceFromEnvar(
	k, v string, resourcesByName map[string]drawio.Resource,
) *drawio.Resource {
	value := toKebabFromEnvar(k, v, envarSuffixDBHost)

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

func (t *Transformer) processEventSourceMapping(conf *terraform.Resource) {
	eventSourceARN := resourceByARN(conf.Attributes["event_source_arn"].(string))
	functionName := resourceByARN(conf.Attributes["function_name"].(string))

	t.relationshipsMap[eventSourceARN] = append(t.relationshipsMap[eventSourceARN], functionName)

	t.tryToCreateResourceByARN(eventSourceARN)
}

func (t *Transformer) tryToCreateResourceByARN(eventSourceARN resourceARN) {
	switch eventSourceARN.key {
	case arnKinesisKey:
		value := toPascalFromEnvar(eventSourceARN.name, eventSourceARN.name, envarSuffixKinesisStreamURL)

		if _, ok := t.kinesisResourcesByName[value]; !ok {
			resource := drawio.NewGenericResource(fmt.Sprintf("%d", t.id), value, drawio.KinesisType)
			t.id++

			t.resources = append(t.resources, resource)
			t.kinesisResourcesByName[value] = resource
		}
	}
}

func (t *Transformer) processGoogleBQResourceFromEnvar(
	k, v string, resourcesByName map[string]drawio.Resource,
) *drawio.Resource {
	value := replaceVars(v, t.tfConfig.Locals)
	value = toKebabFromEnvar(k, value, envarSuffixGoogleBQ)

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
	k, v string, resourcesByName map[string]drawio.Resource,
) *drawio.Resource {
	value := toPascalFromEnvar(k, v, envarSuffixKinesisStreamURL)

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
	conf *terraform.Resource, resourcesByName map[string]drawio.Resource,
) {
	l := conf.Labels[1]

	if strings.HasSuffix(strings.ToLower(l), suffixKinesis) {
		value := toPascalFromEnvar(l, l, suffixKinesis)

		resource := drawio.NewGenericResource(fmt.Sprintf("%d", t.id), value, drawio.KinesisType)
		t.id++

		t.resources = append(t.resources, resource)
		resourcesByName[value] = resource
	}
}

func (t *Transformer) processLambdaModule(conf *terraform.Module) {
	value := lambdaName(conf.Labels[0], suffixLambda)

	resource := drawio.NewGenericResource(fmt.Sprintf("%d", t.id), value, drawio.LambdaType)
	t.id++

	t.resources = append(t.resources, resource)
	t.lambdaResourcesByName[value] = resource

	for k, v := range conf.Attributes["lambda_function_env_vars"].(map[string]any) {
		switch {
		case strings.HasSuffix(k, envarSuffixDBHost):
			target := t.processDBResourceFromEnvar(k, v.(string), t.dbResourcesByName)
			t.relationships = append(t.relationships,
				drawio.Relationship{Source: resource, Target: *target})
		case strings.HasSuffix(k, envarSuffixGoogleBQ):
			target := t.processGoogleBQResourceFromEnvar(k, v.(string), t.googleBQResourcesByName)
			t.relationships = append(t.relationships,
				drawio.Relationship{Source: resource, Target: *target})
		case strings.HasSuffix(k, envarSuffixKinesisStreamURL):
			target := t.processKinesisResourceFromEnvar(k, v.(string), t.kinesisResourcesByName)
			t.relationships = append(t.relationships,
				drawio.Relationship{Source: resource, Target: *target})
		case strings.HasSuffix(k, envarSuffixS3BucketURL):
			target := t.processS3BucketResourceFromEnvar(k, v.(string), t.s3BucketResourcesByName)
			t.relationships = append(t.relationships,
				drawio.Relationship{Source: resource, Target: *target})
		case strings.HasSuffix(k, envarSuffixSQSQueueURL):
			target := t.processSQSResourceFromEnvar(k, v.(string), t.sqsResourcesByName)
			t.relationships = append(t.relationships,
				drawio.Relationship{Source: resource, Target: *target})
		case strings.HasSuffix(k, envarSuffixRestfulAPI):
			target := t.processRestfulAPIResourceFromEnvar(k, v.(string), t.restfulAPIResourcesByName)
			t.relationships = append(t.relationships,
				drawio.Relationship{Source: resource, Target: *target})
		}
	}
}

func (t *Transformer) processS3BucketResourceFromEnvar(
	k, v string, resourcesByName map[string]drawio.Resource,
) *drawio.Resource {
	value := toKebabFromEnvar(k, v, envarSuffixS3BucketURL)

	resource, ok := resourcesByName[value]

	if !ok {
		resource = drawio.NewGenericResource(fmt.Sprintf("%d", t.id), value, drawio.S3Type)
		t.id++

		resourcesByName[value] = resource
		t.resources = append(t.resources, resource)
	}

	return &resource
}

func (t *Transformer) processS3BucketResource(
	conf *terraform.Resource, s3BucketResourcesByName map[string]drawio.Resource,
) {
	l := conf.Labels[1]

	if strings.HasSuffix(strings.ToLower(l), suffixS3Bucket) {
		value := s3BucketName(l, suffixS3Bucket)

		resource := drawio.NewGenericResource(fmt.Sprintf("%d", t.id), value, drawio.S3Type)
		t.id++

		t.resources = append(t.resources, resource)
		s3BucketResourcesByName[value] = resource
	}
}

func (t *Transformer) processSQSResourceFromEnvar(
	k, v string, resourcesByName map[string]drawio.Resource,
) *drawio.Resource {
	value := toKebabFromEnvar(k, v, envarSuffixSQSQueueURL)

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
	k, v string, resourcesByName map[string]drawio.Resource,
) *drawio.Resource {
	value := toCamelFromEnvar(k, v, envarSuffixRestfulAPI)

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
			switch v := v.(type) {
			// TODO: Implement others
			case string:
				result = strings.ReplaceAll(result, k, v)
			case []string:
				result = fmt.Sprintf("%slocal.%s[]", result, k)
			default:
				result = fmt.Sprintf("%slocal.%s", result, k)
			}
		}
	}

	return result
}

////////////////////////////////////////////////////////////////////////////////

func strTransformFromEnvar(
	key, value, suffix string, f func(s string) string,
) string {
	var result string

	if key == suffix {
		suffixMap := map[string]struct{}{
			labelAWSKinesisStream:  {},
			labelAWSLambdaFunction: {},
			labelAWSS3Bucket:       {},
			labelAWSSQSQueue:       {},
		}

		result = value

		if strings.HasPrefix(result, "var.client-var.environment-") { // TODO:
			result = f(strings.ReplaceAll(result, "var.client-var.environment-", "")) // TODO:
		}

		for s := range suffixMap {
			if strings.HasPrefix(result, s) {
				result = f(resourceByARN(result).name)
				break
			}
		}
	} else {
		result = key

		result = strings.ReplaceAll(result, "_"+suffix, "")
		result = strings.ReplaceAll(result, suffix, "")
		result = f(result)
	}

	return result
}

func toCamelFromEnvar(key, value, suffix string) string {
	return strTransformFromEnvar(key, value, suffix, strcase.ToCamel)
}

func toKebabFromEnvar(key, value, suffix string) string {
	return strTransformFromEnvar(key, value, suffix, strcase.ToKebab)
}

func toPascalFromEnvar(key, value, suffix string) string {
	return strTransformFromEnvar(key, value, suffix, strcase.ToPascal)
}

func lambdaName(str, suffix string) string {
	return strcase.ToCamel(str[:len(str)-len(suffix)])
}

func s3BucketName(str, suffix string) string {
	return strcase.ToKebab(str[:len(str)-len(suffix)])
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

		keyParts := strings.Split(parts[0], "_")

		if parts[0] == "module" {
			// TODO: Add support to more type of modules
			key = arnLambdaKey
		} else if len(keyParts) > 1 {
			key = keyParts[1]
		} else {
			key = strings.Join(keyParts, "_")
		}

		name = parts[1]

		switch key {
		case arnKinesisKey:
			name = toPascalFromEnvar(name, name, suffixKinesis)
		case arnLambdaKey:
			name = toCamelFromEnvar(name, name, suffixLambda)
		case arnSQSKey:
			name = toKebabFromEnvar(name, name, suffixSQS)
		}
	}

	return resourceARN{key: key, name: name}
}
