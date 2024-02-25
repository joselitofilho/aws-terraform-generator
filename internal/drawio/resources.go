package drawio

import (
	"fmt"
	"regexp"
	"strings"
)

// Resource interface for all resource types
type Resource interface {
	ID() string
	Value() string
	ReseourceType() ResourceType
}

// GenericResource represents a generic resource
type GenericResource struct {
	id           string
	value        string
	resourceType ResourceType
}

func (r *GenericResource) ID() string                  { return r.id }
func (r *GenericResource) Value() string               { return r.value }
func (r *GenericResource) ReseourceType() ResourceType { return r.resourceType }

// Relationship struct representing the relationship between resources
type Relationship struct {
	Source Resource
	Target Resource
}

// ResourceCollection represents a collection of resources and their relationships
type ResourceCollection struct {
	Resources     []Resource
	Relationships []Relationship
}

// NewResourceCollection creates a new ResourceCollection
func NewResourceCollection() *ResourceCollection {
	return &ResourceCollection{}
}

// AddResource adds a resource to the collection
func (rc *ResourceCollection) AddResource(resource Resource) {
	rc.Resources = append(rc.Resources, resource)
}

// AddRelationship adds a relationship to the collection
func (rc *ResourceCollection) AddRelationship(source, target Resource) {
	relationship := Relationship{Source: source, Target: target}
	rc.Relationships = append(rc.Relationships, relationship)
}

// ParseResources parses resources from the MxFile
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

// createResource creates a resource based on cell data
func createResource(id, value, style string) Resource {
	reDatabase := regexp.MustCompile(`mxgraph.flowchart.database|mxgraph.aws3.dynamo_db|mxgraph.aws4.database|` +
		`mxgraph.aws4.documentdb_with_mongodb_compatibility`)

	switch {
	case strings.Contains(style, "mxgraph.aws3.lambda"):
		return &GenericResource{id: id, value: value, resourceType: LambdaType}
	case strings.Contains(style, "mxgraph.aws3.sqs"):
		return &GenericResource{id: id, value: value, resourceType: SQSType}
	case strings.Contains(style, "mxgraph.aws3.sns"):
		return &GenericResource{id: id, value: value, resourceType: SNSType}
	case strings.Contains(style, "mxgraph.aws4.event_time_based"):
		return &GenericResource{id: id, value: value, resourceType: CronType}
	case strings.Contains(style, "mxgraph.aws3.api_gateway"):
		return &GenericResource{id: id, value: value, resourceType: APIGatewayType}
	case strings.Contains(style, "mxgraph.aws4.endpoint"):
		return &GenericResource{id: id, value: value, resourceType: EndpointType}
	case strings.Contains(style, "mxgraph.aws3.s3"):
		return &GenericResource{id: id, value: value, resourceType: S3Type}
	case reDatabase.MatchString(style):
		return &GenericResource{id: id, value: value, resourceType: DatabaseType}
	case strings.Contains(style, "mxgraph.veeam2.restful_api"):
		return &GenericResource{id: id, value: value, resourceType: RestfulAPIType}
	default:
		if strings.Contains(style, "shape") {
			fmt.Printf("%+v\n", []string{id, value, style})
		}
		return nil
	}
}
