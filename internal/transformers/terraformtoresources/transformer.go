package terraformtoresources

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/diagram-code-generator/resources/pkg/resources"

	"github.com/joselitofilho/aws-terraform-generator/internal/fmtcolor"
	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
	"github.com/joselitofilho/aws-terraform-generator/internal/generators/terraform"
	awsresources "github.com/joselitofilho/aws-terraform-generator/internal/resources"
)

const suffixLambda = "_lambda"

type Transformer struct {
	yamlConfig *config.Config
	tfConfig   *terraform.Config

	resources     []resources.Resource
	relationships []resources.Relationship

	apiGatewayResourcesByName map[string]resources.Resource
	dbResourcesByName         map[string]resources.Resource
	googleBQResourcesByName   map[string]resources.Resource
	kinesisResourcesByName    map[string]resources.Resource
	lambdaResourcesByName     map[string]resources.Resource
	restfulAPIResourcesByName map[string]resources.Resource
	s3BucketResourcesByName   map[string]resources.Resource
	sqsResourcesByName        map[string]resources.Resource

	cronResourcesByLabel     map[string]resources.Resource
	endpointResourcesByLabel map[string]resources.Resource
	kinesisResourcesByLabel  map[string]resources.Resource
	lambdaResourcesByLabel   map[string]resources.Resource
	s3BucketResourcesByLabel map[string]resources.Resource
	sqsResourcesByLabel      map[string]resources.Resource

	apigIntegrationRouteMap map[awsresources.ResourceARN][]awsresources.ResourceARN
	resourceAPIGIntegration map[awsresources.ResourceARN]awsresources.ResourceARN

	relationshipsMap map[awsresources.ResourceARN][]awsresources.ResourceARN

	id int
}

func NewTransformer(yamlConfig *config.Config, tfConfig *terraform.Config) *Transformer {
	return &Transformer{
		yamlConfig: yamlConfig,
		tfConfig:   tfConfig,

		resources:     []resources.Resource{},
		relationships: []resources.Relationship{},

		apiGatewayResourcesByName: map[string]resources.Resource{},
		dbResourcesByName:         map[string]resources.Resource{},
		googleBQResourcesByName:   map[string]resources.Resource{},
		kinesisResourcesByName:    map[string]resources.Resource{},
		lambdaResourcesByName:     map[string]resources.Resource{},
		restfulAPIResourcesByName: map[string]resources.Resource{},
		s3BucketResourcesByName:   map[string]resources.Resource{},
		sqsResourcesByName:        map[string]resources.Resource{},

		cronResourcesByLabel:     map[string]resources.Resource{},
		endpointResourcesByLabel: map[string]resources.Resource{},
		kinesisResourcesByLabel:  map[string]resources.Resource{},
		lambdaResourcesByLabel:   map[string]resources.Resource{},
		s3BucketResourcesByLabel: map[string]resources.Resource{},
		sqsResourcesByLabel:      map[string]resources.Resource{},

		apigIntegrationRouteMap: map[awsresources.ResourceARN][]awsresources.ResourceARN{},
		resourceAPIGIntegration: map[awsresources.ResourceARN]awsresources.ResourceARN{},

		relationshipsMap: map[awsresources.ResourceARN][]awsresources.ResourceARN{},

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
				for _, apig := range t.apigIntegrationRouteMap[integration] {
					updatedSource := t.getResourceByARN(apig)

					t.relationships = append(t.relationships, resources.Relationship{Source: updatedSource, Target: target})
				}

				continue
			}

			t.relationships = append(t.relationships, resources.Relationship{Source: source, Target: target})
		}
	}
}

func (t *Transformer) getResourceByARN(arn awsresources.ResourceARN) (resource resources.Resource) {
	switch arn.Type {
	case awsresources.LabelAWSAPIGatewayAPI:
		resource = t.endpointResourcesByLabel[arn.Label]
	case awsresources.LabelAWSAPIGatewayRoute:
		resource = t.apiGatewayResourcesByName[arn.Name]
	case awsresources.LabelAWSCron:
		resource = t.cronResourcesByLabel[arn.Label]
	case awsresources.LabelAWSEndpoint:
		resource = t.endpointResourcesByLabel[arn.Label]
	case awsresources.LabelAWSKinesisStream:
		if arn.Label == "" {
			resource = t.kinesisResourcesByName[arn.Name]
		} else {
			resource = t.kinesisResourcesByLabel[arn.Label]
		}
	case awsresources.LabelAWSLambdaFunction:
		if arn.Label == "" {
			resource = t.lambdaResourcesByName[arn.Name]
		} else {
			resource = t.lambdaResourcesByLabel[arn.Label]
		}
	case awsresources.LabelAWSS3Bucket:
		if arn.Label == "" {
			resource = t.s3BucketResourcesByName[arn.Name]
		} else {
			resource = t.s3BucketResourcesByLabel[arn.Label]
		}
	case awsresources.LabelAWSSQSQueue:
		if arn.Label == "" {
			resource = t.sqsResourcesByName[arn.Name]
		} else {
			resource = t.sqsResourcesByLabel[arn.Label]
		}
	}

	return resource
}

func (t *Transformer) hasResourceMatched(res resources.Resource, filters config.Filters) bool {
	if res == nil {
		return false
	}

	filter, hasFilter := filters[awsresources.ParseResourceType(res.ResourceType())]
	if !hasFilter {
		return true
	}

	match := len(filter.Match) == 0

	for _, pattern := range filter.Match {
		regex, err := regexp.Compile(pattern)
		if err != nil {
			fmtcolor.Yellow.Println("error compiling match regex:", err)
			continue
		}

		if regex.MatchString(res.Value()) {
			match = true
			break
		}
	}

	for _, pattern := range filter.NotMatch {
		regex, err := regexp.Compile(pattern)
		if err != nil {
			fmtcolor.Yellow.Println("error compiling not_match regex:", err)
			continue
		}

		if regex.MatchString(res.Value()) {
			match = false
			break
		}
	}

	return match
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
			case awsresources.LabelAWSAPIGatewayRoute:
				t.processAPIGatewayRoute(tfResourceConf)
			case awsresources.LabelAWSAPIGatewayIntegration:
				t.processAPIGatewayIntegration(tfResourceConf)
			case awsresources.LabelAWSCloudwatchEventTarget:
				t.processCloudwatchEventTarget(tfResourceConf)
			case awsresources.LabelAWSCron:
				t.processCronResource(tfResourceConf)
			case awsresources.LabelAWSEndpoint:
				t.processEndpointResource(tfResourceConf)
			case awsresources.LabelAWSKinesisStream:
				t.processKinesisResource(tfResourceConf)
			case awsresources.LabelAWSLambdaEventSourceMapping:
				t.processEventSourceMapping(tfResourceConf)
			case awsresources.LabelAWSLambdaFunction:
				t.processLambdaResource(tfResourceConf)
			case awsresources.LabelAWSS3Bucket:
				t.processS3BucketResource(tfResourceConf)
			case awsresources.LabelAWSSQSQueue:
				t.processSQSResource(tfResourceConf)
			}
		}
	}
}

func (t *Transformer) processAPIGatewayRoute(conf *terraform.Resource) {
	routeKeyValue := replaceVars(conf.Attributes["route_key"].(string), t.tfConfig.Variables, t.tfConfig.Locals,
		t.yamlConfig.Draw.ReplaceableTexts)

	routeKeyARN := awsresources.ParseResourceARN(routeKeyValue, awsresources.APIGatewayType)
	routeKeyARN.Label = conf.Labels[1]

	routeKeyValue = routeKeyARN.Name

	_, ok := t.apiGatewayResourcesByName[routeKeyValue]
	if !ok {
		resource := resources.NewGenericResource(fmt.Sprintf("%d", t.id), routeKeyValue, awsresources.APIGatewayType.String())
		t.id++

		t.resources = append(t.resources, resource)
		t.apiGatewayResourcesByName[routeKeyValue] = resource
	}

	apiIDValue := replaceVars(conf.Attributes["api_id"].(string), t.tfConfig.Variables, t.tfConfig.Locals,
		t.yamlConfig.Draw.ReplaceableTexts)
	apiIDARN := awsresources.ParseResourceARN(apiIDValue, awsresources.EndpointType)

	targetValue := replaceVars(conf.Attributes["target"].(string), t.tfConfig.Variables, t.tfConfig.Locals,
		t.yamlConfig.Draw.ReplaceableTexts)
	targetValue = strings.ReplaceAll(strings.ReplaceAll(targetValue, "${", ""), "}", "")
	targetValueParts := strings.Split(strings.Split(targetValue, "/")[1], ".")
	targetARN := awsresources.ResourceARN{Type: targetValueParts[0], Label: targetValueParts[1]}

	t.relationshipsMap[apiIDARN] = append(t.relationshipsMap[apiIDARN], routeKeyARN)
	t.apigIntegrationRouteMap[targetARN] = append(t.apigIntegrationRouteMap[targetARN], routeKeyARN)
}

func (t *Transformer) processAPIGatewayIntegration(conf *terraform.Resource) {
	integrationARN := awsresources.ResourceARN{Type: conf.Labels[0], Label: conf.Labels[1]}

	apiIDValue := replaceVars(conf.Attributes["api_id"].(string), t.tfConfig.Variables, t.tfConfig.Locals,
		t.yamlConfig.Draw.ReplaceableTexts)
	apiIDARN := awsresources.ParseResourceARN(apiIDValue, awsresources.EndpointType)

	integrationURIValue := replaceVars(conf.Attributes["integration_uri"].(string), t.tfConfig.Variables,
		t.tfConfig.Locals, t.yamlConfig.Draw.ReplaceableTexts)
	integrationURIARN := awsresources.ParseResourceARN(integrationURIValue, awsresources.LambdaType)

	t.relationshipsMap[apiIDARN] = append(t.relationshipsMap[apiIDARN], integrationURIARN)
	t.resourceAPIGIntegration[integrationURIARN] = integrationARN
}

func (t *Transformer) processCloudwatchEventTarget(conf *terraform.Resource) {
	t.processResourceRelationships(conf, "rule", "arn", awsresources.CronType, awsresources.LambdaType)
}

func (t *Transformer) processCronResource(conf *terraform.Resource) {
	value, ok := conf.Attributes["schedule_expression"]
	if !ok {
		fmtcolor.Yellow.Printf("it is not cron: %s\n", conf.Labels)
		return
	}

	label := conf.Labels[1]
	if _, ok := t.cronResourcesByLabel[label]; !ok {
		resType := awsresources.CronType
		resource := resources.NewGenericResource(fmt.Sprintf("%d", t.id),
			awsresources.ParseResourceARN(value.(string), resType).Name, resType.String())
		t.id++

		t.resources = append(t.resources, resource)
		t.cronResourcesByLabel[label] = resource
	}
}

func (t *Transformer) processDBResourceFromEnvar(
	v string, resourcesByName map[string]resources.Resource,
) resources.Resource {
	return t.processResourceFromEnvar(v, awsresources.DatabaseType, resourcesByName)
}

func (t *Transformer) processEndpointResource(conf *terraform.Resource) {
	label := conf.Labels[1]
	if _, ok := t.endpointResourcesByLabel[label]; !ok {
		value := replaceVars(conf.Attributes["domain_name"].(string), t.tfConfig.Variables, t.tfConfig.Locals,
			t.yamlConfig.Draw.ReplaceableTexts)
		value = awsresources.ParseResourceARN(value, awsresources.EndpointType).Name

		resource := resources.NewGenericResource(fmt.Sprintf("%d", t.id), value, awsresources.EndpointType.String())
		t.id++

		t.resources = append(t.resources, resource)
		t.endpointResourcesByLabel[label] = resource
	}
}

func (t *Transformer) processEventSourceMapping(conf *terraform.Resource) {
	t.processResourceRelationships(conf, "event_source_arn", "function_name",
		awsresources.UnknownType, awsresources.LambdaType)
}

func (t *Transformer) processGoogleBQResourceFromEnvar(
	v string, resourcesByName map[string]resources.Resource,
) resources.Resource {
	return t.processResourceFromEnvar(v, awsresources.GoogleBQType, resourcesByName)
}

func (t *Transformer) processKinesisResource(conf *terraform.Resource) {
	t.processResource(conf, awsresources.KinesisType, "name", t.kinesisResourcesByName, t.kinesisResourcesByLabel)
}

func (t *Transformer) processLambda(attributes, envars map[string]any, label string) {
	var (
		name string

		restType = awsresources.LambdaType
	)

	for k, v := range attributes {
		if strings.HasSuffix(k, "function_name") {
			value := replaceVars(v.(string), t.tfConfig.Variables, t.tfConfig.Locals,
				t.yamlConfig.Draw.ReplaceableTexts)
			name = awsresources.ParseResourceARN(value, restType).Name

			break
		}
	}

	if name == "" {
		// TODO: Review and create a test for this.
		return
	}

	resource, ok := t.lambdaResourcesByName[name]
	if !ok {
		resource = resources.NewGenericResource(fmt.Sprintf("%d", t.id), name, restType.String())
		t.id++

		t.resources = append(t.resources, resource)
		t.lambdaResourcesByName[name] = resource
		t.lambdaResourcesByLabel[label] = resource
	}

	lambdaARN := awsresources.ResourceARN{
		Type: awsresources.LabelAWSLambdaFunction, Name: resource.Value(), Label: label}

	for k, v := range envars {
		switch {
		case strings.HasSuffix(k, awsresources.EnvarSuffixDBHost):
			target := t.processDBResourceFromEnvar(v.(string), t.dbResourcesByName)
			t.relationships = append(t.relationships,
				resources.Relationship{Source: resource, Target: target})
		case strings.HasSuffix(k, awsresources.EnvarSuffixGoogleBQ):
			target := t.processGoogleBQResourceFromEnvar(v.(string), t.googleBQResourcesByName)
			t.relationships = append(t.relationships,
				resources.Relationship{Source: resource, Target: target})
		case strings.HasSuffix(k, awsresources.EnvarSuffixKinesisStreamURL):
			targetArn := t.processResourceARNFromEnvar(v.(string), awsresources.KinesisType)
			t.relationshipsMap[lambdaARN] = append(t.relationshipsMap[lambdaARN], targetArn)
		case strings.HasSuffix(k, awsresources.EnvarSuffixRestfulAPI):
			target := t.processRestfulAPIResourceFromEnvar(v.(string), t.restfulAPIResourcesByName)
			t.relationships = append(t.relationships,
				resources.Relationship{Source: resource, Target: target})
		case strings.HasSuffix(k, awsresources.EnvarSuffixS3BucketURL),
			strings.HasSuffix(k, awsresources.EnvarSuffixS3BucketName):
			targetArn := t.processResourceARNFromEnvar(v.(string), awsresources.S3Type)
			t.relationshipsMap[lambdaARN] = append(t.relationshipsMap[lambdaARN], targetArn)
		case strings.HasSuffix(k, awsresources.EnvarSuffixSQSQueueURL):
			targetArn := t.processResourceARNFromEnvar(v.(string), awsresources.SQSType)
			t.relationshipsMap[lambdaARN] = append(t.relationshipsMap[lambdaARN], targetArn)
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
	conf *terraform.Resource, resourceType awsresources.ResourceType, attributeName string,
	resourcesByName, resourcesByLabel map[string]resources.Resource,
) {
	label := conf.Labels[1]
	value := replaceVars(conf.Attributes[attributeName].(string), t.tfConfig.Variables, t.tfConfig.Locals,
		t.yamlConfig.Draw.ReplaceableTexts)
	resourceARN := awsresources.ParseResourceARN(value, resourceType)
	resourceARN.Label = label

	name := resourceARN.Name
	if _, ok := resourcesByName[name]; !ok {
		resource := resources.NewGenericResource(fmt.Sprintf("%d", t.id), name, resourceType.String())
		t.id++

		t.resources = append(t.resources, resource)

		resourcesByName[name] = resource
		resourcesByLabel[label] = resource
	}
}

func (t *Transformer) processResourceRelationships(
	conf *terraform.Resource, sourceAttribute string, targetAttribute string,
	sourceType awsresources.ResourceType, targetType awsresources.ResourceType,
) {
	sourceValue := replaceVars(conf.Attributes[sourceAttribute].(string), t.tfConfig.Variables, t.tfConfig.Locals,
		t.yamlConfig.Draw.ReplaceableTexts)
	sourceARN := awsresources.ParseResourceARN(sourceValue, sourceType)

	targetValue := replaceVars(conf.Attributes[targetAttribute].(string), t.tfConfig.Variables, t.tfConfig.Locals,
		t.yamlConfig.Draw.ReplaceableTexts)
	targetARN := awsresources.ParseResourceARN(targetValue, targetType)

	t.relationshipsMap[sourceARN] = append(t.relationshipsMap[sourceARN], targetARN)
}

func (t *Transformer) processS3BucketResource(conf *terraform.Resource) {
	t.processResource(conf, awsresources.S3Type, "bucket", t.s3BucketResourcesByName, t.s3BucketResourcesByLabel)
}

func (t *Transformer) processSQSResource(conf *terraform.Resource) {
	t.processResource(conf, awsresources.SQSType, "name", t.sqsResourcesByName, t.sqsResourcesByLabel)
}

func (t *Transformer) processRestfulAPIResourceFromEnvar(
	v string, resourcesByName map[string]resources.Resource,
) resources.Resource {
	return t.processResourceFromEnvar(v, awsresources.RestfulAPIType, resourcesByName)
}

func (t *Transformer) processResourceFromEnvar(
	v string, restType awsresources.ResourceType, resourcesByName map[string]resources.Resource,
) resources.Resource {
	value := replaceVars(v, t.tfConfig.Variables, t.tfConfig.Locals,
		t.yamlConfig.Draw.ReplaceableTexts)
	value = awsresources.ParseResourceARN(value, restType).Name

	resource, ok := resourcesByName[value]

	if !ok {
		resource = resources.NewGenericResource(fmt.Sprintf("%d", t.id), value, restType.String())
		t.id++

		t.resources = append(t.resources, resource)
		resourcesByName[value] = resource
	}

	return resource
}

func (t *Transformer) processResourceARNFromEnvar(
	v string, restType awsresources.ResourceType,
) awsresources.ResourceARN {
	value := replaceVars(v, t.tfConfig.Variables, t.tfConfig.Locals,
		t.yamlConfig.Draw.ReplaceableTexts)

	return awsresources.ParseResourceARN(value, restType)
}
