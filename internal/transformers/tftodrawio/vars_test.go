package tftodrawio

import (
	"testing"

	"github.com/joselitofilho/aws-terraform-generator/internal/generators/terraform"

	"github.com/stretchr/testify/require"
)

func Test_replaceVars(t *testing.T) {
	type args struct {
		str      string
		tfLocals []*terraform.Local
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "",
			args: args{
				str: "local.api_domain",
				tfLocals: []*terraform.Local{{Attributes: map[string]any{
					"environment": "dev",
					"api_domain":  "location-api.flyingtiger-local.environment.xiatechs.co.uk",
				}}},
			},
			want: "location-api.flyingtiger-dev.xiatechs.co.uk",
		},
	}

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			got := replaceVars(tc.args.str, tc.args.tfLocals)

			require.Equal(t, tc.want, got)
		})
	}
}
