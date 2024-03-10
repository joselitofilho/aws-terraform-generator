// Package main lambda body
package main

import (
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	exampleCheckerLambda := newExampleCheckerLambda()

	lambda.Start(exampleCheckerLambda.run)
}
