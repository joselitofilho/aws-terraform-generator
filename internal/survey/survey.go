//go:generate mockgen -package survey -destination survey_mock.go -source survey.go

package survey

import "github.com/AlecAivazis/survey/v2"

// SurveyAsker defines the interface for asking survey questions.
type SurveyAsker interface {
	Ask(questions []*survey.Question, answers interface{}, opts ...survey.AskOpt) error
	AskOne(p survey.Prompt, response interface{}, opts ...survey.AskOpt) error
}

// RealSurveyAsker is an implementation of SurveyAsker that uses the actual survey library.
type RealSurveyAsker struct{}

// Ask asks survey questions using the actual survey library.
func (r *RealSurveyAsker) Ask(questions []*survey.Question, answers interface{}, opts ...survey.AskOpt) error {
	return survey.Ask(questions, answers, opts...)
}

// AskOne asks a single survey question using the actual survey library.
func (r *RealSurveyAsker) AskOne(prompt survey.Prompt, response interface{}, opts ...survey.AskOpt) error {
	return survey.AskOne(prompt, response, opts...)
}
