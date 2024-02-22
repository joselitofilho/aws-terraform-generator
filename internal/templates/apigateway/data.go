package apigateway

import "github.com/joselitofilho/aws-terraform-generator/internal/templates"

type Data struct {
	StackName string
	APIDomain string
}

type LambdaData struct {
	ModuleLambdaSource string
	StackName          string
	Name               string
	NameSnakeCase      string
	NamePascalCase     string
	Description        string
	Envars             map[string]string
	Verb               string
	Path               string
	Files              map[string]templates.File
}
