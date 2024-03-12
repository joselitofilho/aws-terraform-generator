package transformers

import "strings"

func ReplaceSuffix(value, suffix string, fn func(s string) string) string {
	value = strings.ReplaceAll(value, "_"+suffix, "")
	value = strings.ReplaceAll(value, suffix, "")

	if fn != nil {
		value = fn(value)
	}

	return value
}
