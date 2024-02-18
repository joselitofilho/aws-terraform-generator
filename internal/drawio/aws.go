package drawio

import (
	"strings"
)

// Resource interface for all resource types
type Resource interface {
	ID() string
	Value() string
	Description() string
}

// Lambda represents Lambda resources: mxgraph.aws3.lambda
type Lambda struct {
	id          string
	value       string
	description string
}

func NewLambda(id, name, description string) Lambda {
	return Lambda{id: id, value: name, description: description}
}
func (l Lambda) ID() string          { return l.id }
func (l Lambda) Value() string       { return l.value }
func (l Lambda) Description() string { return l.description }

// SQS represents SQS resources: mxgraph.aws3.sqs
type SQS struct {
	id          string
	value       string
	description string
}

func NewSQS(id, name, description string) SQS {
	return SQS{id: id, value: name, description: description}
}
func (s SQS) ID() string          { return s.id }
func (s SQS) Value() string       { return s.value }
func (s SQS) Description() string { return s.description }

// Cron represents Cron resources: mxgraph.aws4.event_time_based
type Cron struct {
	id          string
	value       string
	description string
}

func NewCron(id, name, description string) Cron {
	return Cron{id: id, value: name, description: description}
}
func (c Cron) ID() string          { return c.id }
func (c Cron) Value() string       { return c.value }
func (c Cron) Description() string { return c.description }

// APIGateway represents APIGateway resources: mxgraph.aws3.api_gateway
type APIGateway struct {
	id          string
	value       string
	description string
}

func NewAPIGateway(id, name, description string) APIGateway {
	return APIGateway{id: id, value: name, description: description}
}
func (a APIGateway) ID() string          { return a.id }
func (a APIGateway) Value() string       { return a.value }
func (a APIGateway) Description() string { return a.description }

// Endpoint represents Endpoint resources: mxgraph.aws4.endpoint
type Endpoint struct {
	id          string
	value       string
	description string
}

func NewEndpoint(id, name, description string) Endpoint {
	return Endpoint{id: id, value: name, description: description}
}
func (e Endpoint) ID() string          { return e.id }
func (e Endpoint) Value() string       { return e.value }
func (e Endpoint) Description() string { return e.description }

// Database represents Database resources: mxgraph.flowchart.database
type Database struct {
	id          string
	value       string
	description string
}

func NewDatabase(id, name, description string) Database {
	return Database{id: id, value: name, description: description}
}
func (d Database) ID() string          { return d.id }
func (d Database) Value() string       { return d.value }
func (d Database) Description() string { return d.description }

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
	Databases     []Database
	Relationships []Relationship
}

func ParseResources(mxFile *MxFile) (*ResourceCollection, error) {
	resources := &ResourceCollection{}

	for _, cell := range mxFile.Diagram.MxGraphModel.Root.MxCells {
		switch {
		case strings.Contains(cell.Style, "mxgraph.aws3.lambda"):
			lambda := NewLambda(cell.Id, cell.Value, "") // You can parse description if it's available in the diagram
			resources.Lambdas = append(resources.Lambdas, lambda)
		case strings.Contains(cell.Style, "mxgraph.aws3.sqs"):
			sqs := NewSQS(cell.Id, cell.Value, "") // You can parse description if it's available in the diagram
			resources.SQSs = append(resources.SQSs, sqs)
		case strings.Contains(cell.Style, "mxgraph.aws4.event_time_based"):
			cron := NewCron(cell.Id, cell.Value, "") // You can parse description if it's available in the diagram
			resources.Crons = append(resources.Crons, cron)
		case strings.Contains(cell.Style, "mxgraph.aws3.api_gateway"):
			apiGateway := NewAPIGateway(cell.Id, cell.Value, "") // You can parse description if it's available in the diagram
			resources.APIGateways = append(resources.APIGateways, apiGateway)
		case strings.Contains(cell.Style, "mxgraph.aws4.endpoint"):
			endpoint := NewEndpoint(cell.Id, cell.Value, "") // You can parse description if it's available in the diagram
			resources.Endpoints = append(resources.Endpoints, endpoint)
		case strings.Contains(cell.Style, "mxgraph.flowchart.database"):
			database := NewDatabase(cell.Id, cell.Value, "") // You can parse description if it's available in the diagram
			resources.Databases = append(resources.Databases, database)
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
	for _, database := range resources.Databases {
		if database.ID() == id {
			return database
		}
	}
	return nil
}
