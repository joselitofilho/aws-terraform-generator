package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMergeStringMap(t *testing.T) {
	type args struct {
		left  map[string]string
		right map[string]string
	}

	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "new values",
			args: args{
				left:  map[string]string{"a": "1", "b": "2"},
				right: map[string]string{"c": "3", "d": "4"},
			},
			want: map[string]string{"a": "1", "b": "2", "c": "3", "d": "4"},
		},
		{
			name: "same values",
			args: args{
				left:  map[string]string{"a": "1", "b": "2"},
				right: map[string]string{"a": "1", "b": "2"},
			},
			want: map[string]string{"a": "1", "b": "2"},
		},
		{
			name: "at least one different value",
			args: args{
				left:  map[string]string{"a": "1", "b": "2"},
				right: map[string]string{"b": "2", "c": "3"},
			},
			want: map[string]string{"a": "1", "b": "2", "c": "3"},
		},
		{
			name: "right overrides left",
			args: args{
				left:  map[string]string{"a": "1", "b": "2"},
				right: map[string]string{"b": "22", "c": "3"},
			},
			want: map[string]string{"a": "1", "b": "22", "c": "3"},
		},
	}

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			got := MergeStringMap(tc.args.left, tc.args.right)

			require.Equal(t, tc.want, got)
		})
	}
}
