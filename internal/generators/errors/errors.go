package errors

import "errors"

var (
	// ErrDrawIOParser represents a failure in the drawio XML parser.
	ErrDrawIOParser = errors.New("drawio XML parser fails")

	// ErrDrawIOToResourcesTransformer represents a failure in the drawio to resources transformer.
	ErrDrawIOToResourcesTransformer = errors.New("drawio to resources transformer fails")

	// ErrYAMLParser represents a failure in the YAML parser.
	ErrYAMLParser = errors.New("YAML parser fails")
)
