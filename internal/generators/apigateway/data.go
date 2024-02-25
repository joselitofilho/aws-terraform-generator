package apigateway

import "github.com/joselitofilho/aws-terraform-generator/internal/generators"

type Data struct {
	StackName string
	APIDomain string
}

type LambdaData struct {
	Name        string
	AsModule    bool
	Source      string
	RoleName    string
	Runtime     string
	StackName   string
	Description string
	Envars      map[string]string
	Verb        string
	Path        string
	Files       map[string]generators.File
}
