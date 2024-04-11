package config

import (
	"testing"

	awsresources "github.com/joselitofilho/aws-terraform-generator/internal/resources"
	"github.com/stretchr/testify/require"
)

func TestImages_ToStringMap(t *testing.T) {
	tests := []struct {
		name string
		m    Images
		want map[string]string
	}{
		{
			name: "happy path",
			m: Images{
				awsresources.LambdaType: "./paht/to/image.svg",
			},
			want: map[string]string{
				"Lambda": "./paht/to/image.svg",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.m.ToStringMap()

			require.Equal(t, tt.want, got)
		})
	}
}
