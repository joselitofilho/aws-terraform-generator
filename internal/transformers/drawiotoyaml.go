package transformers

import (
	"fmt"
	"strings"

	"github.com/ettle/strcase"
	"github.com/joselitofilho/aws-terraform-generator/internal/drawio"
	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
)

func TransformDrawIOToYAML(yamlConfig *config.Config, resources *drawio.ResourceCollection) (*config.Config, error) {
	endpointsByAPIGatewayID := map[string]drawio.Resource{}
	apiGatewaysByID := map[string]drawio.Resource{}
	cronsByLambdaID := map[string]drawio.Resource{}
	sqsTriggersByLambdaID := map[string][]drawio.Resource{}
	envars := map[string]map[string]string{}
	snsMap := map[string]config.SNS{}

	initEnvarsIfNecessaryByKey := func(target map[string]map[string]string, key string) {
		if _, ok := target[key]; !ok {
			target[key] = map[string]string{}
		}
	}

	resourcesByTypeMap := map[drawio.ResourceType][]drawio.Resource{}
	for _, resource := range resources.Resources {
		resourcesByTypeMap[resource.ReseourceType()] = append(resourcesByTypeMap[resource.ReseourceType()], resource)
	}

	for _, rel := range resources.Relationships {
		switch rel.Target.ReseourceType() {
		case drawio.LambdaType:
			lambdaID := rel.Target.ID()

			switch rel.Source.ReseourceType() {
			case drawio.CronType:
				cronsByLambdaID[lambdaID] = rel.Source
			case drawio.SQSType:
				sqsTriggersByLambdaID[lambdaID] = append(sqsTriggersByLambdaID[lambdaID], rel.Source)
			case drawio.SNSType:
				sns, ok := snsMap[rel.Source.ID()]
				if !ok {
					sns = config.SNS{Name: rel.Source.Value()}
				}

				sns.Lambdas = append(sns.Lambdas, config.SNSResource{Name: rel.Source.Value()})

				snsMap[rel.Source.ID()] = sns
			}
		case drawio.APIGatewayType:
			apiGatewayID := rel.Target.ID()
			apiGatewaysByID[apiGatewayID] = rel.Target
			endpointsByAPIGatewayID[apiGatewayID] = rel.Source
		case drawio.SQSType:
			switch rel.Source.ReseourceType() {
			case drawio.LambdaType:
				initEnvarsIfNecessaryByKey(envars, rel.Source.ID())

				envars[rel.Source.ID()]["SQS_QUEUE_URL"] =
					fmt.Sprintf("aws_sqs_queue.%s_sqs.id", strcase.ToSnake(rel.Target.Value()))
			case drawio.SNSType:
				sns, ok := snsMap[rel.Source.ID()]
				if !ok {
					sns = config.SNS{Name: rel.Source.Value()}
				}

				sns.SQSs = append(sns.SQSs, config.SNSResource{Name: rel.Source.Value()})

				snsMap[rel.Source.ID()] = sns
			}
		case drawio.DatabaseType:
			switch rel.Source.ReseourceType() {
			case drawio.LambdaType:
				initEnvarsIfNecessaryByKey(envars, rel.Source.ID())

				envars[rel.Source.ID()]["DOCDB_HOST"] = "var.docdb_host"
				envars[rel.Source.ID()]["DOCDB_USER"] = "var.docdb_user"
				envars[rel.Source.ID()]["DOCDB_PASSWORD_SECRET"] = "var.docdb_password_secret"
			}
		case drawio.RestfulAPIType:
			switch rel.Source.ReseourceType() {
			case drawio.LambdaType:
				initEnvarsIfNecessaryByKey(envars, rel.Source.ID())

				restfulAPIName := strings.ToLower(rel.Target.Value())

				envars[rel.Source.ID()][fmt.Sprintf("%s_API_BASE_URL", strcase.ToSNAKE(restfulAPIName))] =
					fmt.Sprintf("var.%s_api_base_url", strcase.ToSnake(restfulAPIName))
				envars[rel.Source.ID()][fmt.Sprintf("%s_HOST", strcase.ToSNAKE(restfulAPIName))] =
					fmt.Sprintf("var.%s_host", strcase.ToSnake(restfulAPIName))
				envars[rel.Source.ID()][fmt.Sprintf("%s_USER", strcase.ToSNAKE(restfulAPIName))] =
					fmt.Sprintf("var.%s_user", strcase.ToSnake(restfulAPIName))
			}
		case drawio.S3Type:
			switch rel.Source.ReseourceType() {
			case drawio.LambdaType:
				initEnvarsIfNecessaryByKey(envars, rel.Source.ID())

				bucketName := strings.ToLower(rel.Target.Value())

				envars[rel.Source.ID()][fmt.Sprintf("%s_S3_BUCKET", strcase.ToSNAKE(bucketName))] =
					fmt.Sprintf("aws_s3_bucket.%s_bucket.bucket", strcase.ToSnake(bucketName))
				envars[rel.Source.ID()][fmt.Sprintf("%s_S3_DIRECTORY", strcase.ToSNAKE(bucketName))] =
					fmt.Sprintf("%s_files", strings.ToLower(strcase.ToSnake(rel.Target.Value())))
			}
		case drawio.SNSType:
			switch rel.Source.ReseourceType() {
			case drawio.S3Type:
				sns, ok := snsMap[rel.Target.ID()]
				if !ok {
					sns = config.SNS{Name: rel.Target.Value()}
				}

				sns.BucketName = rel.Source.Value()

				snsMap[rel.Target.ID()] = sns
			}
		}
	}

	defaultFiles := []config.File{{Name: "lambda.go"}, {Name: "main.go"}}
	lambdas := []config.Lambda{}
	apiGatewayLambdas := map[string][]config.APIGatewayLambda{}

	for _, lambda := range resourcesByTypeMap[drawio.LambdaType] {
		isAPIGatewayLambda := false

		for _, rel := range resources.Relationships {
			if rel.Target.ID() == lambda.ID() {
				switch rel.Source.ReseourceType() {
				case drawio.APIGatewayType:
					isAPIGatewayLambda = true

					apiGatewayID := rel.Source.ID()

					envarsList := []map[string]string{}
					for key, value := range envars[lambda.ID()] {
						envarsList = append(envarsList, map[string]string{key: value})
					}

					apiGatewayLambdas[apiGatewayID] = append(apiGatewayLambdas[apiGatewayID], config.APIGatewayLambda{
						Source:      yamlConfig.Diagram.Modules.Lambda,
						Name:        lambda.Value(),
						Description: fmt.Sprintf("%s lambda", lambda.Value()),
						Envars:      envarsList,
						Verb:        strings.Split(rel.Source.Value(), " ")[0],
						Path:        strings.Split(rel.Source.Value(), " ")[1],
						Files:       defaultFiles,
					})
				}
			}

			if isAPIGatewayLambda {
				break
			}
		}

		if !isAPIGatewayLambda {
			crons := []config.Cron{}
			if cron, ok := cronsByLambdaID[lambda.ID()]; ok {
				crons = append(crons, config.Cron{
					ScheduleExpression: cron.Value(),
					IsEnabled:          "true",
				})
			}

			envarsList := []map[string]string{}
			for key, value := range envars[lambda.ID()] {
				envarsList = append(envarsList, map[string]string{key: value})
			}

			sqsTriggers := []config.SQSTrigger{}
			for _, sqsTrigger := range sqsTriggersByLambdaID[lambda.ID()] {
				sqsTriggers = append(sqsTriggers, config.SQSTrigger{
					SourceARN: fmt.Sprintf("aws_sqs_queue.%s_sqs.arn", strcase.ToSnake(sqsTrigger.Value())),
				})
			}

			lambdas = append(lambdas, config.Lambda{
				Source:      yamlConfig.Diagram.Modules.Lambda,
				Name:        lambda.Value(),
				Description: fmt.Sprintf("%s lambda", lambda.Value()),
				Envars:      envarsList,
				SQSTriggers: sqsTriggers,
				Files:       defaultFiles,
				Crons:       crons,
			})
		}
	}

	apiGateways := []config.APIGateway{
		{
			StackName: yamlConfig.Diagram.StackName,
			APIG:      true,
		},
	}

	for id := range apiGatewaysByID {
		apiGateways[0].APIDomain = endpointsByAPIGatewayID[id].Value()
		apiGateways[0].Lambdas = append(apiGateways[0].Lambdas, apiGatewayLambdas[id]...)
	}

	sqss := []config.SQS{}

	for _, sqs := range resourcesByTypeMap[drawio.SQSType] {
		sqss = append(sqss, config.SQS{Name: sqs.Value(), MaxReceiveCount: 10})
	}

	buckets := []config.S3{}

	for _, bucket := range resourcesByTypeMap[drawio.S3Type] {
		buckets = append(buckets, config.S3{Name: bucket.Value()})
	}

	restfulAPIs := []config.RestfulAPI{}
	restfulAPINames := map[string]struct{}{}

	for _, restfulAPI := range resourcesByTypeMap[drawio.RestfulAPIType] {
		name := restfulAPI.Value()
		if _, ok := restfulAPINames[name]; !ok {
			restfulAPIs = append(restfulAPIs, config.RestfulAPI{Name: name})
			restfulAPINames[name] = struct{}{}
		}
	}

	snss := make([]config.SNS, 0, len(snsMap))
	for _, sns := range snsMap {
		snss = append(snss, sns)
	}

	return &config.Config{
		Lambdas:     lambdas,
		APIGateways: apiGateways,
		SQSs:        sqss,
		Buckets:     buckets,
		RestfulAPIs: restfulAPIs,
		SNSs:        snss,
	}, nil
}
