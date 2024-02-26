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

func buildLambdaToDatabase(envars map[string]map[string]string, lambda, database drawio.Resource) {
	initEnvarsIfNecessaryByKey(envars, lambda.ID())

	dbName := strings.ToLower(database.Value())

	envars[lambda.ID()][fmt.Sprintf("%sDB_HOST",
		strcase.ToSNAKE(dbName))] = fmt.Sprintf("var.%sdb_host", strcase.ToSnake(dbName))
	envars[lambda.ID()][fmt.Sprintf("%sDB_USER",
		strcase.ToSNAKE(dbName))] = fmt.Sprintf("var.%sdb_user", strcase.ToSnake(dbName))
	envars[lambda.ID()][fmt.Sprintf("%sDB_PASSWORD_SECRET",
		strcase.ToSNAKE(dbName))] = fmt.Sprintf("var.%sdb_password_secret", strcase.ToSnake(dbName))
}

func buildLambdaToRestfulAPI(envars map[string]map[string]string, lambda, restfulAPI drawio.Resource) {
	initEnvarsIfNecessaryByKey(envars, lambda.ID())

	restfulAPIName := strings.ToLower(restfulAPI.Value())

	envars[lambda.ID()][fmt.Sprintf("%s_API_BASE_URL",
		strcase.ToSNAKE(restfulAPIName))] = fmt.Sprintf("var.%s_api_base_url", strcase.ToSnake(restfulAPIName))
	envars[lambda.ID()][fmt.Sprintf("%s_HOST",
		strcase.ToSNAKE(restfulAPIName))] = fmt.Sprintf("var.%s_host", strcase.ToSnake(restfulAPIName))
	envars[lambda.ID()][fmt.Sprintf("%s_USER",
		strcase.ToSNAKE(restfulAPIName))] = fmt.Sprintf("var.%s_user", strcase.ToSnake(restfulAPIName))
}

func buildLambdaToS3(envars map[string]map[string]string, lambda, s3Bucket drawio.Resource) {
	initEnvarsIfNecessaryByKey(envars, lambda.ID())

	bucketName := strings.ToLower(s3Bucket.Value())

	envars[lambda.ID()][fmt.Sprintf("%s_S3_BUCKET",
		strcase.ToSNAKE(bucketName))] = fmt.Sprintf("aws_s3_bucket.%s_bucket.bucket", strcase.ToSnake(bucketName))
	envars[lambda.ID()][fmt.Sprintf("%s_S3_DIRECTORY",
		strcase.ToSNAKE(bucketName))] = fmt.Sprintf("%q_files", strings.ToLower(strcase.ToSnake(lambda.Value())))
}

func buildLambdaToSQS(envars map[string]map[string]string, lambda, sqs drawio.Resource,
) {
	initEnvarsIfNecessaryByKey(envars, lambda.ID())

	sqsName := strings.ToLower(sqs.Value())

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

func initEnvarsIfNecessaryByKey(target map[string]map[string]string, key string) {
	if _, ok := target[key]; !ok {
		target[key] = map[string]string{}
	}
}
