package drawiotoyaml

import (
	"fmt"
	"strings"

	"github.com/ettle/strcase"

	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
	"github.com/joselitofilho/aws-terraform-generator/internal/resources"
)

func buildCronToLambda(cronsByLambdaID map[string]resources.Resource, cron, lambda resources.Resource) {
	cronsByLambdaID[lambda.ID()] = cron
}

func buildEndpointToAPIGateway(
	apiGatewaysByID map[string]resources.Resource,
	endpointsByAPIGatewayID map[string]resources.Resource,
	endpoint resources.Resource,
	apiGateways resources.Resource,
) {
	apiGatewayID := apiGateways.ID()
	apiGatewaysByID[apiGatewayID] = apiGateways
	endpointsByAPIGatewayID[apiGatewayID] = endpoint
}

func buildKinesisToLambda(
	kinesisTriggersByLambdaID map[string][]resources.Resource, kinesis, lambda resources.Resource,
) {
	lambdaID := lambda.ID()
	kinesisTriggersByLambdaID[lambdaID] = append(kinesisTriggersByLambdaID[lambdaID], kinesis)
}

func buildLambdaVars(lambda, target resources.Resource, variables []string, envars map[string]map[string]string) {
	targetName := initLambdaEnvarsAndGetTargetName(lambda, target, envars)

	for _, v := range variables {
		envars[lambda.ID()][strcase.ToSNAKE(fmt.Sprintf("%s%s",
			targetName, v))] = "var." + strcase.ToSnake(fmt.Sprintf("%s%s", targetName, v))
	}
}

func buildLambdaToDatabase(lambda, database resources.Resource, envars map[string]map[string]string) {
	buildLambdaVars(lambda, database, []string{"DB_HOST", "DB_USER", "DB_PASSWORD_SECRET"}, envars)
}

func buildLambdaToGoogleBQ(lambda, googleBQ resources.Resource, envars map[string]map[string]string) {
	buildLambdaVars(lambda, googleBQ,
		[]string{"BQ_PROJECT_ID", "BQ_API_KEY_SECRET", "BQ_PARTITION_FIELD", "BQ_CLUSTERING_FIELDS"}, envars)
}

func buildLambdaToKinesis(lambda, kinesis resources.Resource, envars map[string]map[string]string) {
	kinesisName := initLambdaEnvarsAndGetTargetName(lambda, kinesis, envars)

	envars[lambda.ID()][fmt.Sprintf("%s_KINESIS_STREAM_URL",
		strcase.ToSNAKE(kinesisName))] = fmt.Sprintf("aws_kinesis_stream.%s_kinesis.name", strcase.ToSnake(kinesisName))
}

func buildLambdaToRestfulAPI(lambda, restfulAPI resources.Resource, envars map[string]map[string]string) {
	buildLambdaVars(lambda, restfulAPI, []string{"API_BASE_URL", "HOST", "USER"}, envars)
}

func buildLambdaToS3(lambda, s3Bucket resources.Resource, envars map[string]map[string]string) {
	bucketName := initLambdaEnvarsAndGetTargetName(lambda, s3Bucket, envars)

	envars[lambda.ID()][fmt.Sprintf("%s_S3_BUCKET",
		strcase.ToSNAKE(bucketName))] = fmt.Sprintf("aws_s3_bucket.%s_bucket.bucket", strcase.ToSnake(bucketName))
	envars[lambda.ID()][fmt.Sprintf("%s_S3_DIRECTORY",
		strcase.ToSNAKE(bucketName))] = fmt.Sprintf(`"%s_files"`, strings.ToLower(strcase.ToSnake(lambda.Value())))
}

func buildLambdaToSQS(lambda, sqs resources.Resource, envars map[string]map[string]string) {
	sqsName := initLambdaEnvarsAndGetTargetName(lambda, sqs, envars)

	envars[lambda.ID()][fmt.Sprintf("%s_SQS_QUEUE_URL",
		strcase.ToSNAKE(sqsName))] = fmt.Sprintf("aws_sqs_queue.%s_sqs.name", strcase.ToSnake(sqsName))
}

func buildS3ToSNS(snsMap map[string]config.SNS, s3Bucket, sns resources.Resource) {
	snsConfig := snsMap[sns.ID()]
	snsConfig.BucketName = s3Bucket.Value()
	snsMap[sns.ID()] = snsConfig
}

func buildSNSToLambda(snsMap map[string]config.SNS, sns, lambda resources.Resource) {
	snsConfig := snsMap[sns.ID()]
	snsConfig.Lambdas = append(snsConfig.Lambdas, config.SNSResource{
		Name:   lambda.Value(),
		Events: []string{"s3:ObjectCreated:*"},
	})
	snsMap[sns.ID()] = snsConfig
}

func buildSNSToSQS(snsMap map[string]config.SNS, sns, sqs resources.Resource) {
	snsConfig := snsMap[sns.ID()]
	snsConfig.SQSs = append(snsConfig.SQSs, config.SNSResource{
		Name:   sqs.Value(),
		Events: []string{"s3:ObjectCreated:*"},
	})
	snsMap[sns.ID()] = snsConfig
}

func buildSQSToLambda(sqsTriggersByLambdaID map[string][]resources.Resource, sqs, lambda resources.Resource) {
	lambdaID := lambda.ID()
	sqsTriggersByLambdaID[lambdaID] = append(sqsTriggersByLambdaID[lambdaID], sqs)
}

func initEnvarsIfNecessaryByKey(key string, envars map[string]map[string]string) {
	if _, ok := envars[key]; !ok {
		envars[key] = map[string]string{}
	}
}

func initLambdaEnvarsAndGetTargetName(lambda, target resources.Resource, envars map[string]map[string]string) string {
	initEnvarsIfNecessaryByKey(lambda.ID(), envars)
	return target.Value()
}
