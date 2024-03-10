package main

import (
	"context"
)

type exampleCheckerLambda struct{}

func newExampleCheckerLambda() *exampleCheckerLambda {
	return &exampleCheckerLambda{}
}

func (l *exampleCheckerLambda) run(ctx context.Context) error {
	// TODO: Implement

	return nil
}
