package graphviz

import (
	_ "embed"
	"testing"

	"github.com/joselitofilho/aws-terraform-generator/internal/resources"
	"github.com/stretchr/testify/require"
)

var (
	//go:embed testdata/happy_path.dot
	happyPath []byte

	//go:embed testdata/lr_orientation.dot
	lrOrientation []byte

	//go:embed testdata/source_or_target_nil.dot
	sourceOrTargetNil []byte
)

func TestBuild(t *testing.T) {
	type args struct {
		resc             *resources.ResourceCollection
		resourceImageMap map[resources.ResourceType]string
		config           Config
	}

	lambdaResource := resources.NewGenericResource("1", "MyLambda", resources.LambdaType)
	sqsResource := resources.NewGenericResource("2", "my-queue", resources.SQSType)

	reourceImageMap := map[resources.ResourceType]string{
		"lambda": "images/lambda.svg",
		"sqs":    "images/sqs.svg",
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "happy path",
			args: args{
				resc: &resources.ResourceCollection{
					Resources: []resources.Resource{lambdaResource, sqsResource},
					Relationships: []resources.Relationship{{
						Source: lambdaResource,
						Target: sqsResource,
					}},
				},
				resourceImageMap: reourceImageMap,
				config:           Config{},
			},
			want: string(happyPath),
		},
		{
			name: "left-right orientation",
			args: args{
				resc: &resources.ResourceCollection{
					Resources: []resources.Resource{lambdaResource, sqsResource},
					Relationships: []resources.Relationship{{
						Source: lambdaResource,
						Target: sqsResource,
					}},
				},
				resourceImageMap: reourceImageMap,
				config:           Config{Orientation: "LR"},
			},
			want: string(lrOrientation),
		},
		{
			name: "empty graph",
			args: args{
				resc:             resources.NewResourceCollection(),
				resourceImageMap: map[resources.ResourceType]string{},
				config:           Config{},
			},
			want: "digraph  {\n\t\n\t\n}\n",
		},
		{
			name: "skip edge for nil sources",
			args: args{
				resc: &resources.ResourceCollection{
					Resources: []resources.Resource{lambdaResource, sqsResource},
					Relationships: []resources.Relationship{{
						Source: nil,
						Target: sqsResource,
					}},
				},
				resourceImageMap: reourceImageMap,
				config:           Config{},
			},
			want: string(sourceOrTargetNil),
		},
		{
			name: "skip edge for nil targets",
			args: args{
				resc: &resources.ResourceCollection{
					Resources: []resources.Resource{lambdaResource, sqsResource},
					Relationships: []resources.Relationship{{
						Source: lambdaResource,
						Target: nil,
					}},
				},
				resourceImageMap: reourceImageMap,
				config:           Config{},
			},
			want: string(sourceOrTargetNil),
		},
	}

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			got := Build(tc.args.resc, tc.args.resourceImageMap, tc.args.config)

			require.Equal(t, tc.want, got)
		})
	}
}
