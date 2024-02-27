package transformers

import (
	"fmt"
	"strings"

	"github.com/ettle/strcase"
	"github.com/joselitofilho/aws-terraform-generator/internal/drawio"
	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
)

func buildCronToLambda(cronsByLambdaID map[string]drawio.Resource, cron, lambda drawio.Resource) {
	cronsByLambdaID[lambda.ID()] = cron
}

func buildEndpointToAPIGateway(
	apiGatewaysByID map[string]drawio.Resource,
	endpointsByAPIGatewayID map[string]drawio.Resource,
	endpoint drawio.Resource,
	apiGateways drawio.Resource,
) {
	apiGatewayID := apiGateways.ID()
	apiGatewaysByID[apiGatewayID] = apiGateways
	endpointsByAPIGatewayID[apiGatewayID] = endpoint
}

func buildLambdaVars(envars map[string]map[string]string, lambda, target drawio.Resource, variables []string) {
	targetName := initLambdaEnvarsAndGetTargetName(envars, lambda, target)

	for _, v := range variables {
		envars[lambda.ID()][fmt.Sprintf("%s%s",
			strcase.ToSNAKE(targetName), v)] = strcase.ToSnake(fmt.Sprintf("var.%s%s", targetName, v))
	}
}

func buildLambdaToDatabase(envars map[string]map[string]string, lambda, database drawio.Resource) {
	buildLambdaVars(envars, lambda, database, []string{"DB_HOST", "DB_USER", "DB_PASSWORD_SECRET"})
}

func buildLambdaToRestfulAPI(envars map[string]map[string]string, lambda, restfulAPI drawio.Resource) {
	buildLambdaVars(envars, lambda, restfulAPI, []string{"_API_BASE_URL", "_HOST", "_USER"})
}

func buildLambdaToS3(envars map[string]map[string]string, lambda, s3Bucket drawio.Resource) {
	bucketName := initLambdaEnvarsAndGetTargetName(envars, lambda, s3Bucket)

	envars[lambda.ID()][fmt.Sprintf("%s_S3_BUCKET",
		strcase.ToSNAKE(bucketName))] = fmt.Sprintf("aws_s3_bucket.%s_bucket.bucket", strcase.ToSnake(bucketName))
	envars[lambda.ID()][fmt.Sprintf("%s_S3_DIRECTORY",
		strcase.ToSNAKE(bucketName))] = fmt.Sprintf("%q_files", strings.ToLower(strcase.ToSnake(lambda.Value())))
}

func buildLambdaToSQS(envars map[string]map[string]string, lambda, sqs drawio.Resource,
) {
	sqsName := initLambdaEnvarsAndGetTargetName(envars, lambda, sqs)

	envars[lambda.ID()][fmt.Sprintf("%s_SQS_QUEUE_URL",
		strcase.ToSNAKE(sqsName))] = fmt.Sprintf("aws_sqs_queue.%s_sqs.id", strcase.ToSnake(sqsName))
}

func buildSNSToLambda(snsMap map[string]config.SNS, sns drawio.Resource) {
	snsConfig, ok := snsMap[sns.ID()]
	if !ok {
		snsConfig = config.SNS{Name: sns.Value()}
	}

	snsConfig.Lambdas = append(snsConfig.Lambdas, config.SNSResource{Name: sns.Value()})

	snsMap[sns.ID()] = snsConfig
}

func buildSNSToSQS(snsMap map[string]config.SNS, sqs drawio.Resource) {
	snsConfig, ok := snsMap[sqs.ID()]
	if !ok {
		snsConfig = config.SNS{Name: sqs.Value()}
	}

	snsConfig.SQSs = append(snsConfig.SQSs, config.SNSResource{Name: sqs.Value()})

	snsMap[sqs.ID()] = snsConfig
}

func buildSQSToLambda(sqsTriggersByLambdaID map[string][]drawio.Resource, sqs, lambda drawio.Resource) {
	lambdaID := lambda.ID()
	sqsTriggersByLambdaID[lambdaID] = append(sqsTriggersByLambdaID[lambdaID], sqs)
}

func buildS3ToSNS(snsMap map[string]config.SNS, s3Bucket, sns drawio.Resource) {
	snsConfig, ok := snsMap[sns.ID()]
	if !ok {
		snsConfig = config.SNS{Name: sns.Value()}
	}

	snsConfig.BucketName = s3Bucket.Value()

	snsMap[sns.ID()] = snsConfig
}

func initEnvarsIfNecessaryByKey(envars map[string]map[string]string, key string) {
	if _, ok := envars[key]; !ok {
		envars[key] = map[string]string{}
	}
}

func initLambdaEnvarsAndGetTargetName(envars map[string]map[string]string, lambda, target drawio.Resource) string {
	initEnvarsIfNecessaryByKey(envars, lambda.ID())
	return strings.ToLower(target.Value())
}
