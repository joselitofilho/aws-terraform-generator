package resources

import "github.com/ettle/strcase"

const (
	EnvarSuffixDBHost           = "DB_HOST"
	EnvarSuffixGoogleBQ         = "BQ_PROJECT_ID"
	EnvarSuffixKinesisStreamURL = "KINESIS_STREAM_URL"
	EnvarSuffixS3BucketURL      = "S3_BUCKET"
	EnvarSuffixS3BucketName     = "BUCKET_NAME"
	EnvarSuffixSQSQueueURL      = "SQS_QUEUE_URL"
	EnvarSuffixRestfulAPI       = "API_BASE_URL"
)

var (
	ToDatabaseCase   = strcase.ToKebab
	ToGoogleBQCase   = strcase.ToKebab
	ToKinesisCase    = strcase.ToPascal
	ToLambdaCase     = strcase.ToCamel
	ToS3BucketCase   = strcase.ToKebab
	ToSQSCase        = strcase.ToKebab
	ToSNSCase        = strcase.ToKebab
	ToRestfulAPICase = strcase.ToPascal
)

var SuffixByResource = map[ResourceType]string{
	KinesisType: "kinesis",
	S3Type:      "bucket",
	SQSType:     "sqs",
}
