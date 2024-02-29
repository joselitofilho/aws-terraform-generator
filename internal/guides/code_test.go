package guides

import (
	"errors"
	"testing"

	"github.com/AlecAivazis/survey/v2"

	surveyasker "github.com/joselitofilho/aws-terraform-generator/internal/survey"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

var errDummy = errors.New("dummy error")

func TestGuideCode(t *testing.T) {
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
		want         *CodeAnswers
		targetErr    error
	}{
		{
			name: "happy path",
			args: args{
				surveyAsker: &fakeCodeSurveyAsker{
					TB: t,
					Answers: CodeAnswers{
						StackName: "teststack",
						Config:    "diagram.yaml",
						Output:    ".//testoutput",
					},
					AskPerCall: []survey.Prompt{
						&survey.Input{
							Message: "Enter the stack name:",
							Default: "teststack",
						},
						&survey.Select{
							Message: "Choose a config:",
							Default: 1,
							Options: []string{"config.yaml", "diagram.yaml"},
						},
						&survey.Input{
							Message: "Enter the output folder:",
							Default: "./output",
						},
					},
				},
				workdir: "./testoutput/teststack",
				fileMap: map[string][]string{"config": {"config.yaml", "diagram.yaml"}},
			},
			prepareMocks: func(_ surveyasker.Asker) {},
			want: &CodeAnswers{
				StackName: "teststack",
				Config:    "testoutput/teststack/diagram.yaml",
				Output:    "./testoutput",
			},
		},
		{
			name: "when there is no config files should return an error",
			args: args{
				workdir: "./testoutput/teststack",
				fileMap: map[string][]string{"config": {}},
			},
			prepareMocks: func(_ surveyasker.Asker) {},
			targetErr:    ErrDirDoesNotContainAnyConfigFile,
		},
		{
			name: "when survey to enter the stack name fails should return an error",
			args: args{
				surveyAsker: surveyasker.NewMockAsker(ctrl),
				workdir:     "./testoutput/teststack",
				fileMap:     map[string][]string{"config": {"diagram.yaml"}},
			},
			prepareMocks: func(a surveyasker.Asker) {
				msa := a.(*surveyasker.MockAsker)

				msa.EXPECT().AskOne(gomock.Any(), gomock.Any()).Return(errDummy)
			},
			targetErr: errDummy,
		},
		{
			name: "when survey to choose the config fails should return an error",
			args: args{
				surveyAsker: surveyasker.NewMockAsker(ctrl),
				workdir:     "./testoutput/teststack",
				fileMap:     map[string][]string{"config": {"diagram.yaml"}},
			},
			prepareMocks: func(a surveyasker.Asker) {
				msa := a.(*surveyasker.MockAsker)

				gomock.InOrder(
					msa.EXPECT().AskOne(gomock.Any(), gomock.Any()).Return(nil),
					msa.EXPECT().AskOne(gomock.Any(), gomock.Any()).Return(errDummy),
				)
			},
			targetErr: errDummy,
		},
		{
			name: "when survey to enter the output fails should return an error",
			args: args{
				surveyAsker: surveyasker.NewMockAsker(ctrl),
				workdir:     "./testoutput/teststack",
				fileMap:     map[string][]string{"config": {"diagram.yaml"}},
			},
			prepareMocks: func(a surveyasker.Asker) {
				msa := a.(*surveyasker.MockAsker)

				gomock.InOrder(
					msa.EXPECT().AskOne(gomock.Any(), gomock.Any()).Return(nil),
					msa.EXPECT().AskOne(gomock.Any(), gomock.Any()).Return(nil),
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

			got, err := GuideCode(tc.args.surveyAsker, tc.args.workdir, tc.args.fileMap)

			require.ErrorIs(t, err, tc.targetErr)
			require.Equal(t, tc.want, got)
		})
	}
}

type fakeCodeSurveyAsker struct {
	TB         testing.TB
	Answers    CodeAnswers
	AskPerCall []survey.Prompt
	callCount  int
}

func (f *fakeCodeSurveyAsker) Ask(_ []*survey.Question, _ any, _ ...survey.AskOpt) error {
	return nil
}

func (f *fakeCodeSurveyAsker) AskOne(prompt survey.Prompt, response any, _ ...survey.AskOpt) error {
	f.callCount++

	ans := response.(*string)
	require.Empty(f.TB, ans)

	switch f.callCount {
	case 1:
		require.IsType(f.TB, &survey.Input{}, prompt)
		require.Equal(f.TB, f.AskPerCall[f.callCount-1], prompt)

		*ans = f.Answers.StackName
	case 2:
		require.IsType(f.TB, &survey.Select{}, prompt)
		require.Equal(f.TB, f.AskPerCall[f.callCount-1], prompt)

		*ans = f.Answers.Config
	case 3:
		require.IsType(f.TB, &survey.Input{}, prompt)
		require.Equal(f.TB, f.AskPerCall[f.callCount-1], prompt)

		*ans = f.Answers.Output
	}

	return nil
}
