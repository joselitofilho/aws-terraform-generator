package resources_to_yaml

import (
	"fmt"
	"strings"

	"github.com/ettle/strcase"

	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
	"github.com/joselitofilho/aws-terraform-generator/internal/resources"
)

func (t *Transformer) buildCronToLambda(cron, lambda resources.Resource) {
	t.cronsByLambdaID[lambda.ID()] = cron
}

func (t *Transformer) buildEndpointToAPIGateway(endpoint, apiGateway resources.Resource) {
	apiGatewayID := apiGateway.ID()
	t.apiGatewaysByID[apiGatewayID] = apiGateway
	t.endpointsByAPIGatewayID[apiGatewayID] = endpoint
}

func (t *Transformer) buildKinesisToLambda(kinesis, lambda resources.Resource,
) {
	lambdaID := lambda.ID()
	t.kinesisTriggersByLambdaID[lambdaID] = append(t.kinesisTriggersByLambdaID[lambdaID], kinesis)
}

func (t *Transformer) buildLambdaVars(lambda, target resources.Resource, variables []string) {
	targetName := t.initLambdaEnvarsAndGetTargetName(lambda, target)

	for _, v := range variables {
		t.envars[lambda.ID()][strcase.ToSNAKE(fmt.Sprintf("%s%s",
			targetName, v))] = "var." + strcase.ToSnake(fmt.Sprintf("%s%s", targetName, v))
	}
}

func (t *Transformer) buildLambdaToDatabase(lambda, database resources.Resource) {
	t.buildLambdaVars(lambda, database, []string{"DB_HOST", "DB_USER", "DB_PASSWORD_SECRET"})
}

func (t *Transformer) buildLambdaToGoogleBQ(lambda, googleBQ resources.Resource) {
	t.buildLambdaVars(lambda, googleBQ,
		[]string{"BQ_PROJECT_ID", "BQ_API_KEY_SECRET", "BQ_PARTITION_FIELD", "BQ_CLUSTERING_FIELDS"})
}

func (t *Transformer) buildLambdaToKinesis(lambda, kinesis resources.Resource) {
	kinesisName := t.initLambdaEnvarsAndGetTargetName(lambda, kinesis)

	t.envars[lambda.ID()][fmt.Sprintf("%s_KINESIS_STREAM_URL",
		strcase.ToSNAKE(kinesisName))] = fmt.Sprintf("aws_kinesis_stream.%s_kinesis.name", strcase.ToSnake(kinesisName))
}

func (t *Transformer) buildLambdaToRestfulAPI(lambda, restfulAPI resources.Resource) {
	t.buildLambdaVars(lambda, restfulAPI, []string{"API_BASE_URL", "HOST", "USER"})
}

func (t *Transformer) buildLambdaToS3(lambda, s3Bucket resources.Resource) {
	bucketName := t.initLambdaEnvarsAndGetTargetName(lambda, s3Bucket)

	t.envars[lambda.ID()][fmt.Sprintf("%s_S3_BUCKET",
		strcase.ToSNAKE(bucketName))] = fmt.Sprintf("aws_s3_bucket.%s_bucket.bucket", strcase.ToSnake(bucketName))
	t.envars[lambda.ID()][fmt.Sprintf("%s_S3_DIRECTORY",
		strcase.ToSNAKE(bucketName))] = fmt.Sprintf(`"%s_files"`, strings.ToLower(strcase.ToSnake(lambda.Value())))
}

func (t *Transformer) buildLambdaToSQS(lambda, sqs resources.Resource) {
	sqsName := t.initLambdaEnvarsAndGetTargetName(lambda, sqs)

	t.envars[lambda.ID()][fmt.Sprintf("%s_SQS_QUEUE_URL",
		strcase.ToSNAKE(sqsName))] = fmt.Sprintf("aws_sqs_queue.%s_sqs.name", strcase.ToSnake(sqsName))
}

func (t *Transformer) buildS3ToSNS(s3Bucket, sns resources.Resource) {
	snsConfig := t.snsMap[sns.ID()]
	snsConfig.BucketName = s3Bucket.Value()
	t.snsMap[sns.ID()] = snsConfig
}

func (t *Transformer) buildSNSToLambda(sns, lambda resources.Resource) {
	snsConfig := t.snsMap[sns.ID()]
	snsConfig.Lambdas = append(snsConfig.Lambdas, config.SNSResource{
		Name:   lambda.Value(),
		Events: []string{"s3:ObjectCreated:*"},
	})
	t.snsMap[sns.ID()] = snsConfig
}

func (t *Transformer) buildSNSToSQS(sns, sqs resources.Resource) {
	snsConfig := t.snsMap[sns.ID()]
	snsConfig.SQSs = append(snsConfig.SQSs, config.SNSResource{
		Name:   sqs.Value(),
		Events: []string{"s3:ObjectCreated:*"},
	})
	t.snsMap[sns.ID()] = snsConfig
}

func (t *Transformer) buildSQSToLambda(sqs, lambda resources.Resource) {
	lambdaID := lambda.ID()
	t.sqsTriggersByLambdaID[lambdaID] = append(t.sqsTriggersByLambdaID[lambdaID], sqs)
}

func (t *Transformer) initEnvarsIfNecessaryByKey(key string) {
	if _, ok := t.envars[key]; !ok {
		t.envars[key] = map[string]string{}
	}
}

func (t *Transformer) initLambdaEnvarsAndGetTargetName(lambda, target resources.Resource) string {
	t.initEnvarsIfNecessaryByKey(lambda.ID())
	return target.Value()
}
