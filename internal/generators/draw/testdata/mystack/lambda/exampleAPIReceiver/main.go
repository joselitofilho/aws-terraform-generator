// Package main lambda body
package main

import (
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	exampleAPIReceiverLambda := newExampleApiReceiverLambda()

	lambda.Start(exampleAPIReceiverLambda.run)
}
