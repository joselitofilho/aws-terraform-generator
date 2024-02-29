//go:generate mockgen -package survey -destination survey_mock.go -source survey.go

package survey

import "github.com/AlecAivazis/survey/v2"

// Asker defines the interface for asking survey questions.
type Asker interface {
	Ask(questions []*survey.Question, answers any, opts ...survey.AskOpt) error
	AskOne(p survey.Prompt, response any, opts ...survey.AskOpt) error
}

// RealAsker is an implementation of SurveyAsker that uses the actual survey library.
type RealAsker struct{}

// Ask asks survey questions using the actual survey library.
func (r *RealAsker) Ask(questions []*survey.Question, answers any, opts ...survey.AskOpt) error {
	return survey.Ask(questions, answers, opts...)
}

// AskOne asks a single survey question using the actual survey library.
func (r *RealAsker) AskOne(prompt survey.Prompt, response any, opts ...survey.AskOpt) error {
	return survey.AskOne(prompt, response, opts...)
}
