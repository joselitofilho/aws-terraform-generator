package drawiotoresources

import (
	"regexp"
	"strings"

	"github.com/joselitofilho/aws-terraform-generator/internal/drawio"
	"github.com/joselitofilho/aws-terraform-generator/internal/resources"
)

// Transform parses resources from the MxFile.
func Transform(mxFile *drawio.MxFile) (*resources.ResourceCollection, error) {
	resc := resources.NewResourceCollection()

	for i := range mxFile.Diagram.MxGraphModel.Root.MxCells {
		cell := mxFile.Diagram.MxGraphModel.Root.MxCells[i]

		resource := createResource(cell.ID, cell.Value, cell.Style)
		if resource != nil {
			resc.AddResource(resource)
		}
	}

	resourcesMap := map[string]resources.Resource{}
	for _, resource := range resc.Resources {
		resourcesMap[resource.ID()] = resource
	}

	for i := range mxFile.Diagram.MxGraphModel.Root.MxCells {
		cell := mxFile.Diagram.MxGraphModel.Root.MxCells[i]
		if cell.Source != "" && cell.Target != "" {
			source := resourcesMap[cell.Source]
			target := resourcesMap[cell.Target]

			if source != nil && target != nil {
				resc.AddRelationship(source, target)
			}
		}
	}

	return resc, nil
}

// createResource creates a resource based on cell data.
func createResource(id, value, style string) resources.Resource {
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
		return resources.NewGenericResource(id, value, resources.APIGatewayType)
	case strings.Contains(style, "mxgraph.aws4.event_time_based"):
		return resources.NewGenericResource(id, value, resources.CronType)
	case reDatabase.MatchString(style):
		return resources.NewGenericResource(id, value, resources.DatabaseType)
	case strings.Contains(style, "mxgraph.aws4.endpoint"):
		return resources.NewGenericResource(id, value, resources.EndpointType)
	case reGoogleBQ.MatchString(style):
		return resources.NewGenericResource(id, value, resources.GoogleBQType)
	case reKinesis.MatchString(style):
		return resources.NewGenericResource(id, value, resources.KinesisType)
	case resLambda.MatchString(style):
		return resources.NewGenericResource(id, value, resources.LambdaType)
	case resRestfulAPI.MatchString(style):
		return resources.NewGenericResource(id, value, resources.RestfulAPIType)
	case reS3.MatchString(style):
		return resources.NewGenericResource(id, value, resources.S3Type)
	case reSQS.MatchString(style):
		return resources.NewGenericResource(id, value, resources.SQSType)
	case reSNS.MatchString(style):
		return resources.NewGenericResource(id, value, resources.SNSType)
	default:
		return nil
	}
}
