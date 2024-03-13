package terraformtoresources

import (
	"testing"

	"github.com/joselitofilho/aws-terraform-generator/internal/generators/terraform"

	"github.com/stretchr/testify/require"
)

func Test_replaceVars(t *testing.T) {
	type args struct {
		str             string
		tfVariables     []*terraform.Variable
		tfLocals        []*terraform.Local
		replaceableStrs map[string]string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "local string",
			args: args{
				str: "local.api_domain",
				tfLocals: []*terraform.Local{{Attributes: map[string]any{
					"environment": "dev",
					"api_domain":  "location-api.flyingtiger-local.environment.xiatechs.co.uk",
				}}},
			},
			want: "location-api.flyingtiger-dev.xiatechs.co.uk",
		},
		{
			name: "local string array",
			args: args{
				str: "local.api_domain",
				tfLocals: []*terraform.Local{{Attributes: map[string]any{
					"environment": []string{"dev", "prd"},
					"api_domain":  "location-api.flyingtiger-local.environment.xiatechs.co.uk",
				}}},
			},
			want: "location-api.flyingtiger-dev.xiatechs.co.uk",
		},
		{
			name: "local empty string array",
			args: args{
				str: "local.api_domain",
				tfLocals: []*terraform.Local{{Attributes: map[string]any{
					"environment": []string{},
					"api_domain":  "location-api.flyingtiger-local.environment.xiatechs.co.uk",
				}}},
			},
			want: "location-api.flyingtiger-local.environment.xiatechs.co.uk",
		},
		{
			name: "local string map",
			args: args{
				str: "local.api_domain",
				tfLocals: []*terraform.Local{{Attributes: map[string]any{
					"environment": map[string]any{"dev": struct{}{}, "prd": struct{}{}},
					"api_domain":  "location-api.flyingtiger-local.environment.xiatechs.co.uk",
				}}},
			},
			want: "location-api.flyingtiger-dev.xiatechs.co.uk",
		},
		{
			name: "local empty string map",
			args: args{
				str: "local.api_domain",
				tfLocals: []*terraform.Local{{Attributes: map[string]any{
					"environment": map[string]any{},
					"api_domain":  "location-api.flyingtiger-local.environment.xiatechs.co.uk",
				}}},
			},
			want: "location-api.flyingtiger-local.environment.xiatechs.co.uk",
		},
		{
			name: "local other types",
			args: args{
				str: "local.api_domain",
				tfLocals: []*terraform.Local{{Attributes: map[string]any{
					"environment": 1,
					"api_domain":  "location-api.flyingtiger-local.environment.xiatechs.co.uk",
				}}},
			},
			want: "location-api.flyingtiger-local.environment.xiatechs.co.uk",
		},
		{
			name: "var string",
			args: args{
				str: "var.api_domain",
				tfVariables: []*terraform.Variable{{Attributes: map[string]any{
					"environment": "dev",
					"api_domain":  "location-api.flyingtiger-var.environment.xiatechs.co.uk",
				}}},
			},
			want: "location-api.flyingtiger-dev.xiatechs.co.uk",
		},
		{
			name: "replaceable texts",
			args: args{
				str: "var.ze-location",
				replaceableStrs: map[string]string{
					"var.ze-": "",
				},
			},
			want: "location",
		},
	}

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			got := replaceVars(tc.args.str, tc.args.tfVariables, tc.args.tfLocals, tc.args.replaceableStrs)

			require.Equal(t, tc.want, got)
		})
	}
}
