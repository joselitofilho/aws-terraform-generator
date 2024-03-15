package resources

import (
	"fmt"
	"strings"
)

var (
	arnKinesisKey  = "kinesis"
	arnLambdaKey   = "lambda"
	arnS3BucketKey = "s3"
	arnSQSKey      = "sqs"
	arnSNSKey      = "sns"
)

var arnKeySuffix = map[string]string{
	arnKinesisKey:  "stream",
	arnLambdaKey:   "function",
	arnS3BucketKey: "bucket",
	arnSQSKey:      "queue",
	arnSNSKey:      "topic",
}

var labelByResourceType = map[ResourceType]string{
	APIGatewayType: LabelAWSAPIGatewayRoute,
	CronType:       LabelAWSCron,
	EndpointType:   LabelAWSEndpoint,
	KinesisType:    LabelAWSKinesisStream,
	LambdaType:     LabelAWSLambdaFunction,
	S3Type:         LabelAWSS3Bucket,
	SQSType:        LabelAWSSQSQueue,
	SNSType:        LabelAWSSNSTopic,
	UnknownType:    "",
}

type ResourceARN struct {
	Type  string
	Name  string
	Label string
}

func (r *ResourceARN) LabelOrName() (value string) {
	value = r.Label
	if value == "" {
		value = r.Name
	}

	return value
}

func ParseResourceARN(arn string, suggestedResType ResourceType) ResourceARN {
	var arnType, name, label string

	switch {
	case strings.HasPrefix(arn, "arn:"):
		arnType, name = parseColumnARNTypeAndName(arn)
	case strings.HasPrefix(arn, "http"):
		arnType, name = parseHTTPARNTypeAndName(arn)
	default:
		arnType, name, label = parseARNTypeAndName(arn)
	}

	if suggestedResType == UnknownType {
		suggestedResType = inferResourceType(arnType)
	}

	if arnType == "" {
		arnType = labelByResourceType[suggestedResType]
	}

	return ResourceARN{Type: arnType, Name: name, Label: label}
}

func parseColumnARNTypeAndName(arn string) (arnType, name string) {
	parts := strings.Split(arn, ":")
	arnType = fmt.Sprintf("aws_%s_%s", parts[2], arnKeySuffix[parts[2]])

	if arnType == LabelAWSKinesisStream {
		parts = strings.Split(arn, "/")
	}

	name = parts[len(parts)-1]

	return arnType, name
}

func parseHTTPARNTypeAndName(arn string) (arnType, name string) {
	parts := strings.Split(arn, "//")
	parts = strings.Split(parts[1], "/")

	resStrType := strings.Split(parts[0], ".")[0]
	arnType = fmt.Sprintf("aws_%s_%s", resStrType, arnKeySuffix[resStrType])

	name = parts[len(parts)-1]

	return arnType, name
}

func parseARNTypeAndName(arn string) (arnType, name, label string) {
	parts := strings.Split(arn, ".")

	if len(parts) > 0 && parts[0] == "module" {
		// TODO: Implement other modules.
		arnType = LabelAWSLambdaFunction
		label = parts[1]
	} else if len(parts) > 1 && strings.HasPrefix(parts[0], "aws_") {
		arnType = parts[0]
		label = parts[1]
	} else {
		name = arn
	}

	return arnType, name, label
}

func inferResourceType(arnType string) ResourceType {
	switch arnType {
	case LabelAWSKinesisStream:
		return KinesisType
	case LabelAWSLambdaFunction:
		return LambdaType
	case LabelAWSS3Bucket:
		return S3Type
	case LabelAWSSNSTopic:
		return SNSType
	case LabelAWSSQSQueue:
		return SQSType
	default:
		return UnknownType
	}
}
