package main

import (
	"context"
)

type exampleReceiverLambda struct{}

func newExampleReceiverLambda() *exampleReceiverLambda {
	return &exampleReceiverLambda{}
}

func (l *exampleReceiverLambda) run(ctx context.Context) error {
	// TODO: Implement

	return nil
}
