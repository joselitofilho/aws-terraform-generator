package resources

import (
	"regexp"
	"strings"

	"github.com/diagram-code-generator/resources/pkg/resources"
)

type AWSResourceFactory struct{}

// CreateResource creates a resource based on cell data.
func (f *AWSResourceFactory) CreateResource(id, value, style string) resources.Resource {
	reAPIGateway := regexp.MustCompile("mxgraph.aws3.api_gateway|mxgraph.aws4.api_gateway")
	reDatabase := regexp.MustCompile(`mxgraph.flowchart.database|mxgraph.aws3.dynamo_db|mxgraph.aws4.database|` +
		`mxgraph.aws4.documentdb_with_mongodb_compatibility`)
	reGoogleBQ := regexp.MustCompile("mxgraph.gcp2.big_query|google_bigquery")
	reKinesis := regexp.MustCompile(`mxgraph.aws3.kinesis|mxgraph.aws4.kinesis`)
	resLambda := regexp.MustCompile(`mxgraph.aws3.lambda|mxgraph.aws4.lambda`)
	resRestfulAPI := regexp.MustCompile(`mxgraph.veeam2.restful_api|mxgraph.veeam.2d.restful_apis`)
	reS3 := regexp.MustCompile(`mxgraph.aws3.s3|mxgraph.aws4.s3`)
	reSQS := regexp.MustCompile(`mxgraph.aws3.sqs|mxgraph.aws4.sqs`)
	reSNS := regexp.MustCompile(`mxgraph.aws3.sns|mxgraph.aws4.sns`)

	switch {
	case reAPIGateway.MatchString(style):
		return resources.NewGenericResource(id, value, APIGatewayType.String())
	case strings.Contains(style, "mxgraph.aws4.event_time_based"):
		return resources.NewGenericResource(id, value, CronType.String())
	case reDatabase.MatchString(style):
		return resources.NewGenericResource(id, value, DatabaseType.String())
	case strings.Contains(style, "mxgraph.aws4.endpoint"):
		return resources.NewGenericResource(id, value, EndpointType.String())
	case reGoogleBQ.MatchString(style):
		return resources.NewGenericResource(id, value, GoogleBQType.String())
	case reKinesis.MatchString(style):
		return resources.NewGenericResource(id, value, KinesisType.String())
	case resLambda.MatchString(style):
		return resources.NewGenericResource(id, value, LambdaType.String())
	case resRestfulAPI.MatchString(style):
		return resources.NewGenericResource(id, value, RestfulAPIType.String())
	case reS3.MatchString(style):
		return resources.NewGenericResource(id, value, S3Type.String())
	case reSQS.MatchString(style):
		return resources.NewGenericResource(id, value, SQSType.String())
	case reSNS.MatchString(style):
		return resources.NewGenericResource(id, value, SNSType.String())
	default:
		return nil
	}
}
