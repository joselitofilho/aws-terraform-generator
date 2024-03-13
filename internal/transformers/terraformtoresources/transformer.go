package terraformtoresources

import (
	"fmt"
	"strings"

	"github.com/joselitofilho/aws-terraform-generator/internal/fmtcolor"
	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
	"github.com/joselitofilho/aws-terraform-generator/internal/generators/terraform"
	"github.com/joselitofilho/aws-terraform-generator/internal/resources"
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
	lambdaResourcesByLabel    map[string]resources.Resource
	restfulAPIResourcesByName map[string]resources.Resource
	s3BucketResourcesByName   map[string]resources.Resource
	sqsResourcesByName        map[string]resources.Resource

	endpointAPIGatewayMap   map[ResourceARN][]ResourceARN
	resourceAPIGIntegration map[ResourceARN]ResourceARN

	relationshipsMap map[ResourceARN][]ResourceARN

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
		lambdaResourcesByLabel:    map[string]resources.Resource{},
		restfulAPIResourcesByName: map[string]resources.Resource{},
		s3BucketResourcesByName:   map[string]resources.Resource{},
		sqsResourcesByName:        map[string]resources.Resource{},

		endpointAPIGatewayMap:   map[ResourceARN][]ResourceARN{},
		resourceAPIGIntegration: map[ResourceARN]ResourceARN{},

		relationshipsMap: map[ResourceARN][]ResourceARN{},

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

func (t *Transformer) getResourceByARN(arn ResourceARN) (resource resources.Resource) {
	switch arn.Key {
	case arnAPIGateway:
		resource = t.apiGatewayResourcesByName[arn.Name]
	case arnCloudwatchKey:
		resource = t.cronResourcesByName[arn.Name]
	case arnEndpoint:
		resource = t.endpointResourcesByName[arn.Name]
	case arnKinesisKey:
		resource = t.kinesisResourcesByName[arn.Name]
	case arnLambdaKey:
		if arn.Label == "" {
			resource = t.lambdaResourcesByName[arn.Name]
		} else {
			resource = t.lambdaResourcesByLabel[arn.Label]
		}
	case arnSQSKey:
		resource = t.sqsResourcesByName[arn.Name]
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
				t.processAPIGatewayRoute(tfResourceConf)
			case labelAWSAPIGatewayIntegration:
				t.processAPIGatewayIntegration(tfResourceConf)
			case labelAWSCloudwatchEventTarget:
				t.processCloudwatchEventTarget(tfResourceConf)
			case labelAWSCron:
				t.processCronResource(tfResourceConf)
			case labelAWSEndpoint:
				t.processEndpointResource(tfResourceConf)
			case labelAWSKinesisStream:
				t.processKinesisResource(tfResourceConf)
			case labelAWSLambdaEventSourceMapping:
				t.processEventSourceMapping(tfResourceConf)
			case labelAWSLambdaFunction:
				t.processLambdaResource(tfResourceConf)
			case labelAWSS3Bucket:
				t.processS3BucketResource(tfResourceConf)
			case labelAWSSQSQueue:
				t.processSQSResource(tfResourceConf)
			}
		}
	}
}

func (t *Transformer) processAPIGatewayRoute(conf *terraform.Resource) {
	routeKeyValue := replaceVars(conf.Attributes["route_key"].(string), t.tfConfig.Variables, t.tfConfig.Locals,
		t.yamlConfig.Draw.ReplaceableTexts)
	routeKeyValue = ResourceByARN(routeKeyValue).Name

	_, ok := t.apiGatewayResourcesByName[routeKeyValue]
	if !ok {
		resource := resources.NewGenericResource(fmt.Sprintf("%d", t.id), routeKeyValue, resources.APIGatewayType)
		t.id++

		t.resources = append(t.resources, resource)
		t.apiGatewayResourcesByName[routeKeyValue] = resource
	}

	apiIDValue := replaceVars(conf.Attributes["api_id"].(string), t.tfConfig.Variables, t.tfConfig.Locals,
		t.yamlConfig.Draw.ReplaceableTexts)
	apiIDARN := ResourceByARN(apiIDValue)

	routeKeyARNlabel := conf.Labels[0]
	routeKeyARN := ResourceARN{Key: strings.Split(routeKeyARNlabel, "_")[2], Name: routeKeyValue,
		Label: routeKeyARNlabel}

	targetValue := replaceVars(conf.Attributes["target"].(string), t.tfConfig.Variables, t.tfConfig.Locals,
		t.yamlConfig.Draw.ReplaceableTexts)
	targetARN := ResourceARN{Key: strings.Split(targetValue, ".")[1]}

	t.relationshipsMap[apiIDARN] = append(t.relationshipsMap[apiIDARN], routeKeyARN)
	t.endpointAPIGatewayMap[targetARN] = append(t.endpointAPIGatewayMap[targetARN], routeKeyARN)
}

func (t *Transformer) processAPIGatewayIntegration(conf *terraform.Resource) {
	label := ResourceARN{Key: conf.Labels[1]}

	apiIDValue := replaceVars(conf.Attributes["api_id"].(string), t.tfConfig.Variables, t.tfConfig.Locals,
		t.yamlConfig.Draw.ReplaceableTexts)
	apiIDARN := ResourceByARN(apiIDValue)

	integrationURIValue := replaceVars(conf.Attributes["integration_uri"].(string), t.tfConfig.Variables,
		t.tfConfig.Locals, t.yamlConfig.Draw.ReplaceableTexts)
	integrationURIARN := ResourceByARN(integrationURIValue)

	// TODO: tryToCreateLambdaResourceByARN

	t.relationshipsMap[apiIDARN] = append(t.relationshipsMap[apiIDARN], integrationURIARN)

	t.resourceAPIGIntegration[integrationURIARN] = label
}

func (t *Transformer) processCloudwatchEventTarget(conf *terraform.Resource) {
	ruleValue := replaceVars(conf.Attributes["rule"].(string), t.tfConfig.Variables, t.tfConfig.Locals,
		t.yamlConfig.Draw.ReplaceableTexts)
	ruleARN := ResourceByARN(ruleValue)

	arnValue := replaceVars(conf.Attributes["arn"].(string), t.tfConfig.Variables, t.tfConfig.Locals,
		t.yamlConfig.Draw.ReplaceableTexts)
	arn := ResourceByARN(arnValue)

	// TODO: tryToCreateLambdaResourceByARN

	t.relationshipsMap[ruleARN] = append(t.relationshipsMap[ruleARN], arn)
}

func (t *Transformer) processCronResource(conf *terraform.Resource) {
	value, ok := conf.Attributes["schedule_expression"]
	if !ok {
		fmtcolor.Yellow.Printf("it is not cron: %s\n", conf.Labels)
		return
	}

	label := conf.Labels[1]
	if _, ok := t.cronResourcesByName[label]; !ok {
		resource := resources.NewGenericResource(fmt.Sprintf("%d", t.id),
			ResourceByARN(value.(string)).Name, resources.CronType)
		t.id++

		t.resources = append(t.resources, resource)
		t.cronResourcesByName[label] = resource
	}
}

func (t *Transformer) processDBResourceFromEnvar(
	k, v string, resourcesByName map[string]resources.Resource,
) resources.Resource {
	return t.processResourceFromEnvar(
		k, v, resources.EnvarSuffixDBHost, resources.DatabaseType, resources.ToDatabaseCase, resourcesByName)
}

func (t *Transformer) processEndpointResource(conf *terraform.Resource) {
	label := conf.Labels[1]
	if _, ok := t.endpointResourcesByName[label]; !ok {
		value := replaceVars(conf.Attributes["domain_name"].(string), t.tfConfig.Variables, t.tfConfig.Locals,
			t.yamlConfig.Draw.ReplaceableTexts)
		value = ResourceByARN(value).Name

		resource := resources.NewGenericResource(fmt.Sprintf("%d", t.id), value, resources.EndpointType)
		t.id++

		t.resources = append(t.resources, resource)
		t.endpointResourcesByName[label] = resource
	}
}

func (t *Transformer) processEventSourceMapping(conf *terraform.Resource) {
	eventSourceValue := replaceVars(conf.Attributes["event_source_arn"].(string), t.tfConfig.Variables,
		t.tfConfig.Locals, t.yamlConfig.Draw.ReplaceableTexts)
	eventSourceARN := ResourceByARN(eventSourceValue)

	functionNameValue := replaceVars(conf.Attributes["function_name"].(string), t.tfConfig.Variables, t.tfConfig.Locals,
		t.yamlConfig.Draw.ReplaceableTexts)
	functionNameARN := ResourceByARN(functionNameValue)

	t.relationshipsMap[eventSourceARN] = append(t.relationshipsMap[eventSourceARN], functionNameARN)

	// TODO: tryToCreateSQSResourceByARN
	// TODO: tryToCreateLambdaResourceByARN

	t.tryToCreateResourceByARN(eventSourceARN)
}

func (t *Transformer) processGoogleBQResourceFromEnvar(
	k, v string, resourcesByName map[string]resources.Resource,
) resources.Resource {
	return t.processResourceFromEnvar(
		k, v, resources.EnvarSuffixGoogleBQ, resources.GoogleBQType, resources.ToGoogleBQCase, resourcesByName)
}

func (t *Transformer) processKinesisResourceFromEnvar(
	k, v string, resourcesByName map[string]resources.Resource,
) resources.Resource {
	return t.processResourceFromEnvar(
		k, v, resources.EnvarSuffixKinesisStreamURL, resources.KinesisType, resources.ToKinesisCase, resourcesByName)
}

func (t *Transformer) processKinesisResource(conf *terraform.Resource) {
	t.processResource(conf, resources.KinesisType, "name", suffixKinesis, resources.ToKinesisCase,
		t.kinesisResourcesByName)
}

func (t *Transformer) processLambda(attributes, envars map[string]any, label string) {
	value := strTransformFromKeyValue(label, label, suffixLambda, resources.ToLambdaCase)

	for k, v := range attributes {
		if strings.HasSuffix(k, "function_name") {
			value = replaceVars(v.(string), t.tfConfig.Variables, t.tfConfig.Locals,
				t.yamlConfig.Draw.ReplaceableTexts)
			value = ResourceByARN(value).Name
			value = strTransformFromKeyValue(value, value, suffixLambda, resources.ToLambdaCase)

			break
		}
	}

	resource, ok := t.lambdaResourcesByName[value]
	if !ok {
		resource = resources.NewGenericResource(fmt.Sprintf("%d", t.id), value, resources.LambdaType)
		t.id++

		t.resources = append(t.resources, resource)
		t.lambdaResourcesByName[value] = resource
		t.lambdaResourcesByLabel[label] = resource
	}

	for k, v := range envars {
		switch {
		case strings.HasSuffix(k, resources.EnvarSuffixDBHost):
			target := t.processDBResourceFromEnvar(k, v.(string), t.dbResourcesByName)
			t.relationships = append(t.relationships,
				resources.Relationship{Source: resource, Target: target})
		case strings.HasSuffix(k, resources.EnvarSuffixGoogleBQ):
			target := t.processGoogleBQResourceFromEnvar(k, v.(string), t.googleBQResourcesByName)
			t.relationships = append(t.relationships,
				resources.Relationship{Source: resource, Target: target})
		case strings.HasSuffix(k, resources.EnvarSuffixKinesisStreamURL):
			target := t.processKinesisResourceFromEnvar(k, v.(string), t.kinesisResourcesByName)
			t.relationships = append(t.relationships,
				resources.Relationship{Source: resource, Target: target})
		case strings.HasSuffix(k, resources.EnvarSuffixS3BucketURL):
			target := t.processS3BucketResourceFromEnvar(
				k, v.(string), resources.EnvarSuffixS3BucketURL, t.s3BucketResourcesByName)
			t.relationships = append(t.relationships,
				resources.Relationship{Source: resource, Target: target})
		case strings.HasSuffix(k, resources.EnvarSuffixS3BucketName):
			target := t.processS3BucketResourceFromEnvar(
				k, v.(string), resources.EnvarSuffixS3BucketName, t.s3BucketResourcesByName)
			t.relationships = append(t.relationships,
				resources.Relationship{Source: resource, Target: target})
		case strings.HasSuffix(k, resources.EnvarSuffixSQSQueueURL):
			target := t.processSQSResourceFromEnvar(k, v.(string), t.sqsResourcesByName)
			t.relationships = append(t.relationships,
				resources.Relationship{Source: resource, Target: target})
		case strings.HasSuffix(k, resources.EnvarSuffixRestfulAPI):
			target := t.processRestfulAPIResourceFromEnvar(k, v.(string), t.restfulAPIResourcesByName)
			t.relationships = append(t.relationships,
				resources.Relationship{Source: resource, Target: target})
		}
	}
}

func (t *Transformer) processLambdaModule(conf *terraform.Module) {
	envars := map[string]any{}
	if vars, ok := conf.Attributes["lambda_function_env_vars"]; ok {
		envars = vars.(map[string]any)
	}

	t.processLambda(conf.Attributes, envars, conf.Labels[0])
}

func (t *Transformer) processLambdaResource(conf *terraform.Resource) {
	envars := map[string]any{}

	if environment, ok := conf.Attributes["environment"]; ok {
		if vars, ok := environment.(map[string]map[string]any)["variables"]; ok {
			for k, v := range vars {
				envars[k] = v
			}
		}
	}

	t.processLambda(conf.Attributes, envars, conf.Labels[1])
}

func (t *Transformer) processResource(
	conf *terraform.Resource, resourceType resources.ResourceType, attributeName string, suffix string,
	caseTransformer func(string) string, resourcesMap map[string]resources.Resource,
) {
	value := replaceVars(conf.Attributes[attributeName].(string), t.tfConfig.Variables, t.tfConfig.Locals,
		t.yamlConfig.Draw.ReplaceableTexts)
	value = ResourceByARN(value).Name
	value = strTransformFromKeyValue(value, value, suffix, caseTransformer)

	if _, ok := resourcesMap[value]; !ok {
		resource := resources.NewGenericResource(fmt.Sprintf("%d", t.id), value, resourceType)
		t.id++

		t.resources = append(t.resources, resource)
		resourcesMap[value] = resource
	}
}

func (t *Transformer) processS3BucketResourceFromEnvar(
	k, v, suffix string, resourcesByName map[string]resources.Resource,
) resources.Resource {
	value := replaceVars(v, t.tfConfig.Variables, t.tfConfig.Locals,
		t.yamlConfig.Draw.ReplaceableTexts)
	value = ResourceByARN(value).Name
	value = strTransformFromKeyValue(k, value, suffix, resources.ToS3BucketCase)

	resource, ok := resourcesByName[value]

	if !ok {
		resource = resources.NewGenericResource(fmt.Sprintf("%d", t.id), value, resources.S3Type)
		t.id++

		resourcesByName[value] = resource
		t.resources = append(t.resources, resource)
	}

	return resource
}

func (t *Transformer) processS3BucketResource(conf *terraform.Resource) {
	t.processResource(conf, resources.S3Type, "bucket", suffixS3Bucket, resources.ToS3BucketCase,
		t.s3BucketResourcesByName)
}

func (t *Transformer) processSQSResourceFromEnvar(
	k, v string, resourcesByName map[string]resources.Resource,
) resources.Resource {
	return t.processResourceFromEnvar(
		k, v, resources.EnvarSuffixSQSQueueURL, resources.SQSType, resources.ToSQSCase, resourcesByName)
}

func (t *Transformer) processSQSResource(conf *terraform.Resource) {
	t.processResource(conf, resources.SQSType, "name", suffixSQS, resources.ToSQSCase, t.sqsResourcesByName)
}

func (t *Transformer) processRestfulAPIResourceFromEnvar(
	k, v string, resourcesByName map[string]resources.Resource,
) resources.Resource {
	return t.processResourceFromEnvar(
		k, v, resources.EnvarSuffixRestfulAPI, resources.RestfulAPIType, resources.ToRestfulAPICase, resourcesByName)
}

func (t *Transformer) tryToCreateResourceByARN(eventSourceARN ResourceARN) {
	if eventSourceARN.Key == arnKinesisKey {
		value := replaceVars(eventSourceARN.Name, t.tfConfig.Variables, t.tfConfig.Locals,
			t.yamlConfig.Draw.ReplaceableTexts)
		value = ResourceByARN(value).Name
		value = strTransformFromKeyValue(eventSourceARN.Name,
			value, resources.EnvarSuffixKinesisStreamURL, resources.ToKinesisCase)

		if _, ok := t.kinesisResourcesByName[value]; !ok {
			resource := resources.NewGenericResource(fmt.Sprintf("%d", t.id), value, resources.KinesisType)
			t.id++

			t.resources = append(t.resources, resource)
			t.kinesisResourcesByName[value] = resource
		}
	}
}

func (t *Transformer) processResourceFromEnvar(
	k, v, suffix string, restType resources.ResourceType, fn func(s string) string,
	resourcesByName map[string]resources.Resource,
) resources.Resource {
	value := replaceVars(v, t.tfConfig.Variables, t.tfConfig.Locals,
		t.yamlConfig.Draw.ReplaceableTexts)
	value = ResourceByARN(value).Name
	value = strTransformFromKeyValue(k, value, suffix, fn)

	resource, ok := resourcesByName[value]

	if !ok {
		resource = resources.NewGenericResource(fmt.Sprintf("%d", t.id), value, restType)
		t.id++

		resourcesByName[value] = resource
		t.resources = append(t.resources, resource)
	}

	return resource
}
