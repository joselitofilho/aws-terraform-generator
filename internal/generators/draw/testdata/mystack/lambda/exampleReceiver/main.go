// Package main lambda body
package main

import (
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	exampleReceiverLambda := newExampleReceiverLambda()

	lambda.Start(exampleReceiverLambda.run)
}
