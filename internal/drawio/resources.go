package drawio

import (
	"regexp"
	"strings"
)

// Resource interface for all resource types.
type Resource interface {
	ID() string
	Value() string
	ResourceType() ResourceType
}

// GenericResource represents a generic resource.
type GenericResource struct {
	id           string
	value        string
	resourceType ResourceType
}

func NewGenericResource(id, value string, resourceType ResourceType) *GenericResource {
	return &GenericResource{id: id, value: value, resourceType: resourceType}
}
func (r *GenericResource) ID() string                 { return r.id }
func (r *GenericResource) Value() string              { return r.value }
func (r *GenericResource) ResourceType() ResourceType { return r.resourceType }

// Relationship struct representing the relationship between resources.
type Relationship struct {
	Source Resource
	Target Resource
}

// ResourceCollection represents a collection of resources and their relationships.
type ResourceCollection struct {
	Resources     []Resource
	Relationships []Relationship
}

// NewResourceCollection creates a new ResourceCollection.
func NewResourceCollection() *ResourceCollection {
	return &ResourceCollection{}
}

// AddResource adds a resource to the collection.
func (rc *ResourceCollection) AddResource(resource Resource) {
	rc.Resources = append(rc.Resources, resource)
}

// AddRelationship adds a relationship to the collection.
func (rc *ResourceCollection) AddRelationship(source, target Resource) {
	relationship := Relationship{Source: source, Target: target}
	rc.Relationships = append(rc.Relationships, relationship)
}

// ParseResources parses resources from the MxFile.
func ParseResources(mxFile *MxFile) (*ResourceCollection, error) {
	resources := NewResourceCollection()

	for i := range mxFile.Diagram.MxGraphModel.Root.MxCells {
		cell := mxFile.Diagram.MxGraphModel.Root.MxCells[i]

		resource := createResource(cell.ID, cell.Value, cell.Style)
		if resource != nil {
			resources.AddResource(resource)
		}
	}

	resourcesMap := map[string]Resource{}
	for _, resource := range resources.Resources {
		resourcesMap[resource.ID()] = resource
	}

	for i := range mxFile.Diagram.MxGraphModel.Root.MxCells {
		cell := mxFile.Diagram.MxGraphModel.Root.MxCells[i]
		if cell.Source != "" && cell.Target != "" {
			source := resourcesMap[cell.Source]
			target := resourcesMap[cell.Target]

			if source != nil && target != nil {
				resources.AddRelationship(source, target)
			}
		}
	}

	return resources, nil
}

// createResource creates a resource based on cell data.
func createResource(id, value, style string) Resource {
	reAPIGateway := regexp.MustCompile("mxgraph.aws3.api_gateway|mxgraph.aws4.api_gateway")
	reDatabase := regexp.MustCompile(`mxgraph.flowchart.database|mxgraph.aws3.dynamo_db|mxgraph.aws4.database|` +
		`mxgraph.aws4.documentdb_with_mongodb_compatibility`)
	reKinesis := regexp.MustCompile(`mxgraph.aws3.kinesis|mxgraph.aws4.kinesis`)
	resLambda := regexp.MustCompile(`mxgraph.aws3.lambda|mxgraph.aws4.lambda`)
	resRestfulAPI := regexp.MustCompile(`mxgraph.veeam2.restful_api|mxgraph.veeam.2d.restful_apis`)
	reS3 := regexp.MustCompile(`mxgraph.aws3.s3|mxgraph.aws4.s3`)
	reSQS := regexp.MustCompile(`mxgraph.aws3.sqs|mxgraph.aws4.sqs`)
	reSNS := regexp.MustCompile(`mxgraph.aws3.sns|mxgraph.aws4.sns`)

	switch {
	case reAPIGateway.MatchString(style):
		return &GenericResource{id: id, value: value, resourceType: APIGatewayType}
	case strings.Contains(style, "mxgraph.aws4.event_time_based"):
		return &GenericResource{id: id, value: value, resourceType: CronType}
	case reDatabase.MatchString(style):
		return &GenericResource{id: id, value: value, resourceType: DatabaseType}
	case strings.Contains(style, "mxgraph.aws4.endpoint"):
		return &GenericResource{id: id, value: value, resourceType: EndpointType}
	case reKinesis.MatchString(style):
		return &GenericResource{id: id, value: value, resourceType: KinesisType}
	case resLambda.MatchString(style):
		return &GenericResource{id: id, value: value, resourceType: LambdaType}
	case resRestfulAPI.MatchString(style):
		return &GenericResource{id: id, value: value, resourceType: RestfulAPIType}
	case reS3.MatchString(style):
		return &GenericResource{id: id, value: value, resourceType: S3Type}
	case reSQS.MatchString(style):
		return &GenericResource{id: id, value: value, resourceType: SQSType}
	case reSNS.MatchString(style):
		return &GenericResource{id: id, value: value, resourceType: SNSType}
	default:
		return nil
	}
}
