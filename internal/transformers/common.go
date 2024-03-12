package transformers

import "strings"

func ReplaceSuffix(value, suffix string, fn func(s string) string) string {
	if strings.HasSuffix(value, "_"+suffix) {
		value = strings.TrimSuffix(value, "_"+suffix)
	} else {
		value = strings.TrimSuffix(value, suffix)
	}

	if fn != nil {
		value = fn(value)
	}

	return value
}
