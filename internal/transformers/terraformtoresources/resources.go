package terraformtoresources

import (
	"regexp"
	"strings"

	"github.com/joselitofilho/aws-terraform-generator/internal/fmtcolor"
	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
	"github.com/joselitofilho/aws-terraform-generator/internal/resources"
	"github.com/joselitofilho/aws-terraform-generator/internal/transformers"
)

type ResourceARN struct {
	Key   string
	Name  string
	Label string
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

func ResourceByARN(arn string) ResourceARN {
	var key, name, label string

	if strings.HasPrefix(arn, "arn:") {
		parts := strings.Split(arn, ":")

		key = parts[2]

		if key == arnKinesisKey {
			parts = strings.Split(arn, "/")
		}

		name = parts[len(parts)-1]
	} else if strings.HasPrefix(arn, "http") {
		parts := strings.Split(arn, "//")
		parts = strings.Split(parts[1], "/")

		key = strings.Split(parts[0], ".")[0]

		name = parts[len(parts)-1]
	} else {
		parts := strings.Split(arn, ".")

		if len(parts) > 0 && parts[0] == "module" {
			// TODO: Add support to more type of modules
			key = arnLambdaKey
			name = parts[1]
			label = parts[1]
		} else if len(parts) > 1 && strings.HasPrefix(parts[0], "aws_") {
			label = parts[1]
			keyParts := strings.Split(parts[0], "_")

			if len(keyParts) > 1 {
				key = keyParts[1]
			} else {
				key = strings.Join(keyParts, "_")
			}

			name = parts[1]
		} else {
			name = arn
		}

		switch key {
		case arnKinesisKey:
			name = strTransformFromKeyValue(name, name, suffixKinesis, resources.ToKinesisCase)
		case arnLambdaKey:
			name = strTransformFromKeyValue(name, name, suffixLambda, resources.ToLambdaCase)
		case arnSQSKey:
			name = strTransformFromKeyValue(name, name, suffixSQS, resources.ToSQSCase)
		}
	}

	return ResourceARN{Key: key, Name: name, Label: label}
}

func strTransformFromKeyValue(
	key, value, suffix string, f func(s string) string,
) string {
	if key == suffix {
		suffixMap := map[string]struct{}{
			labelAWSKinesisStream:  {},
			labelAWSLambdaFunction: {},
			labelAWSS3Bucket:       {},
			labelAWSSQSQueue:       {},
		}

		result := value

		for s := range suffixMap {
			if strings.HasPrefix(result, s) {
				result = ResourceByARN(result).Name
				break
			}
		}

		return f(result)
	}

	return transformers.ReplaceSuffix(key, suffix, f)
}
