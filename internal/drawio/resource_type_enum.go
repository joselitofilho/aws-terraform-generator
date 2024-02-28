package drawio

import "strings"

// ResourceType represents the type of a resource
type ResourceType int

const (
	// APIGatewayType represents the API Gateway resource type
	APIGatewayType ResourceType = iota
	// CronType represents the Cron resource type
	CronType
	// DatabaseType represents the Database resource type
	DatabaseType
	// EndpointType represents the Endpoint resource type
	EndpointType
	// KinesisType represents the Kinesis resource type
	KinesisType
	// LambdaType represents the Lambda resource type
	LambdaType
	// RestfulAPIType represents the Restful API resource type
	RestfulAPIType
	// S3Type represents the S3 resource type
	S3Type
	// SNSType represents the SNS resource type
	SNSType
	// SQSType represents the SQS resource type
	SQSType
	// UnknownType represents an unknown resource type
	UnknownType
)

// String returns the string representation of a ResourceType
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

// ParseResourceType parses a ResourceType from a string
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
