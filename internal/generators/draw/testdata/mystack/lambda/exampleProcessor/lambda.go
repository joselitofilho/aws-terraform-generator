package main

import (
	"context"
)

type exampleProcessorLambda struct{}

func newExampleProcessorLambda() *exampleProcessorLambda {
	return &exampleProcessorLambda{}
}

func (l *exampleProcessorLambda) run(ctx context.Context) error {
	// TODO: Implement

	return nil
}
