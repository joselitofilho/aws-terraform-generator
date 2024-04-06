package errors

import "errors"

var (
	// ErrDrawIOParser represents a failure in the drawio XML parser.
	ErrDrawIOParser = errors.New("drawio XML parser fails")

	// ErrYAMLParser represents a failure in the YAML parser.
	ErrYAMLParser = errors.New("YAML parser fails")
)
