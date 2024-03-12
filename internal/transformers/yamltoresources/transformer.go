package yamltoresources

import (
	"errors"
	"fmt"
	"strings"

	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
	"github.com/joselitofilho/aws-terraform-generator/internal/resources"
	"github.com/joselitofilho/aws-terraform-generator/internal/transformers"
)

var ErrEmptyConfig = errors.New("config file is empty")

type Transformer struct {
	yamlConfig *config.Config

	apigatewayByName map[string]resources.Resource
	endpointByName   map[string]struct{}
	kinesisByName    map[string]struct{}
	lambdaByName     map[string]resources.Resource
	restfulAPIByName map[string]struct{}
	s3BucketByName   map[string]struct{}
	sqsByName        map[string]struct{}
	snsByName        map[string]struct{}
}

func NewTransformer(yamlConfig *config.Config) *Transformer {
	return &Transformer{
		yamlConfig: yamlConfig,

		apigatewayByName: map[string]resources.Resource{},
		endpointByName:   map[string]struct{}{},
		kinesisByName:    map[string]struct{}{},
		lambdaByName:     map[string]resources.Resource{},
		restfulAPIByName: map[string]struct{}{},
		s3BucketByName:   map[string]struct{}{},
		sqsByName:        map[string]struct{}{},
		snsByName:        map[string]struct{}{},
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
		if _, ok := t.restfulAPIByName[res.Name]; !ok {
			rscs = append(rscs, resources.NewGenericResource(fmt.Sprintf("%d", id), res.Name, resources.RestfulAPIType))
			id++

			t.restfulAPIByName[res.Name] = struct{}{}
		}
	}

	for _, res := range t.yamlConfig.Buckets {
		if _, ok := t.s3BucketByName[res.Name]; !ok {
			rscs = append(rscs, resources.NewGenericResource(fmt.Sprintf("%d", id), res.Name, resources.S3Type))
			id++

			t.s3BucketByName[res.Name] = struct{}{}
		}
	}

	for _, res := range t.yamlConfig.SQSs {
		if _, ok := t.sqsByName[res.Name]; !ok {
			rscs = append(rscs, resources.NewGenericResource(fmt.Sprintf("%d", id), res.Name, resources.SQSType))
			id++

			t.sqsByName[res.Name] = struct{}{}
		}
	}

	for _, res := range t.yamlConfig.SNSs {
		if _, ok := t.snsByName[res.Name]; !ok {
			rscs = append(rscs, resources.NewGenericResource(fmt.Sprintf("%d", id), res.Name, resources.SNSType))
			id++

			t.snsByName[res.Name] = struct{}{}
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

			t.transformLambda(&config.Lambda{Name: l.Name}, rscs, relationships, id)

			*relationships = append(*relationships, resources.Relationship{
				Source: apigRes,
				Target: t.lambdaByName[l.Name],
			})
		}

		endpointValue := res.APIDomain

		if _, ok := t.endpointByName[endpointValue]; !ok {
			*rscs = append(*rscs,
				resources.NewGenericResource(fmt.Sprintf("%d", *id), endpointValue, resources.EndpointType))
			*id++

			t.endpointByName[endpointValue] = struct{}{}
		}
	}
}

func (t *Transformer) transformKinesis(rscs *[]resources.Resource, id *int) {
	for i := range t.yamlConfig.Kinesis {
		res := t.yamlConfig.Kinesis[i]

		if _, ok := t.kinesisByName[res.Name]; !ok {
			*rscs = append(*rscs, resources.NewGenericResource(fmt.Sprintf("%d", *id), res.Name, resources.KinesisType))
			*id++

			t.kinesisByName[res.Name] = struct{}{}
		}
	}
}

func (t *Transformer) transformLambda(
	res *config.Lambda, rscs *[]resources.Resource, relationships *[]resources.Relationship, id *int,
) {
	if _, ok := t.lambdaByName[res.Name]; ok {
		return
	}

	lambda := resources.NewGenericResource(fmt.Sprintf("%d", *id), res.Name, resources.LambdaType)
	*rscs = append(*rscs, lambda)
	*id++

	t.lambdaByName[res.Name] = lambda

	for _, r := range res.Crons {
		cron := resources.NewGenericResource(fmt.Sprintf("%d", *id), r.ScheduleExpression, resources.CronType)
		*id++

		*relationships = append(*relationships, resources.Relationship{Source: cron, Target: lambda})
	}

	t.transformLambdaEnvars(res, lambda, relationships, id)

	for _, r := range res.KinesisTriggers {
		cron := resources.NewGenericResource(fmt.Sprintf("%d", *id), r.SourceARN, resources.KinesisType)
		*id++

		*relationships = append(*relationships, resources.Relationship{Source: cron, Target: lambda})
	}

	for _, r := range res.SQSTriggers {
		cron := resources.NewGenericResource(fmt.Sprintf("%d", *id), r.SourceARN, resources.SQSType)
		*id++

		*relationships = append(*relationships, resources.Relationship{Source: cron, Target: lambda})
	}
}

func (t *Transformer) transformLambdas(rscs *[]resources.Resource, relationships *[]resources.Relationship, id *int) {
	for i := range t.yamlConfig.Lambdas {
		res := t.yamlConfig.Lambdas[i]

		t.transformLambda(&res, rscs, relationships, id)
	}
}

func (*Transformer) transformLambdaEnvars(
	res *config.Lambda, lambda *resources.GenericResource, relationships *[]resources.Relationship, id *int,
) {
	for _, envars := range res.Envars {
		for k := range envars {
			var (
				value   string
				resType resources.ResourceType
			)

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

			if value != "" {
				r := resources.NewGenericResource(fmt.Sprintf("%d", *id), value, resType)
				*id++

				*relationships = append(*relationships, resources.Relationship{Source: lambda, Target: r})
			}
		}
	}
}
