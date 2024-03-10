// Package main lambda body
package main

import (
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	exampleProcessorLambda := newExampleProcessorLambda()

	lambda.Start(exampleProcessorLambda.run)
}
