package yamltoresources

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/ettle/strcase"

	"github.com/diagram-code-generator/resources/pkg/resources"

	"github.com/joselitofilho/aws-terraform-generator/internal/fmtcolor"
	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
	awsresources "github.com/joselitofilho/aws-terraform-generator/internal/resources"
	"github.com/joselitofilho/aws-terraform-generator/internal/transformers"
)

var ErrEmptyConfig = errors.New("config file is empty")

type Transformer struct {
	yamlConfig *config.Config

	apigatewayByName map[string]resources.Resource
	cronByName       map[string]resources.Resource
	databaseByName   map[string]resources.Resource
	endpointByName   map[string]resources.Resource
	googleBQByName   map[string]resources.Resource
	kinesisByName    map[string]resources.Resource
	lambdaByName     map[string]resources.Resource
	restfulAPIByName map[string]resources.Resource
	s3BucketByName   map[string]resources.Resource
	snsByName        map[string]resources.Resource
	sqsByName        map[string]resources.Resource

	relationshipsMap map[awsresources.ResourceARN][]awsresources.ResourceARN
}

func NewTransformer(yamlConfig *config.Config) *Transformer {
	return &Transformer{
		yamlConfig: yamlConfig,

		apigatewayByName: map[string]resources.Resource{},
		cronByName:       map[string]resources.Resource{},
		databaseByName:   map[string]resources.Resource{},
		endpointByName:   map[string]resources.Resource{},
		googleBQByName:   map[string]resources.Resource{},
		kinesisByName:    map[string]resources.Resource{},
		lambdaByName:     map[string]resources.Resource{},
		restfulAPIByName: map[string]resources.Resource{},
		s3BucketByName:   map[string]resources.Resource{},
		snsByName:        map[string]resources.Resource{},
		sqsByName:        map[string]resources.Resource{},

		relationshipsMap: map[awsresources.ResourceARN][]awsresources.ResourceARN{},
	}
}

func (t *Transformer) Transform() (*resources.ResourceCollection, error) {
	if t.yamlConfig == nil {
		return nil, fmt.Errorf("%w", ErrEmptyConfig)
	}

	id := 1
	rscs := []resources.Resource{}
	relationships := []resources.Relationship{}

	t.transformAPIGateways(&rscs, &relationships, &id)
	t.extractKinesisResources(&rscs, &id)
	t.transformLambdas(&rscs, &relationships, &id)
	t.extractRestfulAPIResources(&rscs, &id)
	t.extractS3BucketResources(&rscs, &id)
	t.extractSNSBucketResources(&rscs, &id)
	t.extractSQSResources(&rscs, &id)

	t.buildRelationships(&relationships)

	return &resources.ResourceCollection{
		Resources:     rscs,
		Relationships: relationships,
	}, nil
}

func (t *Transformer) buildRelationships(relationships *[]resources.Relationship) {
	for sourceARN, rel := range t.relationshipsMap {
		source := t.getResourceByARN(sourceARN)

		for i := range rel {
			target := t.getResourceByARN(rel[i])

			*relationships = append(*relationships, resources.Relationship{Source: source, Target: target})
		}
	}
}

func (t *Transformer) getResourceByARN(arn awsresources.ResourceARN) (resource resources.Resource) {
	key := arn.LabelOrName()

	switch arn.Type {
	case awsresources.LabelAWSAPIGatewayAPI:
		resource = t.endpointByName[key]
	case awsresources.LabelAWSAPIGatewayRoute:
		resource = t.apigatewayByName[key]
	case awsresources.LabelAWSCron:
		resource = t.cronByName[key]
	case awsresources.LabelAWSEndpoint:
		resource = t.endpointByName[key]
	case awsresources.LabelAWSKinesisStream:
		resource = t.kinesisByName[key]
	case awsresources.LabelAWSLambdaFunction:
		resource = t.lambdaByName[key]
	case awsresources.LabelAWSS3Bucket:
		resource = t.s3BucketByName[key]
	case awsresources.LabelAWSSQSQueue:
		resource = t.sqsByName[key]
	}

	return resource
}

func (t *Transformer) extractResourcesByType(
	resourcesList []config.Resource, resourceType awsresources.ResourceType, resourceMap map[string]resources.Resource,
	rscs *[]resources.Resource, id *int,
) {
	for i := range resourcesList {
		res := resourcesList[i]

		resARN := awsresources.ParseResourceARN(res.GetName(), resourceType)
		if resARN.Label == "" &&
			(resourceType == awsresources.KinesisType || resourceType == awsresources.S3Type || resourceType == awsresources.SQSType) {
			arnType := fmt.Sprintf("%s_%s", strcase.ToSnake(resARN.Name), awsresources.SuffixByResource[resourceType])
			resARN.Label = arnType
		}

		key := resARN.LabelOrName()

		if _, ok := resourceMap[key]; !ok {
			resourceRes := resources.NewGenericResource(fmt.Sprintf("%d", *id), resARN.Name, resourceType.String())
			*id++

			*rscs = append(*rscs, resourceRes)

			resourceMap[key] = resourceRes
		}
	}
}

func (t *Transformer) extractKinesisResources(rscs *[]resources.Resource, id *int) {
	configResources := make([]config.Resource, 0, len(t.yamlConfig.Kinesis))
	for i := range t.yamlConfig.Kinesis {
		configResources = append(configResources,
			reflect.ValueOf(&t.yamlConfig.Kinesis[i]).Interface().(config.Resource))
	}

	t.extractResourcesByType(configResources, awsresources.KinesisType, t.kinesisByName, rscs, id)
}

func (t *Transformer) extractRestfulAPIResources(rscs *[]resources.Resource, id *int) {
	configResources := make([]config.Resource, 0, len(t.yamlConfig.RestfulAPIs))
	for i := range t.yamlConfig.RestfulAPIs {
		configResources = append(configResources,
			reflect.ValueOf(&t.yamlConfig.RestfulAPIs[i]).Interface().(config.Resource))
	}

	t.extractResourcesByType(configResources, awsresources.RestfulAPIType, t.restfulAPIByName, rscs, id)
}

func (t *Transformer) extractS3BucketResources(rscs *[]resources.Resource, id *int) {
	configResources := make([]config.Resource, 0, len(t.yamlConfig.Buckets))
	for i := range t.yamlConfig.Buckets {
		configResources = append(configResources,
			reflect.ValueOf(&t.yamlConfig.Buckets[i]).Interface().(config.Resource))
	}

	t.extractResourcesByType(configResources, awsresources.S3Type, t.s3BucketByName, rscs, id)
}

func (t *Transformer) extractSNSBucketResources(rscs *[]resources.Resource, id *int) {
	configResources := make([]config.Resource, 0, len(t.yamlConfig.SNSs))
	for i := range t.yamlConfig.SNSs {
		configResources = append(configResources,
			reflect.ValueOf(&t.yamlConfig.SNSs[i]).Interface().(config.Resource))
	}

	t.extractResourcesByType(configResources, awsresources.SNSType, t.snsByName, rscs, id)
}

func (t *Transformer) extractSQSResources(rscs *[]resources.Resource, id *int) {
	configResources := make([]config.Resource, 0, len(t.yamlConfig.SQSs))
	for i := range t.yamlConfig.SQSs {
		configResources = append(configResources,
			reflect.ValueOf(&t.yamlConfig.SQSs[i]).Interface().(config.Resource))
	}

	t.extractResourcesByType(configResources, awsresources.SQSType, t.sqsByName, rscs, id)
}

func (t *Transformer) transformAPIGateways(
	rscs *[]resources.Resource, relationships *[]resources.Relationship, id *int,
) {
	for _, res := range t.yamlConfig.APIGateways {
		endpointValue := res.APIDomain

		endpointRes, ok := t.endpointByName[endpointValue]
		if !ok {
			endpointRes = resources.NewGenericResource(
				fmt.Sprintf("%d", *id), endpointValue, awsresources.EndpointType.String())
			*rscs = append(*rscs, endpointRes)
			*id++

			t.endpointByName[endpointValue] = endpointRes
		}

		for i := range res.Lambdas {
			l := res.Lambdas[i]
			apigValue := fmt.Sprintf("%s %s", l.Verb, l.Path)

			apigRes, ok := t.apigatewayByName[apigValue]

			if !ok {
				apigRes = resources.NewGenericResource(
					fmt.Sprintf("%d", *id), apigValue, awsresources.APIGatewayType.String())
				*rscs = append(*rscs, apigRes)
				*id++

				t.apigatewayByName[apigValue] = apigRes
			}

			lambdaName := awsresources.ToLambdaCase(l.Name)

			t.transformLambda(&config.Lambda{Name: lambdaName, Envars: l.Envars}, rscs, relationships, id)

			*relationships = append(*relationships,
				resources.Relationship{
					Source: apigRes,
					Target: t.lambdaByName[lambdaName],
				}, resources.Relationship{
					Source: endpointRes,
					Target: apigRes,
				})
		}
	}
}

func (t *Transformer) transformLambda(
	res *config.Lambda, rscs *[]resources.Resource, relationships *[]resources.Relationship, id *int,
) {
	if _, ok := t.lambdaByName[res.Name]; ok {
		return
	}

	lambdaARN := awsresources.ParseResourceARN(res.Name, awsresources.LambdaType)
	lambdaName := awsresources.ToLambdaCase(lambdaARN.Name)

	lambda, ok := t.lambdaByName[lambdaName]
	if !ok {
		lambda = resources.NewGenericResource(fmt.Sprintf("%d", *id), lambdaName, awsresources.LambdaType.String())
		*id++

		*rscs = append(*rscs, lambda)

		t.lambdaByName[lambdaName] = lambda
	}

	for _, r := range res.Crons {
		name := r.ScheduleExpression

		sourceRes, ok := t.cronByName[name]
		if !ok {
			sourceRes = resources.NewGenericResource(fmt.Sprintf("%d", *id), name, awsresources.CronType.String())
			*id++

			*rscs = append(*rscs, sourceRes)

			t.cronByName[name] = sourceRes
		}

		*relationships = append(*relationships, resources.Relationship{Source: sourceRes, Target: lambda})
	}

	t.transformLambdaEnvars(res, lambda, lambdaARN, rscs, relationships, id)

	for _, r := range res.KinesisTriggers {
		kinesisARN := awsresources.ParseResourceARN(r.SourceARN, awsresources.KinesisType)
		t.relationshipsMap[kinesisARN] = append(t.relationshipsMap[kinesisARN], lambdaARN)
	}

	for _, r := range res.SQSTriggers {
		sqsARN := awsresources.ParseResourceARN(r.SourceARN, awsresources.SQSType)
		t.relationshipsMap[sqsARN] = append(t.relationshipsMap[sqsARN], lambdaARN)
	}
}

func (t *Transformer) transformLambdas(rscs *[]resources.Resource, relationships *[]resources.Relationship, id *int) {
	for i := range t.yamlConfig.Lambdas {
		res := t.yamlConfig.Lambdas[i]

		t.transformLambda(&res, rscs, relationships, id)
	}
}

func (t *Transformer) transformLambdaEnvars(
	res *config.Lambda, lambda resources.Resource, lambdaARN awsresources.ResourceARN,
	rscs *[]resources.Resource, relationships *[]resources.Relationship, id *int,
) {
	for k, v := range res.Envars {
		value, resType := t.getValueTypeFromEnvar(k)

		switch resType {
		case awsresources.DatabaseType:
			t.fromLambdaToResource(value, lambda, t.databaseByName, id, resType, rscs, relationships)
		case awsresources.GoogleBQType:
			t.fromLambdaToResource(value, lambda, t.googleBQByName, id, resType, rscs, relationships)
		case awsresources.KinesisType:
			targetARN := awsresources.ParseResourceARN(v, resType)
			t.relationshipsMap[lambdaARN] = append(t.relationshipsMap[lambdaARN], targetARN)
		case awsresources.S3Type:
			targetARN := awsresources.ParseResourceARN(v, resType)
			t.relationshipsMap[lambdaARN] = append(t.relationshipsMap[lambdaARN], targetARN)
		case awsresources.SQSType:
			targetARN := awsresources.ParseResourceARN(v, resType)
			t.relationshipsMap[lambdaARN] = append(t.relationshipsMap[lambdaARN], targetARN)
		case awsresources.RestfulAPIType:
			t.fromLambdaToResource(value, lambda, t.restfulAPIByName, id, resType, rscs, relationships)
		default:
			fmtcolor.Yellow.Printf("yaml to resource: unidentified variable: %s=%s\n", k, v)
		}
	}
}

func (*Transformer) getValueTypeFromEnvar(k string) (value string, resType awsresources.ResourceType) {
	switch {
	case strings.HasSuffix(k, awsresources.EnvarSuffixDBHost):
		value = transformers.ReplaceSuffix(k, awsresources.EnvarSuffixDBHost, awsresources.ToDatabaseCase)
		resType = awsresources.DatabaseType
	case strings.HasSuffix(k, awsresources.EnvarSuffixGoogleBQ):
		value = transformers.ReplaceSuffix(k, awsresources.EnvarSuffixGoogleBQ, awsresources.ToGoogleBQCase)
		resType = awsresources.GoogleBQType
	case strings.HasSuffix(k, awsresources.EnvarSuffixKinesisStreamURL):
		value = transformers.ReplaceSuffix(k, awsresources.EnvarSuffixKinesisStreamURL, awsresources.ToKinesisCase)
		resType = awsresources.KinesisType
	case strings.HasSuffix(k, awsresources.EnvarSuffixS3BucketURL):
		value = transformers.ReplaceSuffix(k, awsresources.EnvarSuffixS3BucketURL, awsresources.ToS3BucketCase)
		resType = awsresources.S3Type
	case strings.HasSuffix(k, awsresources.EnvarSuffixS3BucketName):
		value = transformers.ReplaceSuffix(k, awsresources.EnvarSuffixS3BucketName, awsresources.ToS3BucketCase)
		resType = awsresources.S3Type
	case strings.HasSuffix(k, awsresources.EnvarSuffixSQSQueueURL):
		value = transformers.ReplaceSuffix(k, awsresources.EnvarSuffixSQSQueueURL, awsresources.ToSQSCase)
		resType = awsresources.SQSType
	case strings.HasSuffix(k, awsresources.EnvarSuffixRestfulAPI):
		value = transformers.ReplaceSuffix(k, awsresources.EnvarSuffixRestfulAPI, awsresources.ToRestfulAPICase)
		resType = awsresources.RestfulAPIType
	}

	return value, resType
}

func (t *Transformer) fromLambdaToResource(
	value string, lambda resources.Resource, resourceMap map[string]resources.Resource,
	id *int, resType awsresources.ResourceType, rscs *[]resources.Resource, relationships *[]resources.Relationship,
) {
	r, ok := resourceMap[value]
	if !ok {
		r = resources.NewGenericResource(fmt.Sprintf("%d", *id), value, resType.String())
		*id++

		*rscs = append(*rscs, r)

		resourceMap[value] = r
	}

	*relationships = append(*relationships, resources.Relationship{Source: lambda, Target: r})
}
