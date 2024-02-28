package drawio

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseResources(t *testing.T) {
	type args struct {
		mxFile *MxFile
	}

	lambdaResource := NewGenericResource("1", "myReceiver", LambdaType)
	sqsResource := NewGenericResource("2", "my-sqs", SQSType)

	tests := []struct {
		name      string
		args      args
		want      *ResourceCollection
		targetErr error
	}{
		{
			name: "Lambda Resource",
			args: args{
				mxFile: &MxFile{
					Diagram: Diagram{
						MxGraphModel: MxGraphModel{
							Root: Root{
								MxCells: []MxCell{
									{ID: "1", Value: "myReceiver", Style: "mxgraph.aws3.lambda"},
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
			name: "SQS Resource",
			args: args{
				mxFile: &MxFile{
					Diagram: Diagram{
						MxGraphModel: MxGraphModel{
							Root: Root{
								MxCells: []MxCell{
									{ID: "2", Value: "my-sqs", Style: "mxgraph.aws3.sqs"},
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
									{ID: "3", Value: "my-sns", Style: "mxgraph.aws3.sns"},
								},
							},
						},
					},
				},
			},
			want: &ResourceCollection{
				Resources: []Resource{NewGenericResource("3", "my-sns", SNSType)},
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
									{ID: "4", Value: "myScheduler", Style: "mxgraph.aws4.event_time_based"},
								},
							},
						},
					},
				},
			},
			want: &ResourceCollection{
				Resources: []Resource{NewGenericResource("4", "myScheduler", CronType)},
			},
		},
		{
			name: "API Gateway Resource",
			args: args{
				mxFile: &MxFile{
					Diagram: Diagram{
						MxGraphModel: MxGraphModel{
							Root: Root{
								MxCells: []MxCell{
									{ID: "5", Value: "myAPI", Style: "mxgraph.aws3.api_gateway"},
								},
							},
						},
					},
				},
			},
			want: &ResourceCollection{
				Resources: []Resource{NewGenericResource("5", "myAPI", APIGatewayType)},
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
									{ID: "6", Value: "myEndpoint", Style: "mxgraph.aws4.endpoint"},
								},
							},
						},
					},
				},
			},
			want: &ResourceCollection{
				Resources: []Resource{NewGenericResource("6", "myEndpoint", EndpointType)},
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
									{ID: "7", Value: "myBucket", Style: "mxgraph.aws3.s3"},
								},
							},
						},
					},
				},
			},
			want: &ResourceCollection{
				Resources: []Resource{NewGenericResource("7", "myBucket", S3Type)},
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
									{ID: "8", Value: "myDB", Style: "mxgraph.flowchart.database"},
								},
							},
						},
					},
				},
			},
			want: &ResourceCollection{
				Resources: []Resource{NewGenericResource("8", "myDB", DatabaseType)},
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
									{ID: "9", Value: "myRestAPI", Style: "mxgraph.veeam2.restful_api"},
								},
							},
						},
					},
				},
			},
			want: &ResourceCollection{
				Resources: []Resource{NewGenericResource("9", "myRestAPI", RestfulAPIType)},
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
									{ID: "10", Value: "myKinesis", Style: "mxgraph.aws3.kinesis"},
								},
							},
						},
					},
				},
			},
			want: &ResourceCollection{
				Resources: []Resource{NewGenericResource("10", "myKinesis", KinesisType)},
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
									{ID: "1", Value: "myReceiver", Style: "mxgraph.aws3.lambda"},
									{ID: "2", Value: "my-sqs", Style: "mxgraph.aws3.sqs"},
									{ID: "3", Source: "1", Target: "2"},
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
