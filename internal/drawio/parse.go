package drawio

import (
	"encoding/xml"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/joselitofilho/aws-terraform-generator/internal/resources"
)

// MxFile represents the root element of the draw.io XML file.
type MxFile struct {
	XMLName xml.Name `xml:"mxfile"`
	Diagram Diagram  `xml:"diagram"`
}

// Diagram represents the diagram element within the draw.io XML file.
type Diagram struct {
	XMLName      xml.Name     `xml:"diagram"`
	MxGraphModel MxGraphModel `xml:"mxGraphModel"`
}

// MxGraphModel represents the graph model element within the draw.io XML file.
type MxGraphModel struct {
	XMLName xml.Name `xml:"mxGraphModel"`
	Root    Root     `xml:"root"`
}

// Root represents the root element within the graph model of the draw.io XML file.
type Root struct {
	XMLName xml.Name `xml:"root"`
	MxCells []MxCell `xml:"mxCell"`
}

// MxCell represents a cell element within the draw.io XML file.
type MxCell struct {
	XMLName  xml.Name `xml:"mxCell"`
	ID       string   `xml:"id,attr"`
	Value    string   `xml:"value,attr"`
	Style    string   `xml:"style,attr"`
	Parent   string   `xml:"parent,attr"`
	Vertex   bool     `xml:"vertex,attr"`
	Source   string   `xml:"source,attr"`
	Target   string   `xml:"target,attr"`
	Geometry Geometry `xml:"mxGeometry"`
}

// Geometry represents the geometry element within a cell of the draw.io XML file.
type Geometry struct {
	XMLName xml.Name `xml:"mxGeometry"`
	X       float64  `xml:"x,attr"`
	Y       float64  `xml:"y,attr"`
	Width   float64  `xml:"width,attr"`
	Height  float64  `xml:"height,attr"`
}

// ParseXML parses a draw.io XML file and returns an MxFile struct.
func ParseXML(fileName string) (*MxFile, error) {
	xmlFile, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer xmlFile.Close()

	var mxFile MxFile
	if err := xml.NewDecoder(xmlFile).Decode(&mxFile); err != nil {
		return nil, fmt.Errorf("error decoding XML: %w", err)
	}

	return &mxFile, nil
}

// ParseResources parses resources from the MxFile.
func ParseResources(mxFile *MxFile) (*resources.ResourceCollection, error) {
	rscs := resources.NewResourceCollection()

	for i := range mxFile.Diagram.MxGraphModel.Root.MxCells {
		cell := mxFile.Diagram.MxGraphModel.Root.MxCells[i]

		resource := createResource(cell.ID, cell.Value, cell.Style)
		if resource != nil {
			rscs.AddResource(resource)
		}
	}

	resourcesMap := map[string]resources.Resource{}
	for _, resource := range rscs.Resources {
		resourcesMap[resource.ID()] = resource
	}

	for i := range mxFile.Diagram.MxGraphModel.Root.MxCells {
		cell := mxFile.Diagram.MxGraphModel.Root.MxCells[i]
		if cell.Source != "" && cell.Target != "" {
			source := resourcesMap[cell.Source]
			target := resourcesMap[cell.Target]

			if source != nil && target != nil {
				rscs.AddRelationship(source, target)
			}
		}
	}

	return rscs, nil
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
