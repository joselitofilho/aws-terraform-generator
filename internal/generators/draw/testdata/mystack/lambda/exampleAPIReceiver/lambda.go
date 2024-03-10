package main

import (
	"context"
)

type exampleAPIReceiverLambda struct{}

func newExampleApiReceiverLambda() *exampleAPIReceiverLambda {
	return &exampleAPIReceiverLambda{}
}

func (l *exampleAPIReceiverLambda) run(ctx context.Context) error {
	// TODO: Implement

	return nil
}
