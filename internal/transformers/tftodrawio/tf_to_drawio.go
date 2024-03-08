package tftodrawio

import (
	"fmt"
	"strings"

	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
	"github.com/joselitofilho/aws-terraform-generator/internal/generators/terraform"
	"github.com/joselitofilho/aws-terraform-generator/internal/resources"
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

	resources     []resources.Resource
	relationships []resources.Relationship

	apiGatewayResourcesByName map[string]resources.Resource
	cronResourcesByName       map[string]resources.Resource
	dbResourcesByName         map[string]resources.Resource
	endpointResourcesByName   map[string]resources.Resource
	googleBQResourcesByName   map[string]resources.Resource
	kinesisResourcesByName    map[string]resources.Resource
	lambdaResourcesByName     map[string]resources.Resource
	restfulAPIResourcesByName map[string]resources.Resource
	s3BucketResourcesByName   map[string]resources.Resource
	sqsResourcesByName        map[string]resources.Resource

	endpointAPIGatewayMap   map[resourceARN][]resourceARN
	resourceAPIGIntegration map[resourceARN]resourceARN

	relationshipsMap map[resourceARN][]resourceARN

	id int
}

func NewTransformer(yamlConfig *config.Config, tfConfig *terraform.Config) *Transformer {
	return &Transformer{
		yamlConfig: yamlConfig,
		tfConfig:   tfConfig,

		resources:     []resources.Resource{},
		relationships: []resources.Relationship{},

		apiGatewayResourcesByName: map[string]resources.Resource{},
		cronResourcesByName:       map[string]resources.Resource{},
		dbResourcesByName:         map[string]resources.Resource{},
		endpointResourcesByName:   map[string]resources.Resource{},
		googleBQResourcesByName:   map[string]resources.Resource{},
		kinesisResourcesByName:    map[string]resources.Resource{},
		lambdaResourcesByName:     map[string]resources.Resource{},
		restfulAPIResourcesByName: map[string]resources.Resource{},
		s3BucketResourcesByName:   map[string]resources.Resource{},
		sqsResourcesByName:        map[string]resources.Resource{},

		endpointAPIGatewayMap:   map[resourceARN][]resourceARN{},
		resourceAPIGIntegration: map[resourceARN]resourceARN{},

		relationshipsMap: map[resourceARN][]resourceARN{},

		id: 1,
	}
}

func (t *Transformer) Transform() *resources.ResourceCollection {
	t.processTerraformModules()

	t.processTerraformResources()

	t.buildRelationships()

	t.applyFiltersInResources()
	t.applyFiltersInRelationships()

	return &resources.ResourceCollection{Resources: t.resources, Relationships: t.relationships}
}

func (t *Transformer) applyFiltersInResources() {
	filtered := make([]resources.Resource, 0, len(t.resources))

	for _, res := range t.resources {
		if t.hasResourceMatched(res, t.yamlConfig.Draw.Filters) {
			filtered = append(filtered, res)
		}
	}

	t.resources = filtered
}

func (t *Transformer) applyFiltersInRelationships() {
	filtered := make([]resources.Relationship, 0, len(t.relationships))

	for _, rel := range t.relationships {
		sourceMatch := t.hasResourceMatched(rel.Source, t.yamlConfig.Draw.Filters)
		targetMatch := t.hasResourceMatched(rel.Target, t.yamlConfig.Draw.Filters)

		if sourceMatch && targetMatch {
			filtered = append(filtered, rel)
		}
	}

	t.relationships = filtered
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

					t.relationships = append(t.relationships, resources.Relationship{Source: updatedSource, Target: target})
				}

				continue
			}

			t.relationships = append(t.relationships, resources.Relationship{Source: source, Target: target})
		}
	}
}

func (t *Transformer) getResourceByARN(arn resourceARN) (resource resources.Resource) {
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

func (t *Transformer) processAPIGatewayRoute(conf *terraform.Resource, resourcesByName map[string]resources.Resource) {
	value := conf.Attributes["route_key"].(string)

	resource := resources.NewGenericResource(fmt.Sprintf("%d", t.id), value, resources.APIGatewayType)
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
	conf *terraform.Resource, resourcesByName map[string]resources.Resource,
) {
	var value string

	switch v := conf.Attributes["schedule_expression"].(type) {
	case string:
		value = v
	default:
		value = "schedule_expression"
	}

	resource := resources.NewGenericResource(fmt.Sprintf("%d", t.id), value, resources.CronType)
	t.id++

	t.resources = append(t.resources, resource)
	resourcesByName[conf.Labels[1]] = resource
}

func (t *Transformer) processDBResourceFromEnvar(
	k, v string, resourcesByName map[string]resources.Resource,
) resources.Resource {
	value := toKebabFromEnvar(k, v, envarSuffixDBHost)

	resource, ok := resourcesByName[value]

	if !ok {
		resource = resources.NewGenericResource(fmt.Sprintf("%d", t.id), value, resources.DatabaseType)
		t.id++

		resourcesByName[value] = resource
		t.resources = append(t.resources, resource)
	}

	return resource
}

func (t *Transformer) processEndpointResource(
	conf *terraform.Resource, resourcesByName map[string]resources.Resource,
) {
	value := replaceVars(conf.Attributes["domain_name"].(string), t.tfConfig.Locals)

	resource := resources.NewGenericResource(fmt.Sprintf("%d", t.id), value, resources.EndpointType)
	t.id++

	t.resources = append(t.resources, resource)
	resourcesByName[conf.Labels[1]] = resource
}

func (t *Transformer) processEventSourceMapping(conf *terraform.Resource) {
	eventSourceARN := resourceByARN(conf.Attributes["event_source_arn"].(string))
	functionName := resourceByARN(conf.Attributes["function_name"].(string))

	t.relationshipsMap[eventSourceARN] = append(t.relationshipsMap[eventSourceARN], functionName)

	t.tryToCreateKinesisResourceByARN(eventSourceARN)
}

func (t *Transformer) processGoogleBQResourceFromEnvar(
	k, v string, resourcesByName map[string]resources.Resource,
) resources.Resource {
	value := replaceVars(v, t.tfConfig.Locals)
	value = toKebabFromEnvar(k, value, envarSuffixGoogleBQ)

	resource, ok := resourcesByName[value]

	if !ok {
		resource = resources.NewGenericResource(fmt.Sprintf("%d", t.id), value, resources.GoogleBQType)
		t.id++

		resourcesByName[value] = resource
		t.resources = append(t.resources, resource)
	}

	return resource
}

func (t *Transformer) processKinesisResourceFromEnvar(
	k, v string, resourcesByName map[string]resources.Resource,
) resources.Resource {
	value := toPascalFromEnvar(k, v, envarSuffixKinesisStreamURL)

	resource, ok := resourcesByName[value]

	if !ok {
		resource = resources.NewGenericResource(fmt.Sprintf("%d", t.id), value, resources.KinesisType)
		t.id++

		resourcesByName[value] = resource
		t.resources = append(t.resources, resource)
	}

	return resource
}

func (t *Transformer) processKinesisResource(
	conf *terraform.Resource, resourcesByName map[string]resources.Resource,
) {
	l := conf.Labels[1]

	if strings.HasSuffix(strings.ToLower(l), suffixKinesis) {
		value := toPascalFromEnvar(l, l, suffixKinesis)

		resource := resources.NewGenericResource(fmt.Sprintf("%d", t.id), value, resources.KinesisType)
		t.id++

		t.resources = append(t.resources, resource)
		resourcesByName[value] = resource
	}
}

func (t *Transformer) processLambdaModule(conf *terraform.Module) {
	value := lambdaName(conf.Labels[0], suffixLambda)

	resource := resources.NewGenericResource(fmt.Sprintf("%d", t.id), value, resources.LambdaType)
	t.id++

	t.resources = append(t.resources, resource)
	t.lambdaResourcesByName[value] = resource

	for k, v := range conf.Attributes["lambda_function_env_vars"].(map[string]any) {
		switch {
		case strings.HasSuffix(k, envarSuffixDBHost):
			target := t.processDBResourceFromEnvar(k, v.(string), t.dbResourcesByName)
			t.relationships = append(t.relationships,
				resources.Relationship{Source: resource, Target: target})
		case strings.HasSuffix(k, envarSuffixGoogleBQ):
			target := t.processGoogleBQResourceFromEnvar(k, v.(string), t.googleBQResourcesByName)
			t.relationships = append(t.relationships,
				resources.Relationship{Source: resource, Target: target})
		case strings.HasSuffix(k, envarSuffixKinesisStreamURL):
			target := t.processKinesisResourceFromEnvar(k, v.(string), t.kinesisResourcesByName)
			t.relationships = append(t.relationships,
				resources.Relationship{Source: resource, Target: target})
		case strings.HasSuffix(k, envarSuffixS3BucketURL):
			target := t.processS3BucketResourceFromEnvar(k, v.(string), t.s3BucketResourcesByName)
			t.relationships = append(t.relationships,
				resources.Relationship{Source: resource, Target: target})
		case strings.HasSuffix(k, envarSuffixSQSQueueURL):
			target := t.processSQSResourceFromEnvar(k, v.(string), t.sqsResourcesByName)
			t.relationships = append(t.relationships,
				resources.Relationship{Source: resource, Target: target})
		case strings.HasSuffix(k, envarSuffixRestfulAPI):
			target := t.processRestfulAPIResourceFromEnvar(k, v.(string), t.restfulAPIResourcesByName)
			t.relationships = append(t.relationships,
				resources.Relationship{Source: resource, Target: target})
		}
	}
}

func (t *Transformer) processS3BucketResourceFromEnvar(
	k, v string, resourcesByName map[string]resources.Resource,
) resources.Resource {
	value := toKebabFromEnvar(k, v, envarSuffixS3BucketURL)

	resource, ok := resourcesByName[value]

	if !ok {
		resource = resources.NewGenericResource(fmt.Sprintf("%d", t.id), value, resources.S3Type)
		t.id++

		resourcesByName[value] = resource
		t.resources = append(t.resources, resource)
	}

	return resource
}

func (t *Transformer) processS3BucketResource(
	conf *terraform.Resource, s3BucketResourcesByName map[string]resources.Resource,
) {
	l := conf.Labels[1]

	if strings.HasSuffix(strings.ToLower(l), suffixS3Bucket) {
		value := s3BucketName(l, suffixS3Bucket)

		resource := resources.NewGenericResource(fmt.Sprintf("%d", t.id), value, resources.S3Type)
		t.id++

		t.resources = append(t.resources, resource)
		s3BucketResourcesByName[value] = resource
	}
}

func (t *Transformer) processSQSResourceFromEnvar(
	k, v string, resourcesByName map[string]resources.Resource,
) resources.Resource {
	value := toKebabFromEnvar(k, v, envarSuffixSQSQueueURL)

	resource, ok := resourcesByName[value]

	if !ok {
		resource = resources.NewGenericResource(fmt.Sprintf("%d", t.id), value, resources.SQSType)
		t.id++

		resourcesByName[value] = resource
		t.resources = append(t.resources, resource)
	}

	return resource
}

func (t *Transformer) processSQSResource(conf *terraform.Resource, sqsResourcesByName map[string]resources.Resource) {
	l := conf.Labels[1]

	if strings.HasSuffix(strings.ToLower(l), suffixSQS) {
		value := sqsName(l, suffixSQS)

		resource := resources.NewGenericResource(fmt.Sprintf("%d", t.id), value, resources.SQSType)
		t.id++

		t.resources = append(t.resources, resource)
		sqsResourcesByName[value] = resource
	}
}

func (t *Transformer) processRestfulAPIResourceFromEnvar(
	k, v string, resourcesByName map[string]resources.Resource,
) resources.Resource {
	value := toCamelFromEnvar(k, v, envarSuffixRestfulAPI)

	resource, ok := resourcesByName[value]

	if !ok {
		resource = resources.NewGenericResource(fmt.Sprintf("%d", t.id), value, resources.RestfulAPIType)
		t.id++

		resourcesByName[value] = resource
		t.resources = append(t.resources, resource)
	}

	return resource
}

func (t *Transformer) tryToCreateKinesisResourceByARN(eventSourceARN resourceARN) {
	if eventSourceARN.key == arnKinesisKey {
		value := toPascalFromEnvar(eventSourceARN.name, eventSourceARN.name, envarSuffixKinesisStreamURL)

		if _, ok := t.kinesisResourcesByName[value]; !ok {
			resource := resources.NewGenericResource(fmt.Sprintf("%d", t.id), value, resources.KinesisType)
			t.id++

			t.resources = append(t.resources, resource)
			t.kinesisResourcesByName[value] = resource
		}
	}
}
