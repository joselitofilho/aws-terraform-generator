package resources

import "strings"

// ResourceType represents the type of a resource.
type ResourceType string

const (
	// APIGatewayType represents the API Gateway resource type.
	APIGatewayType ResourceType = "apigateway"

	// CronType represents the Cron resource type.
	CronType ResourceType = "cron"

	// DatabaseType represents the Database resource type.
	DatabaseType ResourceType = "database"

	// EndpointType represents the Endpoint resource type.
	EndpointType ResourceType = "endpoint"

	// GoogleBQType represents the Google BigQuery resource type.
	GoogleBQType ResourceType = "googlebq"

	// KinesisType represents the Kinesis resource type.
	KinesisType ResourceType = "kinesis"

	// LambdaType represents the Lambda resource type.
	LambdaType ResourceType = "lambda"

	// RestfulAPIType represents the Restful API resource type.
	RestfulAPIType ResourceType = "restfulapi"

	// S3Type represents the S3 resource type.
	S3Type ResourceType = "s3"

	// SNSType represents the SNS resource type.
	SNSType ResourceType = "sns"

	// SQSType represents the SQS resource type.
	SQSType ResourceType = "sqs"

	// UnknownType represents an unknown resource type.
	UnknownType ResourceType = "unknown"
)

// String returns the string representation of a ResourceType.
func (rt ResourceType) String() string {
	switch rt {
	case APIGatewayType:
		return "APIGateway"
	case CronType:
		return "Cron"
	case DatabaseType:
		return "Database"
	case EndpointType:
		return "Endpoint"
	case GoogleBQType:
		return "GoogleBQ"
	case KinesisType:
		return "Kinesis"
	case LambdaType:
		return "Lambda"
	case RestfulAPIType:
		return "RestfulAPI"
	case S3Type:
		return "S3"
	case SNSType:
		return "SNS"
	case SQSType:
		return "SQS"
	default:
		return "Unknown"
	}
}

// ParseResourceType parses a ResourceType from a string.
func ParseResourceType(s string) ResourceType {
	switch strings.ToLower(s) {
	case "apigateway":
		return APIGatewayType
	case "cron":
		return CronType
	case "database":
		return DatabaseType
	case "endpoint":
		return EndpointType
	case "googlebq":
		return GoogleBQType
	case "kinesis":
		return KinesisType
	case "lambda":
		return LambdaType
	case "restfulapi":
		return RestfulAPIType
	case "s3":
		return S3Type
	case "sns":
		return SNSType
	case "sqs":
		return SQSType
	default:
		return UnknownType
	}
}
