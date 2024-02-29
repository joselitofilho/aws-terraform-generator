package guides

import (
	"testing"

	"github.com/AlecAivazis/survey/v2"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGuideCode(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		surveyAsker *FakeSurveyAsker
		workdir     string
		fileMap     map[string][]string
	}

	tests := []struct {
		name         string
		args         args
		prepareMocks func(*FakeSurveyAsker)
		want         *CodeAnswers
		targetErr    error
	}{
		{
			name: "",
			args: args{
				surveyAsker: &FakeSurveyAsker{TB: t, Answers: CodeAnswers{
					StackName: "teststack",
					Config:    "diagram.yaml",
					Output:    "./testoutput",
				}},
				workdir: "./testoutput/teststask",
				fileMap: map[string][]string{
					"config": {"diagram.yaml"},
				},
			},
			prepareMocks: func(_ *FakeSurveyAsker) {},
			want: &CodeAnswers{
				StackName: "teststack",
				Config:    "testoutput/teststask/diagram.yaml",
				Output:    "./testoutput",
			},
		},
	}

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			tc.prepareMocks(tc.args.surveyAsker)

			got, err := GuideCode(tc.args.surveyAsker, tc.args.workdir, tc.args.fileMap)

			require.ErrorIs(t, err, tc.targetErr)
			require.Equal(t, tc.want, got)
		})
	}
}

type FakeSurveyAsker struct {
	TB        testing.TB
	Answers   CodeAnswers
	callCount int
}

func (f *FakeSurveyAsker) Ask(questions []*survey.Question, response any, _ ...survey.AskOpt) error {
	f.callCount++

	ans := response.(*CodeAnswers)

	switch f.callCount {
	case 1:
		require.Equal(f.TB, "stackName", questions[0].Name)
		ans.StackName = f.Answers.StackName
	case 2:
		require.Equal(f.TB, "config", questions[0].Name)
		ans.Config = f.Answers.Config
	case 3:
		require.Equal(f.TB, "output", questions[0].Name)
		ans.Output = f.Answers.Output
	}

	return nil
}

// AskOne asks a single survey question using the actual survey library.
func (f *FakeSurveyAsker) AskOne(_ survey.Prompt, _ any, _ ...survey.AskOpt) error {
	return nil
}
