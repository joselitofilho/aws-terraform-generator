package terraformtoresources

import (
	"testing"

	hcl "github.com/joselitofilho/hcl-parser-go/pkg/parser/hcl"

	"github.com/stretchr/testify/require"
)

func Test_replaceVars(t *testing.T) {
	type args struct {
		str             string
		tfVariables     []*hcl.Variable
		tfLocals        []*hcl.Local
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
				tfLocals: []*hcl.Local{{Attributes: map[string]any{
					"environment": "dev",
					"api_domain":  "stack-api.domain-local.environment.com",
				}}},
			},
			want: "stack-api.domain-dev.com",
		},
		{
			name: "local string array",
			args: args{
				str: "local.api_domain",
				tfLocals: []*hcl.Local{{Attributes: map[string]any{
					"environment": []string{"dev", "prd"},
					"api_domain":  "stack-api.domain-local.environment.com",
				}}},
			},
			want: "stack-api.domain-dev.com",
		},
		{
			name: "local empty string array",
			args: args{
				str: "local.api_domain",
				tfLocals: []*hcl.Local{{Attributes: map[string]any{
					"environment": []string{},
					"api_domain":  "stack-api.domain-local.environment.com",
				}}},
			},
			want: "stack-api.domain-local.environment.com",
		},
		{
			name: "local string map",
			args: args{
				str: "local.api_domain",
				tfLocals: []*hcl.Local{{Attributes: map[string]any{
					"environment": map[string]any{"dev": struct{}{}, "prd": struct{}{}},
					"api_domain":  "stack-api.domain-local.environment.com",
				}}},
			},
			want: "stack-api.domain-dev.com",
		},
		{
			name: "local empty string map",
			args: args{
				str: "local.api_domain",
				tfLocals: []*hcl.Local{{Attributes: map[string]any{
					"environment": map[string]any{},
					"api_domain":  "stack-api.domain-local.environment.com",
				}}},
			},
			want: "stack-api.domain-local.environment.com",
		},
		{
			name: "local other types",
			args: args{
				str: "local.api_domain",
				tfLocals: []*hcl.Local{{Attributes: map[string]any{
					"environment": 1,
					"api_domain":  "stack-api.domain-local.environment.com",
				}}},
			},
			want: "stack-api.domain-local.environment.com",
		},
		{
			name: "var string",
			args: args{
				str: "var.api_domain",
				tfVariables: []*hcl.Variable{{Attributes: map[string]any{
					"environment": "dev",
					"api_domain":  "stack-api.domain-var.environment.com",
				}}},
			},
			want: "stack-api.domain-dev.com",
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
