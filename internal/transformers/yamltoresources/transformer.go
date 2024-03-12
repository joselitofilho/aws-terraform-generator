package yamltoresources

import (
	"errors"
	"fmt"
	"strings"

	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
	"github.com/joselitofilho/aws-terraform-generator/internal/resources"
	"github.com/joselitofilho/aws-terraform-generator/internal/transformers"
	"github.com/joselitofilho/aws-terraform-generator/internal/transformers/terraformtoresources"
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
	sqsByName        map[string]resources.Resource
	snsByName        map[string]resources.Resource
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
		sqsByName:        map[string]resources.Resource{},
		snsByName:        map[string]resources.Resource{},
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

	t.transformKinesis(&rscs, &id)

	t.transformLambdas(&rscs, &relationships, &id)

	for _, res := range t.yamlConfig.RestfulAPIs {
		name := terraformtoresources.ResourceByARN(res.Name).Name
		name = resources.ToRestfulAPICase(name)

		if _, ok := t.restfulAPIByName[name]; !ok {
			restfulAPIRes := resources.NewGenericResource(fmt.Sprintf("%d", id), name, resources.RestfulAPIType)
			id++

			rscs = append(rscs, restfulAPIRes)

			t.restfulAPIByName[name] = restfulAPIRes
		}
	}

	for _, res := range t.yamlConfig.Buckets {
		name := terraformtoresources.ResourceByARN(res.Name).Name
		name = resources.ToS3BucketCase(name)

		if _, ok := t.s3BucketByName[name]; !ok {
			s3BucketResource := resources.NewGenericResource(fmt.Sprintf("%d", id), name, resources.S3Type)
			id++

			rscs = append(rscs, s3BucketResource)

			t.s3BucketByName[name] = s3BucketResource
		}
	}

	for _, res := range t.yamlConfig.SQSs {
		name := terraformtoresources.ResourceByARN(res.Name).Name
		name = resources.ToSQSCase(name)

		if _, ok := t.sqsByName[name]; !ok {
			sqsResource := resources.NewGenericResource(fmt.Sprintf("%d", id), name, resources.SQSType)
			id++

			rscs = append(rscs, sqsResource)

			t.sqsByName[name] = sqsResource
		}
	}

	for _, res := range t.yamlConfig.SNSs {
		name := terraformtoresources.ResourceByARN(res.Name).Name
		name = resources.ToSNSCase(name)

		if _, ok := t.snsByName[name]; !ok {
			snsResource := resources.NewGenericResource(fmt.Sprintf("%d", id), name, resources.SNSType)
			id++

			rscs = append(rscs, snsResource)

			t.snsByName[name] = snsResource
		}
	}

	return &resources.ResourceCollection{
		Resources:     rscs,
		Relationships: relationships,
	}, nil
}

func (t *Transformer) transformAPIGateways(
	rscs *[]resources.Resource, relationships *[]resources.Relationship, id *int,
) {
	for _, res := range t.yamlConfig.APIGateways {
		endpointValue := res.APIDomain

		endpointRes, ok := t.endpointByName[endpointValue]
		if !ok {
			endpointRes = resources.NewGenericResource(fmt.Sprintf("%d", *id), endpointValue, resources.EndpointType)
			*rscs = append(*rscs, endpointRes)
			*id++

			t.endpointByName[endpointValue] = endpointRes
		}

		for i := range res.Lambdas {
			l := res.Lambdas[i]
			apigValue := fmt.Sprintf("%s %s", l.Verb, l.Path)

			apigRes, ok := t.apigatewayByName[apigValue]

			if !ok {
				apigRes = resources.NewGenericResource(fmt.Sprintf("%d", *id), apigValue, resources.APIGatewayType)
				*rscs = append(*rscs, apigRes)
				*id++

				t.apigatewayByName[apigValue] = apigRes
			}

			lambdaName := resources.ToLambdaCase(l.Name)

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

func (t *Transformer) transformKinesis(rscs *[]resources.Resource, id *int) {
	for i := range t.yamlConfig.Kinesis {
		res := t.yamlConfig.Kinesis[i]
		name := terraformtoresources.ResourceByARN(res.Name).Name
		name = resources.ToKinesisCase(name)

		if _, ok := t.kinesisByName[name]; !ok {
			kinesisResource := resources.NewGenericResource(fmt.Sprintf("%d", *id), name, resources.KinesisType)
			*id++

			*rscs = append(*rscs, kinesisResource)

			t.kinesisByName[name] = kinesisResource
		}
	}
}

func (t *Transformer) transformLambda(
	res *config.Lambda, rscs *[]resources.Resource, relationships *[]resources.Relationship, id *int,
) {
	if _, ok := t.lambdaByName[res.Name]; ok {
		return
	}

	lambdaName := terraformtoresources.ResourceByARN(res.Name).Name
	lambdaName = resources.ToLambdaCase(lambdaName)

	lambda, ok := t.lambdaByName[lambdaName]
	if !ok {
		lambda = resources.NewGenericResource(fmt.Sprintf("%d", *id), lambdaName, resources.LambdaType)
		*id++

		*rscs = append(*rscs, lambda)

		t.lambdaByName[lambdaName] = lambda
	}

	for _, r := range res.Crons {
		name := r.ScheduleExpression

		sourceRes, ok := t.cronByName[name]
		if !ok {
			sourceRes = resources.NewGenericResource(fmt.Sprintf("%d", *id), name, resources.CronType)
			*id++

			*rscs = append(*rscs, sourceRes)

			t.cronByName[name] = sourceRes
		}

		*relationships = append(*relationships, resources.Relationship{Source: sourceRes, Target: lambda})
	}

	t.transformLambdaEnvars(res, lambda, rscs, relationships, id)

	for _, r := range res.KinesisTriggers {
		name := terraformtoresources.ResourceByARN(r.SourceARN).Name
		name = resources.ToKinesisCase(name)

		sourceRes, ok := t.kinesisByName[name]
		if !ok {
			sourceRes = resources.NewGenericResource(fmt.Sprintf("%d", *id), name, resources.KinesisType)
			*id++

			*rscs = append(*rscs, sourceRes)

			t.kinesisByName[name] = sourceRes
		}

		*relationships = append(*relationships, resources.Relationship{Source: sourceRes, Target: lambda})
	}

	for _, r := range res.SQSTriggers {
		name := terraformtoresources.ResourceByARN(r.SourceARN).Name
		name = resources.ToSQSCase(name)

		sourceRes, ok := t.sqsByName[name]
		if !ok {
			sourceRes = resources.NewGenericResource(fmt.Sprintf("%d", *id), name, resources.SQSType)
			*id++

			*rscs = append(*rscs, sourceRes)

			t.sqsByName[name] = sourceRes
		}

		*relationships = append(*relationships, resources.Relationship{Source: sourceRes, Target: lambda})
	}
}

func (t *Transformer) transformLambdas(rscs *[]resources.Resource, relationships *[]resources.Relationship, id *int) {
	for i := range t.yamlConfig.Lambdas {
		res := t.yamlConfig.Lambdas[i]

		t.transformLambda(&res, rscs, relationships, id)
	}
}

func (t *Transformer) transformLambdaEnvars(
	res *config.Lambda, lambda resources.Resource,
	rscs *[]resources.Resource, relationships *[]resources.Relationship, id *int,
) {
	for _, envars := range res.Envars {
		for k := range envars {
			value, resType := t.getValueTypeFromEnvar(k)

			if value != "" {
				switch resType {
				case resources.DatabaseType:
					t.fromLambdaToResource(value, lambda, t.databaseByName, id, resType, rscs, relationships)
				case resources.GoogleBQType:
					t.fromLambdaToResource(value, lambda, t.googleBQByName, id, resType, rscs, relationships)
				case resources.KinesisType:
					t.fromLambdaToResource(value, lambda, t.kinesisByName, id, resType, rscs, relationships)
				case resources.S3Type:
					t.fromLambdaToResource(value, lambda, t.s3BucketByName, id, resType, rscs, relationships)
				case resources.SQSType:
					t.fromLambdaToResource(value, lambda, t.sqsByName, id, resType, rscs, relationships)
				case resources.RestfulAPIType:
					t.fromLambdaToResource(value, lambda, t.restfulAPIByName, id, resType, rscs, relationships)
				default:
					r := resources.NewGenericResource(fmt.Sprintf("%d", *id), value, resType)
					*id++

					*relationships = append(*relationships, resources.Relationship{Source: lambda, Target: r})
				}
			}
		}
	}
}

func (*Transformer) getValueTypeFromEnvar(k string) (value string, resType resources.ResourceType) {
	switch {
	case strings.HasSuffix(k, resources.EnvarSuffixDBHost):
		value = transformers.ReplaceSuffix(k, resources.EnvarSuffixDBHost, resources.ToDatabaseCase)
		resType = resources.DatabaseType
	case strings.HasSuffix(k, resources.EnvarSuffixGoogleBQ):
		value = transformers.ReplaceSuffix(k, resources.EnvarSuffixGoogleBQ, resources.ToGoogleBQCase)
		resType = resources.GoogleBQType
	case strings.HasSuffix(k, resources.EnvarSuffixKinesisStreamURL):
		value = transformers.ReplaceSuffix(k, resources.EnvarSuffixKinesisStreamURL, resources.ToKinesisCase)
		resType = resources.KinesisType
	case strings.HasSuffix(k, resources.EnvarSuffixS3BucketURL):
		value = transformers.ReplaceSuffix(k, resources.EnvarSuffixS3BucketURL, resources.ToS3BucketCase)
		resType = resources.S3Type
	case strings.HasSuffix(k, resources.EnvarSuffixS3BucketName):
		value = transformers.ReplaceSuffix(k, resources.EnvarSuffixS3BucketName, resources.ToS3BucketCase)
		resType = resources.S3Type
	case strings.HasSuffix(k, resources.EnvarSuffixSQSQueueURL):
		value = transformers.ReplaceSuffix(k, resources.EnvarSuffixSQSQueueURL, resources.ToSQSCase)
		resType = resources.SQSType
	case strings.HasSuffix(k, resources.EnvarSuffixRestfulAPI):
		value = transformers.ReplaceSuffix(k, resources.EnvarSuffixRestfulAPI, resources.ToRestfulAPICase)
		resType = resources.RestfulAPIType
	}

	return value, resType
}

func (t *Transformer) fromLambdaToResource(
	value string, lambda resources.Resource, resourceMap map[string]resources.Resource,
	id *int, resType resources.ResourceType, rscs *[]resources.Resource, relationships *[]resources.Relationship,
) {
	r, ok := resourceMap[value]
	if !ok {
		r = resources.NewGenericResource(fmt.Sprintf("%d", *id), value, resType)
		*id++

		*rscs = append(*rscs, r)

		resourceMap[value] = r
	}

	*relationships = append(*relationships, resources.Relationship{Source: lambda, Target: r})
}
