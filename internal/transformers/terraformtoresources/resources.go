package terraformtoresources

import (
	"regexp"
	"strings"

	"github.com/ettle/strcase"

	"github.com/joselitofilho/aws-terraform-generator/internal/fmtcolor"
	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
	"github.com/joselitofilho/aws-terraform-generator/internal/resources"
)

type resourceARN struct {
	key  string
	name string
}

func (t *Transformer) hasResourceMatched(res resources.Resource, filters config.Filters) bool {
	if res == nil {
		return false
	}

	filter, hasFilter := filters[res.ResourceType()]
	if !hasFilter {
		return true
	}

	match := len(filter.Match) == 0

	for _, pattern := range filter.Match {
		regex, err := regexp.Compile(pattern)
		if err != nil {
			fmtcolor.Yellow.Println("error compiling match regex:", err)
			continue
		}

		if regex.MatchString(res.Value()) {
			match = true
			break
		}
	}

	for _, pattern := range filter.NotMatch {
		regex, err := regexp.Compile(pattern)
		if err != nil {
			fmtcolor.Yellow.Println("error compiling not_match regex:", err)
			continue
		}

		if regex.MatchString(res.Value()) {
			match = false
			break
		}
	}

	return match
}

func resourceByARN(arn string) resourceARN {
	var key, name string

	if strings.HasPrefix(arn, "arn:") {
		parts := strings.Split(arn, ":")

		key = parts[2]

		if key == arnKinesisKey {
			parts = strings.Split(arn, "/")
		}

		name = parts[len(parts)-1]
	} else {
		parts := strings.Split(arn, ".")

		keyParts := strings.Split(parts[0], "_")

		if parts[0] == "module" {
			// TODO: Add support to more type of modules
			key = arnLambdaKey
		} else if len(keyParts) > 1 {
			key = keyParts[1]
		} else {
			key = strings.Join(keyParts, "_")
		}

		name = parts[1]

		switch key {
		case arnKinesisKey:
			name = toPascalFromEnvar(name, name, suffixKinesis)
		case arnLambdaKey:
			name = toCamelFromEnvar(name, name, suffixLambda)
		case arnSQSKey:
			name = toKebabFromEnvar(name, name, suffixSQS)
		}
	}

	return resourceARN{key: key, name: name}
}

func strTransformFromEnvar(
	key, value, suffix string, f func(s string) string,
) string {
	var result string

	if key == suffix {
		suffixMap := map[string]struct{}{
			labelAWSKinesisStream:  {},
			labelAWSLambdaFunction: {},
			labelAWSS3Bucket:       {},
			labelAWSSQSQueue:       {},
		}

		result = value

		if strings.HasPrefix(result, "var.client-var.environment-") { // TODO:
			result = f(strings.ReplaceAll(result, "var.client-var.environment-", "")) // TODO:
		}

		for s := range suffixMap {
			if strings.HasPrefix(result, s) {
				result = f(resourceByARN(result).name)
				break
			}
		}
	} else {
		result = key

		result = strings.ReplaceAll(result, "_"+suffix, "")
		result = strings.ReplaceAll(result, suffix, "")
		result = f(result)
	}

	return result
}

func toCamelFromEnvar(key, value, suffix string) string {
	return strTransformFromEnvar(key, value, suffix, strcase.ToCamel)
}

func toKebabFromEnvar(key, value, suffix string) string {
	return strTransformFromEnvar(key, value, suffix, strcase.ToKebab)
}

func toPascalFromEnvar(key, value, suffix string) string {
	return strTransformFromEnvar(key, value, suffix, strcase.ToPascal)
}

func lambdaName(str, suffix string) string {
	return strcase.ToCamel(str[:len(str)-len(suffix)])
}

func s3BucketName(str, suffix string) string {
	return strcase.ToKebab(str[:len(str)-len(suffix)])
}

func sqsName(str, suffix string) string {
	return strcase.ToKebab(str[:len(str)-len(suffix)])
}
