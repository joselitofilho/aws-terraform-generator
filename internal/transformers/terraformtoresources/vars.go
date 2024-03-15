package terraformtoresources

import (
	"fmt"
	"slices"
	"strings"

	"github.com/joselitofilho/aws-terraform-generator/internal/generators/terraform"
)

func buildKeyValueFromLocals(tfLocals []*terraform.Local) map[string]string {
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

func buildKeyValueFromVariables(tfVariables []*terraform.Variable) map[string]string {
	keyValue := map[string]string{}

	for i := range tfVariables {
		for k, v := range tfVariables[i].Attributes {
			varName := fmt.Sprintf("var.%s", k)

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

func replaceVars(
	str string, tfVars []*terraform.Variable, tfLocals []*terraform.Local, replaceableStrs map[string]string,
) string {
	keyValue := buildKeyValueFromVariables(tfVars)
	str = replaceVariables(str, keyValue)

	keyValue = buildKeyValueFromLocals(tfLocals)
	str = replaceLocals(str, keyValue)

	str = replaceStrings(str, replaceableStrs)

	return str
}

func replaceVariables(str string, keyValue map[string]string) string {
	for i := 0; i <= len(keyValue); i++ {
		for varName, finalValue := range keyValue {
			str = strings.ReplaceAll(str, varName, finalValue)
		}

		if !strings.Contains(str, "var.") {
			break
		}
	}

	return str
}

func replaceLocals(str string, keyValue map[string]string) string {
	for i := 0; i <= len(keyValue); i++ {
		for varName, finalValue := range keyValue {
			str = strings.ReplaceAll(str, varName, finalValue)
		}

		if !strings.Contains(str, "local.") {
			break
		}
	}

	return str
}

func replaceStrings(str string, replaceableStrs map[string]string) string {
	for i := 0; i <= len(replaceableStrs); i++ {
		for varName, finalValue := range replaceableStrs {
			str = strings.ReplaceAll(str, varName, finalValue)
		}

		for varName := range replaceableStrs {
			if !strings.Contains(str, varName) {
				break
			}
		}
	}

	return str
}
