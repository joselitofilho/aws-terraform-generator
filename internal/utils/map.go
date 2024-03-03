package utils

func MergeStringMap(left, right map[string]string) map[string]string {
	result := left
	for k, v := range right {
		result[k] = v
	}

	return result
}
