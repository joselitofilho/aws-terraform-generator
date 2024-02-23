package drawio

import "strings"

// ResourceType represents the type of a resource
type ResourceType int

const (
	// LambdaType represents the Lambda resource type
	LambdaType ResourceType = iota
	// SQSType represents the SQS resource type
	SQSType
	// SNSType represents the SNS resource type
	SNSType
	// CronType represents the Cron resource type
	CronType
	// APIGatewayType represents the API Gateway resource type
	APIGatewayType
	// EndpointType represents the Endpoint resource type
	EndpointType
	// S3Type represents the S3 resource type
	S3Type
	// DatabaseType represents the Database resource type
	DatabaseType
	// RestfulAPIType represents the Restful API resource type
	RestfulAPIType
	// UnknownType represents an unknown resource type
	UnknownType
)

// String returns the string representation of a ResourceType
func (rt ResourceType) String() string {
	switch rt {
	case LambdaType:
		return "Lambda"
	case SQSType:
		return "SQS"
	case SNSType:
		return "SNS"
	case CronType:
		return "Cron"
	case APIGatewayType:
		return "APIGateway"
	case EndpointType:
		return "Endpoint"
	case S3Type:
		return "S3"
	case DatabaseType:
		return "Database"
	case RestfulAPIType:
		return "RestfulAPI"
	default:
		return "Unknown"
	}
}

// ParseResourceType parses a ResourceType from a string
func ParseResourceType(s string) ResourceType {
	switch strings.ToLower(s) {
	case "lambda":
		return LambdaType
	case "sqs":
		return SQSType
	case "sns":
		return SNSType
	case "cron":
		return CronType
	case "apigateway":
		return APIGatewayType
	case "endpoint":
		return EndpointType
	case "s3":
		return S3Type
	case "database":
		return DatabaseType
	case "restfulapi":
		return RestfulAPIType
	default:
		return UnknownType
	}
}
