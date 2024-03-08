package tftodrawio

import (
	"fmt"
	"slices"
	"strings"

	"github.com/joselitofilho/aws-terraform-generator/internal/generators/terraform"
)

func buildKeyValueVars(tfLocals []*terraform.Local) map[string]string {
	keyValue := map[string]string{}

	for i := range tfLocals {
		for k, v := range tfLocals[i].Attributes {
			varName := fmt.Sprintf("local.%s", k)

			switch v := v.(type) {
			case string:
				keyValue[varName] = v
			case []string:
				buildSliceStringVars(varName, v, keyValue)
			case map[string]any:
				buildStringAnyMapVars(varName, v, keyValue)
			default:
				keyValue[varName] = varName
			}
		}
	}

	return keyValue
}

func buildSliceStringVars(varName string, values []string, keyValue map[string]string) {
	if len(values) > 0 {
		keyValue[varName] = values[0]
	} else {
		keyValue[varName] = varName
	}
}

func buildStringAnyMapVars(varName string, values map[string]any, keyValue map[string]string) {
	arr := make([]string, 0, len(values))
	for k := range values {
		arr = append(arr, k)
	}

	if len(arr) > 0 {
		slices.Sort(arr)
		keyValue[varName] = arr[0]
	} else {
		keyValue[varName] = varName
	}
}

func replaceVars(str string, tfLocals []*terraform.Local) string {
	keyValue := buildKeyValueVars(tfLocals)

	for {
		// TODO: Cycle variables raise a panic here.
		for varName, finalValue := range keyValue {
			str = strings.ReplaceAll(str, varName, finalValue)
		}

		if !strings.Contains(str, "local.") {
			break
		}
	}

	return str
}
