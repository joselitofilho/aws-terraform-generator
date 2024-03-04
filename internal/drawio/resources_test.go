package drawio

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseResources(t *testing.T) {
	type args struct {
		mxFile *MxFile
	}

	lambdaResource := NewGenericResource("LAMBDA_ID", "myReceiver", LambdaType)
	sqsResource := NewGenericResource("SQS_ID", "my-sqs", SQSType)

	tests := []struct {
		name      string
		args      args
		want      *ResourceCollection
		targetErr error
	}{
		{
			name: "API Gateway Resource",
			args: args{
				mxFile: &MxFile{
					Diagram: Diagram{
						MxGraphModel: MxGraphModel{
							Root: Root{
								MxCells: []MxCell{
									{ID: "APIG_ID", Value: "myAPI", Style: "mxgraph.aws3.api_gateway"},
								},
							},
						},
					},
				},
			},
			want: &ResourceCollection{
				Resources: []Resource{NewGenericResource("APIG_ID", "myAPI", APIGatewayType)},
			},
		},
		{
			name: "Cron Resource",
			args: args{
				mxFile: &MxFile{
					Diagram: Diagram{
						MxGraphModel: MxGraphModel{
							Root: Root{
								MxCells: []MxCell{
									{ID: "CRON_ID", Value: "myScheduler", Style: "mxgraph.aws4.event_time_based"},
								},
							},
						},
					},
				},
			},
			want: &ResourceCollection{
				Resources: []Resource{NewGenericResource("CRON_ID", "myScheduler", CronType)},
			},
		},
		{
			name: "Database Resource",
			args: args{
				mxFile: &MxFile{
					Diagram: Diagram{
						MxGraphModel: MxGraphModel{
							Root: Root{
								MxCells: []MxCell{
									{ID: "DB_ID", Value: "myDB", Style: "mxgraph.flowchart.database"},
								},
							},
						},
					},
				},
			},
			want: &ResourceCollection{
				Resources: []Resource{NewGenericResource("DB_ID", "myDB", DatabaseType)},
			},
		},
		{
			name: "Endpoint Resource",
			args: args{
				mxFile: &MxFile{
					Diagram: Diagram{
						MxGraphModel: MxGraphModel{
							Root: Root{
								MxCells: []MxCell{
									{ID: "ENDPOINT_ID", Value: "myEndpoint", Style: "mxgraph.aws4.endpoint"},
								},
							},
						},
					},
				},
			},
			want: &ResourceCollection{
				Resources: []Resource{NewGenericResource("ENDPOINT_ID", "myEndpoint", EndpointType)},
			},
		},
		{
			name: "GoogleBQ Resource",
			args: args{
				mxFile: &MxFile{
					Diagram: Diagram{
						MxGraphModel: MxGraphModel{
							Root: Root{
								MxCells: []MxCell{
									{ID: "GBC_ID", Value: "myGBC", Style: "google_bigquery"},
								},
							},
						},
					},
				},
			},
			want: &ResourceCollection{
				Resources: []Resource{NewGenericResource("GBC_ID", "myGBC", GoogleBQType)},
			},
		},
		{
			name: "Kinesis Resource",
			args: args{
				mxFile: &MxFile{
					Diagram: Diagram{
						MxGraphModel: MxGraphModel{
							Root: Root{
								MxCells: []MxCell{
									{ID: "KINESIS_ID", Value: "myKinesis", Style: "mxgraph.aws3.kinesis"},
								},
							},
						},
					},
				},
			},
			want: &ResourceCollection{
				Resources: []Resource{NewGenericResource("KINESIS_ID", "myKinesis", KinesisType)},
			},
		},
		{
			name: "Lambda Resource",
			args: args{
				mxFile: &MxFile{
					Diagram: Diagram{
						MxGraphModel: MxGraphModel{
							Root: Root{
								MxCells: []MxCell{
									{ID: "LAMBDA_ID", Value: "myReceiver", Style: "mxgraph.aws3.lambda"},
								},
							},
						},
					},
				},
			},
			want: &ResourceCollection{
				Resources: []Resource{lambdaResource},
			},
		},
		{
			name: "Restful API Resource",
			args: args{
				mxFile: &MxFile{
					Diagram: Diagram{
						MxGraphModel: MxGraphModel{
							Root: Root{
								MxCells: []MxCell{
									{ID: "RESTFULAPI_ID", Value: "myRestAPI", Style: "mxgraph.veeam2.restful_api"},
								},
							},
						},
					},
				},
			},
			want: &ResourceCollection{
				Resources: []Resource{NewGenericResource("RESTFULAPI_ID", "myRestAPI", RestfulAPIType)},
			},
		},
		{
			name: "S3 Resource",
			args: args{
				mxFile: &MxFile{
					Diagram: Diagram{
						MxGraphModel: MxGraphModel{
							Root: Root{
								MxCells: []MxCell{
									{ID: "S3BUCKET_ID", Value: "myBucket", Style: "mxgraph.aws3.s3"},
								},
							},
						},
					},
				},
			},
			want: &ResourceCollection{
				Resources: []Resource{NewGenericResource("S3BUCKET_ID", "myBucket", S3Type)},
			},
		},
		{
			name: "SQS Resource",
			args: args{
				mxFile: &MxFile{
					Diagram: Diagram{
						MxGraphModel: MxGraphModel{
							Root: Root{
								MxCells: []MxCell{
									{ID: "SQS_ID", Value: "my-sqs", Style: "mxgraph.aws3.sqs"},
								},
							},
						},
					},
				},
			},
			want: &ResourceCollection{
				Resources: []Resource{sqsResource},
			},
		},
		{
			name: "SNS Resource",
			args: args{
				mxFile: &MxFile{
					Diagram: Diagram{
						MxGraphModel: MxGraphModel{
							Root: Root{
								MxCells: []MxCell{
									{ID: "SNS_ID", Value: "my-sns", Style: "mxgraph.aws3.sns"},
								},
							},
						},
					},
				},
			},
			want: &ResourceCollection{
				Resources: []Resource{NewGenericResource("SNS_ID", "my-sns", SNSType)},
			},
		},
		{
			name: "Empty MxFile",
			args: args{
				mxFile: &MxFile{
					Diagram: Diagram{
						MxGraphModel: MxGraphModel{
							Root: Root{
								MxCells: []MxCell{},
							},
						},
					},
				},
			},
			want: NewResourceCollection(),
		},
		{
			name: "Two Connected Resources",
			args: args{
				mxFile: &MxFile{
					Diagram: Diagram{
						MxGraphModel: MxGraphModel{
							Root: Root{
								MxCells: []MxCell{
									{ID: "LAMBDA_ID", Value: "myReceiver", Style: "mxgraph.aws3.lambda"},
									{ID: "SQS_ID", Value: "my-sqs", Style: "mxgraph.aws3.sqs"},
									{ID: "3", Source: "LAMBDA_ID", Target: "SQS_ID"},
								},
							},
						},
					},
				},
			},
			want: &ResourceCollection{
				Resources:     []Resource{lambdaResource, sqsResource},
				Relationships: []Relationship{{Source: lambdaResource, Target: sqsResource}},
			},
		},
		{
			name: "Single Unknown Resource",
			args: args{
				mxFile: &MxFile{
					Diagram: Diagram{
						MxGraphModel: MxGraphModel{
							Root: Root{
								MxCells: []MxCell{
									{ID: "1", Value: "Resource A", Style: "styleA"},
								},
							},
						},
					},
				},
			},
			want: NewResourceCollection(),
		},
		{
			name: "Two Connected Unknown Resources",
			args: args{
				mxFile: &MxFile{
					Diagram: Diagram{
						MxGraphModel: MxGraphModel{
							Root: Root{
								MxCells: []MxCell{
									{ID: "1", Value: "Resource A", Style: "styleA"},
									{ID: "2", Value: "Resource B", Style: "styleB"},
									{ID: "3", Source: "1", Target: "2"},
								},
							},
						},
					},
				},
			},
			want: NewResourceCollection(),
		},
		{
			name: "Multiple Unknown Resources",
			args: args{
				mxFile: &MxFile{
					Diagram: Diagram{
						MxGraphModel: MxGraphModel{
							Root: Root{
								MxCells: []MxCell{
									{ID: "1", Value: "Resource A", Style: "styleA"},
									{ID: "2", Value: "Resource B", Style: "styleB"},
								},
							},
						},
					},
				},
			},
			want: NewResourceCollection(),
		},
	}

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			got, err := ParseResources(tc.args.mxFile)

			if tc.targetErr == nil {
				require.NoError(t, err)
				require.Equal(t, tc.want, got)
			} else {
				require.ErrorIs(t, err, tc.targetErr)
			}
		})
	}
}
