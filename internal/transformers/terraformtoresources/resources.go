package terraformtoresources

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/joselitofilho/aws-terraform-generator/internal/fmtcolor"
	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
	"github.com/joselitofilho/aws-terraform-generator/internal/resources"
	"github.com/joselitofilho/aws-terraform-generator/internal/transformers"
)

type ResourceARN struct {
	Type  string
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

func ResourceByARN(arn string, restType resources.ResourceType) ResourceARN {
	var arnType, name, label string

	if strings.HasPrefix(arn, "arn:") {
		parts := strings.Split(arn, ":")
		arnType = fmt.Sprintf("aws_%s_%s", parts[2], arnKeySuffix[parts[2]])

		if arnType == labelAWSKinesisStream {
			parts = strings.Split(arn, "/")
		}

		name = parts[len(parts)-1]
	} else if strings.HasPrefix(arn, "http") {
		parts := strings.Split(arn, "//")
		parts = strings.Split(parts[1], "/")

		resStrType := strings.Split(parts[0], ".")[0]
		arnType = fmt.Sprintf("aws_%s_%s", resStrType, arnKeySuffix[resStrType])

		name = parts[len(parts)-1]
	} else {
		parts := strings.Split(arn, ".")

		if len(parts) > 0 && parts[0] == "module" {
			// TODO: Add support to more type of modules
			arnType = labelAWSLambdaFunction
			name = parts[1]
			label = parts[1]
		} else if len(parts) > 1 && strings.HasPrefix(parts[0], "aws_") {
			arnType = parts[0]
			label = parts[1]
			name = parts[1]
		} else {
			name = arn
		}

		switch arnType {
		case labelAWSKinesisStream:
			name = strTransformFromKeyValue(name, name, suffixKinesis, restType, resources.ToKinesisCase)
		case labelAWSLambdaFunction:
			name = strTransformFromKeyValue(name, name, suffixLambda, restType, resources.ToLambdaCase)
		case labelAWSS3Bucket:
			name = strTransformFromKeyValue(name, name, suffixS3Bucket, restType, resources.ToS3BucketCase)
		case labelAWSSQSQueue:
			name = strTransformFromKeyValue(name, name, suffixSQS, restType, resources.ToSQSCase)
		case labelAWSSNSTopic:
			name = strTransformFromKeyValue(name, name, suffixSNS, restType, resources.ToSNSCase)
		}
	}

	if restType == resources.UnknownType {
		switch arnType {
		case labelAWSKinesisStream:
			restType = resources.KinesisType
		case labelAWSLambdaFunction:
			restType = resources.LambdaType
		case labelAWSS3Bucket:
			restType = resources.S3Type
		case labelAWSSQSQueue:
			restType = resources.SQSType
		case labelAWSSNSTopic:
			restType = resources.SNSType
		}
	}

	if arnType == "" {
		arnType = resourceARNByType[restType]
	}

	return ResourceARN{Type: arnType, Name: name, Label: label}
}

func strTransformFromKeyValue(
	key, value, suffix string, restType resources.ResourceType, f func(s string) string,
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
				result = ResourceByARN(result, restType).Name
				break
			}
		}

		return f(result)
	}

	return transformers.ReplaceSuffix(key, suffix, f)
}
