package drawio

import (
	_ "embed"
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseXML(t *testing.T) {
	type args struct {
		fileName string
	}

	tests := []struct {
		name      string
		args      args
		want      *MxFile
		targetErr error
	}{
		{
			name: "happy path",
			args: args{
				fileName: "testdata/diagram.xml",
			},
			want: &MxFile{
				XMLName: xml.Name{Local: "mxfile"},
				Diagram: Diagram{
					XMLName: xml.Name{Local: "diagram"},
					MxGraphModel: MxGraphModel{
						XMLName: xml.Name{Local: "mxGraphModel"},
						Root: Root{
							XMLName: xml.Name{Local: "root"},
							MxCells: []MxCell{{
								XMLName: xml.Name{Local: "mxCell"},
								ID:      "kVijt7gfVD9ZtySMmpSK-1",
								Value:   "myReceiver",
								Style: "outlineConnect=0;dashed=0;verticalLabelPosition=bottom;verticalAlign=top;" +
									"align=center;html=1;shape=mxgraph.aws3.lambda;fillColor=#F58534;" +
									"gradientColor=none;",
								Parent: "1",
								Vertex: true,
								Geometry: Geometry{
									XMLName: xml.Name{Local: "mxGeometry"},
									X:       850,
									Y:       -1121.5,
									Width:   76.5,
									Height:  93,
								},
							}},
						},
					},
				},
			},
		},
	}

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			got, err := ParseXML(tc.args.fileName)

			if tc.targetErr == nil {
				require.NoError(t, err)
				require.Equal(t, tc.want, got)
			} else {
				require.ErrorIs(t, err, tc.targetErr)
			}
		})
	}
}
