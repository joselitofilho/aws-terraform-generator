package drawio

import (
	"regexp"
	"strings"
)

// Resource interface for all resource types
type Resource interface {
	ID() string
	Value() string
}

// Lambda represents Lambda resource: mxgraph.aws3.lambda
type Lambda struct {
	id    string
	value string
}

func NewLambda(id, name string) Lambda {
	return Lambda{id: id, value: name}
}
func (l Lambda) ID() string    { return l.id }
func (l Lambda) Value() string { return l.value }

// SQS represents SQS resource: mxgraph.aws3.sqs
type SQS struct {
	id    string
	value string
}

func NewSQS(id, name string) SQS {
	return SQS{id: id, value: name}
}
func (s SQS) ID() string    { return s.id }
func (s SQS) Value() string { return s.value }

// Cron represents Cron resource: mxgraph.aws4.event_time_based
type Cron struct {
	id    string
	value string
}

func NewCron(id, name string) Cron {
	return Cron{id: id, value: name}
}
func (c Cron) ID() string    { return c.id }
func (c Cron) Value() string { return c.value }

// APIGateway represents APIGateway resource: mxgraph.aws3.api_gateway
type APIGateway struct {
	id    string
	value string
}

func NewAPIGateway(id, name string) APIGateway {
	return APIGateway{id: id, value: name}
}
func (a APIGateway) ID() string    { return a.id }
func (a APIGateway) Value() string { return a.value }

// Endpoint represents Endpoint resource: mxgraph.aws4.endpoint
type Endpoint struct {
	id    string
	value string
}

func NewEndpoint(id, name string) Endpoint {
	return Endpoint{id: id, value: name}
}
func (e Endpoint) ID() string    { return e.id }
func (e Endpoint) Value() string { return e.value }

// S3 represents S3 resource: mxgraph.aws3.s3
type S3 struct {
	id    string
	value string
}

func NewS3(id, name string) S3 {
	return S3{id: id, value: name}
}
func (s S3) ID() string    { return s.id }
func (s S3) Value() string { return s.value }

// Database represents Database resource: mxgraph.flowchart.database
type Database struct {
	id    string
	value string
}

func NewDatabase(id, name string) Database {
	return Database{id: id, value: name}
}
func (d Database) ID() string    { return d.id }
func (d Database) Value() string { return d.value }

// RestfulAPI represents RestfulAPI resource: mxgraph.veeam2.restful_api
type RestfulAPI struct {
	id    string
	value string
}

func NewRestfulAPI(id, name string) RestfulAPI {
	return RestfulAPI{id: id, value: name}
}
func (r RestfulAPI) ID() string    { return r.id }
func (r RestfulAPI) Value() string { return r.value }

// Relationship struct representing the relationship between resources
type Relationship struct {
	Source Resource
	Target Resource
}

type ResourceCollection struct {
	Lambdas       []Lambda
	SQSs          []SQS
	Crons         []Cron
	APIGateways   []APIGateway
	Endpoints     []Endpoint
	Buckets       []S3
	Databases     []Database
	RestfulAPIs   []RestfulAPI
	Relationships []Relationship
}

func ParseResources(mxFile *MxFile) (*ResourceCollection, error) {
	resources := &ResourceCollection{}

	var reDatabase = regexp.MustCompile(
		`mxgraph.flowchart.database|mxgraph.aws3.dynamo_db|mxgraph.aws4.database|mxgraph.aws4.documentdb_with_mongodb_compatibility`)

	for _, cell := range mxFile.Diagram.MxGraphModel.Root.MxCells {
		switch {
		case strings.Contains(cell.Style, "mxgraph.aws3.lambda"):
			lambda := NewLambda(cell.Id, cell.Value)
			resources.Lambdas = append(resources.Lambdas, lambda)
		case strings.Contains(cell.Style, "mxgraph.aws3.sqs"):
			sqs := NewSQS(cell.Id, cell.Value)
			resources.SQSs = append(resources.SQSs, sqs)
		case strings.Contains(cell.Style, "mxgraph.aws4.event_time_based"):
			cron := NewCron(cell.Id, cell.Value)
			resources.Crons = append(resources.Crons, cron)
		case strings.Contains(cell.Style, "mxgraph.aws3.api_gateway"):
			apiGateway := NewAPIGateway(cell.Id, cell.Value)
			resources.APIGateways = append(resources.APIGateways, apiGateway)
		case strings.Contains(cell.Style, "mxgraph.aws4.endpoint"):
			endpoint := NewEndpoint(cell.Id, cell.Value)
			resources.Endpoints = append(resources.Endpoints, endpoint)
		case strings.Contains(cell.Style, "mxgraph.aws3.s3"):
			storage := NewS3(cell.Id, cell.Value)
			resources.Buckets = append(resources.Buckets, storage)
		case reDatabase.MatchString(cell.Style):
			database := NewDatabase(cell.Id, cell.Value)
			resources.Databases = append(resources.Databases, database)
		case strings.Contains(cell.Style, "mxgraph.veeam2.restful_api"):
			restfulAPI := NewRestfulAPI(cell.Id, cell.Value)
			resources.RestfulAPIs = append(resources.RestfulAPIs, restfulAPI)
		}
	}

	for _, cell := range mxFile.Diagram.MxGraphModel.Root.MxCells {
		if cell.Source != "" && cell.Target != "" {
			source, target := findResourceByID(resources, cell.Source), findResourceByID(resources, cell.Target)
			if source != nil && target != nil {
				relationship := Relationship{Source: source, Target: target}
				resources.Relationships = append(resources.Relationships, relationship)
			}
		}
	}

	return resources, nil
}

func findResourceByID(resources *ResourceCollection, id string) Resource {
	for _, lambda := range resources.Lambdas {
		if lambda.ID() == id {
			return lambda
		}
	}
	for _, sqs := range resources.SQSs {
		if sqs.ID() == id {
			return sqs
		}
	}
	for _, cron := range resources.Crons {
		if cron.ID() == id {
			return cron
		}
	}
	for _, apiGateway := range resources.APIGateways {
		if apiGateway.ID() == id {
			return apiGateway
		}
	}
	for _, endpoint := range resources.Endpoints {
		if endpoint.ID() == id {
			return endpoint
		}
	}
	for _, storage := range resources.Buckets {
		if storage.ID() == id {
			return storage
		}
	}
	for _, database := range resources.Databases {
		if database.ID() == id {
			return database
		}
	}
	for _, restfulAPI := range resources.RestfulAPIs {
		if restfulAPI.ID() == id {
			return restfulAPI
		}
	}
	return nil
}
