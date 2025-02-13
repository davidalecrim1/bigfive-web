package app

import (
	"errors"

	"github.com/google/uuid"
)

var SupportedLanguages = []string{
	"en-us",
}

var ErrUnsupportedLanguage = errors.New("unsupported language provided")

type Choice struct {
	Text  string
	Score int
	Color int
}

type ChoiceKeyed string

const (
	ChoiceKeyedPlus  ChoiceKeyed = "plus"
	ChoiceKeyedMinus ChoiceKeyed = "minus"
)

type ChoiceList = map[ChoiceKeyed][]*Choice

type Question struct {
	Id     string
	Text   string
	Keyed  string
	Domain string
	Facet  int
}

type Items struct {
	*Question
	Num    int
	Choice []*Choice
}

type UserAnswers struct {
	Id     uuid.UUID
	Domain string
	Score  int
	Count  int
	Result string
}

func NewUserAnswers(domain string, score, count int, result string) UserAnswers {
	return UserAnswers{
		Id:     uuid.New(),
		Domain: domain,
		Score:  score,
		Count:  count,
		Result: result,
	}
}
