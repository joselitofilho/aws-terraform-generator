package resources

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFindDifferences(t *testing.T) {
	type args struct {
		rc1 *ResourceCollection
		rc2 *ResourceCollection
	}

	lambda1Resource := NewGenericResource("1", "myReceiver", LambdaType)
	sqs1Resource := NewGenericResource("2", "my-queue", SQSType)

	lambda2Resource := NewGenericResource("1", "myProcessor", LambdaType)
	sqs2Resource := NewGenericResource("2", "my-q", SQSType)

	tests := []struct {
		name                       string
		args                       args
		wantAddedResourcesByType   map[ResourceType][]Resource
		wantRemovedResourcesByType map[ResourceType][]Resource
		wantAddedRelationships     []Relationship
		wantRemovedRelationships   []Relationship
	}{
		{
			name: "happy path",
			args: args{
				rc1: &ResourceCollection{
					Resources: []Resource{lambda1Resource, sqs1Resource},
					Relationships: []Relationship{
						{Source: lambda1Resource, Target: sqs1Resource},
					},
				},
				rc2: &ResourceCollection{
					Resources: []Resource{lambda2Resource, sqs2Resource},
					Relationships: []Relationship{
						{Source: lambda2Resource, Target: sqs2Resource},
					},
				},
			},
			wantAddedResourcesByType: map[ResourceType][]Resource{
				LambdaType: {lambda2Resource},
				SQSType:    {sqs2Resource},
			},
			wantRemovedResourcesByType: map[ResourceType][]Resource{
				LambdaType: {lambda1Resource},
				SQSType:    {sqs1Resource},
			},
			wantAddedRelationships: []Relationship{
				{Source: lambda2Resource, Target: sqs2Resource},
			},
			wantRemovedRelationships: []Relationship{
				{Source: lambda1Resource, Target: sqs1Resource},
			},
		},
		{
			name: "empty",
			args: args{
				rc1: &ResourceCollection{},
				rc2: &ResourceCollection{},
			},
			wantAddedResourcesByType:   map[ResourceType][]Resource{},
			wantRemovedResourcesByType: map[ResourceType][]Resource{},
			wantAddedRelationships:     nil,
			wantRemovedRelationships:   nil,
		},
	}

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			gotAddedResourcesByType, gotRemovedResourcesByType, gotAddedRelationships, gotRemovedRelationships :=
				FindDifferences(tc.args.rc1, tc.args.rc2)

			require.Equal(t, tc.wantAddedResourcesByType, gotAddedResourcesByType)
			require.Equal(t, tc.wantRemovedResourcesByType, gotRemovedResourcesByType)
			require.Equal(t, tc.wantAddedRelationships, gotAddedRelationships)
			require.Equal(t, tc.wantRemovedRelationships, gotRemovedRelationships)
		})
	}
}
