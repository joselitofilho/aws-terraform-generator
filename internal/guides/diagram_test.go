package guides

import (
	"testing"

	"github.com/AlecAivazis/survey/v2"
	surveyasker "github.com/joselitofilho/aws-terraform-generator/internal/survey"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGuideDiagram(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		surveyAsker surveyasker.Asker
		workdir     string
		fileMap     map[string][]string
	}

	tests := []struct {
		name         string
		args         args
		prepareMocks func(surveyasker.Asker)
		want         *DiagramAnswers
		targetErr    error
	}{
		{
			name: "happy path",
			args: args{
				surveyAsker: &fakeDiagramSurveyAsker{
					TB: t,
					Answers: DiagramAnswers{
						Diagram: "diagram.xml",
						Config:  "diagram.config.yaml",
						Output:  "./testoutput//diagram.yaml",
					},
					AskQuestions: []survey.Prompt{
						&survey.Select{
							Message: "Choose a diagram:",
							Options: []string{"digram.xml"},
						},
						&survey.Select{
							Message: "Choose a config:",
							Default: 1,
							Options: []string{"config.yaml", "diagram.config.yaml"},
						},
					},
					AskForOutput: &survey.Input{
						Message: "Enter the output file:",
						Default: "./testoutput/teststack/diagram.yaml",
					},
				},
				workdir: "./testoutput/teststack",
				fileMap: map[string][]string{
					"diagram": {"digram.xml"},
					"config":  {"config.yaml", "diagram.config.yaml"},
				},
			},
			want: &DiagramAnswers{
				Diagram: "./testoutput/teststack/diagram.xml",
				Config:  "./testoutput/teststack/diagram.config.yaml",
				Output:  "./testoutput/diagram.yaml",
			},
			prepareMocks: func(_ surveyasker.Asker) {},
		},
		{
			name: "when there is no config files should return an error",
			args: args{
				workdir: "./testoutput/teststack",
				fileMap: map[string][]string{"diagram": {}},
			},
			prepareMocks: func(_ surveyasker.Asker) {},
			targetErr:    ErrDirDoesNotContainAnyDiagramFile,
		},
		{
			name: "when there is no config files should return an error",
			args: args{
				workdir: "./testoutput/teststack",
				fileMap: map[string][]string{"diagram": {"digram.xml"}, "config": {}},
			},
			prepareMocks: func(_ surveyasker.Asker) {},
			targetErr:    ErrDirDoesNotContainAnyConfigFile,
		},
		{
			name: "when first ask fails should return an error",
			args: args{
				surveyAsker: surveyasker.NewMockAsker(ctrl),
				workdir:     "./testoutput/teststack",
				fileMap:     map[string][]string{"diagram": {"digram.xml"}, "config": {"diagram.config.yaml"}},
			},
			prepareMocks: func(a surveyasker.Asker) {
				msa := a.(*surveyasker.MockAsker)

				msa.EXPECT().Ask(gomock.Any(), gomock.Any()).Return(errDummy)
			},
			targetErr: errDummy,
		},
		{
			name: "when survey to enter the output fails should return an error",
			args: args{
				surveyAsker: surveyasker.NewMockAsker(ctrl),
				workdir:     "./testoutput/teststack",
				fileMap:     map[string][]string{"diagram": {"digram.xml"}, "config": {"diagram.config.yaml"}},
			},
			prepareMocks: func(a surveyasker.Asker) {
				msa := a.(*surveyasker.MockAsker)

				gomock.InOrder(
					msa.EXPECT().Ask(gomock.Any(), gomock.Any()).Return(nil),
					msa.EXPECT().AskOne(gomock.Any(), gomock.Any()).Return(errDummy),
				)
			},
			targetErr: errDummy,
		},
	}

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			tc.prepareMocks(tc.args.surveyAsker)

			got, err := GuideDiagram(tc.args.surveyAsker, tc.args.workdir, tc.args.fileMap)

			if tc.targetErr == nil {
				require.NoError(t, err)
				require.Equal(t, tc.want, got)
			} else {
				require.ErrorIs(t, err, tc.targetErr)
			}
		})
	}
}

type fakeDiagramSurveyAsker struct {
	TB           testing.TB
	Answers      DiagramAnswers
	AskQuestions []survey.Prompt
	AskForOutput survey.Prompt
}

func (f *fakeDiagramSurveyAsker) Ask(questions []*survey.Question, response any, _ ...survey.AskOpt) error {
	require.IsType(f.TB, &survey.Select{}, questions[0].Prompt)
	require.Equal(f.TB, f.AskQuestions[0], questions[0].Prompt)

	require.IsType(f.TB, &survey.Select{}, questions[1].Prompt)
	require.Equal(f.TB, f.AskQuestions[1], questions[1].Prompt)

	ans := response.(*DiagramAnswers)

	ans.Diagram = f.Answers.Diagram
	ans.Config = f.Answers.Config

	return nil
}

func (f *fakeDiagramSurveyAsker) AskOne(prompt survey.Prompt, response any, _ ...survey.AskOpt) error {
	ans := response.(*string)
	require.Empty(f.TB, ans)
	require.IsType(f.TB, &survey.Input{}, f.AskForOutput)
	require.Equal(f.TB, f.AskForOutput, prompt)

	*ans = f.Answers.Output

	return nil
}
