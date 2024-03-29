package drawiotoresources

import (
	"testing"

	drawioxml "github.com/joselitofilho/drawio-parser-go/pkg/parser/xml"

	"github.com/joselitofilho/aws-terraform-generator/internal/resources"

	"github.com/stretchr/testify/require"
)

func TestParseResources(t *testing.T) {
	type args struct {
		mxFile *drawioxml.MxFile
	}

	lambdaResource := resources.NewGenericResource("LAMBDA_ID", "myReceiver", resources.LambdaType)
	sqsResource := resources.NewGenericResource("SQS_ID", "my-sqs", resources.SQSType)

	tests := []struct {
		name      string
		args      args
		want      *resources.ResourceCollection
		targetErr error
	}{
		{
			name: "API Gateway Resource",
			args: args{
				mxFile: &drawioxml.MxFile{
					Diagram: drawioxml.Diagram{
						MxGraphModel: drawioxml.MxGraphModel{
							Root: drawioxml.Root{
								MxCells: []drawioxml.MxCell{
									{ID: "APIG_ID", Value: "myAPI", Style: "mxgraph.aws3.api_gateway"},
								},
							},
						},
					},
				},
			},
			want: &resources.ResourceCollection{
				Resources: []resources.Resource{
					resources.NewGenericResource("APIG_ID", "myAPI", resources.APIGatewayType)},
			},
		},
		{
			name: "Cron Resource",
			args: args{
				mxFile: &drawioxml.MxFile{
					Diagram: drawioxml.Diagram{
						MxGraphModel: drawioxml.MxGraphModel{
							Root: drawioxml.Root{
								MxCells: []drawioxml.MxCell{
									{ID: "CRON_ID", Value: "myScheduler", Style: "mxgraph.aws4.event_time_based"},
								},
							},
						},
					},
				},
			},
			want: &resources.ResourceCollection{
				Resources: []resources.Resource{
					resources.NewGenericResource("CRON_ID", "myScheduler", resources.CronType)},
			},
		},
		{
			name: "Database Resource",
			args: args{
				mxFile: &drawioxml.MxFile{
					Diagram: drawioxml.Diagram{
						MxGraphModel: drawioxml.MxGraphModel{
							Root: drawioxml.Root{
								MxCells: []drawioxml.MxCell{
									{ID: "DB_ID", Value: "myDB", Style: "mxgraph.flowchart.database"},
								},
							},
						},
					},
				},
			},
			want: &resources.ResourceCollection{
				Resources: []resources.Resource{resources.NewGenericResource("DB_ID", "myDB", resources.DatabaseType)},
			},
		},
		{
			name: "Endpoint Resource",
			args: args{
				mxFile: &drawioxml.MxFile{
					Diagram: drawioxml.Diagram{
						MxGraphModel: drawioxml.MxGraphModel{
							Root: drawioxml.Root{
								MxCells: []drawioxml.MxCell{
									{ID: "ENDPOINT_ID", Value: "myEndpoint", Style: "mxgraph.aws4.endpoint"},
								},
							},
						},
					},
				},
			},
			want: &resources.ResourceCollection{
				Resources: []resources.Resource{
					resources.NewGenericResource("ENDPOINT_ID", "myEndpoint", resources.EndpointType)},
			},
		},
		{
			name: "GoogleBQ Resource",
			args: args{
				mxFile: &drawioxml.MxFile{
					Diagram: drawioxml.Diagram{
						MxGraphModel: drawioxml.MxGraphModel{
							Root: drawioxml.Root{
								MxCells: []drawioxml.MxCell{
									{ID: "GBC_ID", Value: "myGBC", Style: "google_bigquery"},
								},
							},
						},
					},
				},
			},
			want: &resources.ResourceCollection{
				Resources: []resources.Resource{
					resources.NewGenericResource("GBC_ID", "myGBC", resources.GoogleBQType)},
			},
		},
		{
			name: "Kinesis Resource",
			args: args{
				mxFile: &drawioxml.MxFile{
					Diagram: drawioxml.Diagram{
						MxGraphModel: drawioxml.MxGraphModel{
							Root: drawioxml.Root{
								MxCells: []drawioxml.MxCell{
									{ID: "KINESIS_ID", Value: "myKinesis", Style: "mxgraph.aws3.kinesis"},
								},
							},
						},
					},
				},
			},
			want: &resources.ResourceCollection{
				Resources: []resources.Resource{
					resources.NewGenericResource("KINESIS_ID", "myKinesis", resources.KinesisType)},
			},
		},
		{
			name: "Lambda Resource",
			args: args{
				mxFile: &drawioxml.MxFile{
					Diagram: drawioxml.Diagram{
						MxGraphModel: drawioxml.MxGraphModel{
							Root: drawioxml.Root{
								MxCells: []drawioxml.MxCell{
									{ID: "LAMBDA_ID", Value: "myReceiver", Style: "mxgraph.aws3.lambda"},
								},
							},
						},
					},
				},
			},
			want: &resources.ResourceCollection{
				Resources: []resources.Resource{lambdaResource},
			},
		},
		{
			name: "Restful API Resource",
			args: args{
				mxFile: &drawioxml.MxFile{
					Diagram: drawioxml.Diagram{
						MxGraphModel: drawioxml.MxGraphModel{
							Root: drawioxml.Root{
								MxCells: []drawioxml.MxCell{
									{ID: "RESTFULAPI_ID", Value: "myRestAPI", Style: "mxgraph.veeam2.restful_api"},
								},
							},
						},
					},
				},
			},
			want: &resources.ResourceCollection{
				Resources: []resources.Resource{
					resources.NewGenericResource("RESTFULAPI_ID", "myRestAPI", resources.RestfulAPIType)},
			},
		},
		{
			name: "S3 Resource",
			args: args{
				mxFile: &drawioxml.MxFile{
					Diagram: drawioxml.Diagram{
						MxGraphModel: drawioxml.MxGraphModel{
							Root: drawioxml.Root{
								MxCells: []drawioxml.MxCell{
									{ID: "S3BUCKET_ID", Value: "myBucket", Style: "mxgraph.aws3.s3"},
								},
							},
						},
					},
				},
			},
			want: &resources.ResourceCollection{
				Resources: []resources.Resource{
					resources.NewGenericResource("S3BUCKET_ID", "myBucket", resources.S3Type)},
			},
		},
		{
			name: "SQS Resource",
			args: args{
				mxFile: &drawioxml.MxFile{
					Diagram: drawioxml.Diagram{
						MxGraphModel: drawioxml.MxGraphModel{
							Root: drawioxml.Root{
								MxCells: []drawioxml.MxCell{
									{ID: "SQS_ID", Value: "my-sqs", Style: "mxgraph.aws3.sqs"},
								},
							},
						},
					},
				},
			},
			want: &resources.ResourceCollection{
				Resources: []resources.Resource{sqsResource},
			},
		},
		{
			name: "SNS Resource",
			args: args{
				mxFile: &drawioxml.MxFile{
					Diagram: drawioxml.Diagram{
						MxGraphModel: drawioxml.MxGraphModel{
							Root: drawioxml.Root{
								MxCells: []drawioxml.MxCell{
									{ID: "SNS_ID", Value: "my-sns", Style: "mxgraph.aws3.sns"},
								},
							},
						},
					},
				},
			},
			want: &resources.ResourceCollection{
				Resources: []resources.Resource{resources.NewGenericResource("SNS_ID", "my-sns", resources.SNSType)},
			},
		},
		{
			name: "Empty MxFile",
			args: args{
				mxFile: &drawioxml.MxFile{
					Diagram: drawioxml.Diagram{
						MxGraphModel: drawioxml.MxGraphModel{
							Root: drawioxml.Root{
								MxCells: []drawioxml.MxCell{},
							},
						},
					},
				},
			},
			want: resources.NewResourceCollection(),
		},
		{
			name: "Two Connected Resources",
			args: args{
				mxFile: &drawioxml.MxFile{
					Diagram: drawioxml.Diagram{
						MxGraphModel: drawioxml.MxGraphModel{
							Root: drawioxml.Root{
								MxCells: []drawioxml.MxCell{
									{ID: "LAMBDA_ID", Value: "myReceiver", Style: "mxgraph.aws3.lambda"},
									{ID: "SQS_ID", Value: "my-sqs", Style: "mxgraph.aws3.sqs"},
									{ID: "3", Source: "LAMBDA_ID", Target: "SQS_ID"},
								},
							},
						},
					},
				},
			},
			want: &resources.ResourceCollection{
				Resources:     []resources.Resource{lambdaResource, sqsResource},
				Relationships: []resources.Relationship{{Source: lambdaResource, Target: sqsResource}},
			},
		},
		{
			name: "Single Unknown Resource",
			args: args{
				mxFile: &drawioxml.MxFile{
					Diagram: drawioxml.Diagram{
						MxGraphModel: drawioxml.MxGraphModel{
							Root: drawioxml.Root{
								MxCells: []drawioxml.MxCell{
									{ID: "1", Value: "Resource A", Style: "styleA"},
								},
							},
						},
					},
				},
			},
			want: resources.NewResourceCollection(),
		},
		{
			name: "Two Connected Unknown Resources",
			args: args{
				mxFile: &drawioxml.MxFile{
					Diagram: drawioxml.Diagram{
						MxGraphModel: drawioxml.MxGraphModel{
							Root: drawioxml.Root{
								MxCells: []drawioxml.MxCell{
									{ID: "1", Value: "Resource A", Style: "styleA"},
									{ID: "2", Value: "Resource B", Style: "styleB"},
									{ID: "3", Source: "1", Target: "2"},
								},
							},
						},
					},
				},
			},
			want: resources.NewResourceCollection(),
		},
		{
			name: "Multiple Unknown Resources",
			args: args{
				mxFile: &drawioxml.MxFile{
					Diagram: drawioxml.Diagram{
						MxGraphModel: drawioxml.MxGraphModel{
							Root: drawioxml.Root{
								MxCells: []drawioxml.MxCell{
									{ID: "1", Value: "Resource A", Style: "styleA"},
									{ID: "2", Value: "Resource B", Style: "styleB"},
								},
							},
						},
					},
				},
			},
			want: resources.NewResourceCollection(),
		},
	}

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			got, err := Transform(tc.args.mxFile)

			if tc.targetErr == nil {
				require.NoError(t, err)
				require.Equal(t, tc.want, got)
			} else {
				require.ErrorIs(t, err, tc.targetErr)
			}
		})
	}
}
