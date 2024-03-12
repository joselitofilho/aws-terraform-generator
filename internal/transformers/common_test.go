package transformers

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReplaceSuffix(t *testing.T) {
	type args struct {
		value  string
		suffix string
		fn     func(s string) string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "suffix stars with _",
			args: args{
				value:  "NAME_SUFFIX",
				suffix: "SUFFIX",
				fn: func(s string) string {
					require.Equal(t, s, "NAME")
					return strings.ToLower(s)
				},
			},
			want: "name",
		},
		{
			name: "suffix does not start with _",
			args: args{
				value:  "NAMESUFFIX",
				suffix: "SUFFIX",
				fn: func(s string) string {
					require.Equal(t, s, "NAME")
					return strings.ToLower(s)
				},
			},
			want: "name",
		},
		{
			name: "when fn is nil",
			args: args{
				value:  "NAME_SUFFIX",
				suffix: "SUFFIX",
				fn:     nil,
			},
			want: "NAME",
		},
	}

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			got := ReplaceSuffix(tc.args.value, tc.args.suffix, tc.args.fn)

			require.Equal(t, tc.want, got)
		})
	}
}
