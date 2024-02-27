package guides

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_replaceDoubleSlash(t *testing.T) {
	type args struct {
		str string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "happy path",
			args: args{str: "a//b"},
			want: "a/b",
		},
	}

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			got := replaceDoubleSlash(tc.args.str)

			require.Equal(t, tc.want, got)
		})
	}
}
