package app

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"slices"
)

const DEFAULT_QUESTIONS_LENGTH = 120

type PersonalityTestDatabase interface {
	SaveTestResults(ctx context.Context, answers []UserAnswers) (string, error)
	GetTestResults(ctx context.Context, id string) ([]UserAnswers, error)
}

type PersonalityTestService struct {
	db PersonalityTestDatabase
}

type svc = PersonalityTestService

func NewPersonalityTestService(db PersonalityTestDatabase) *svc {
	return &svc{
		db: db,
	}
}

func (p *svc) GetItems(language string) ([]Items, error) {
	items := make([]Items, 0, DEFAULT_QUESTIONS_LENGTH)

	if !slices.Contains(SupportedLanguages, language) {
		return nil, ErrUnsupportedLanguage
	}

	questions, err := loadQuestions(language)
	if err != nil {
		return nil, err
	}

	choices, err := loadChoices(language)
	if err != nil {
		return nil, err
	}

	for i, question := range questions {
		num := i + 1

		items = append(items, Items{
			Question: question,
			Num:      num,
			Choice:   choices[ChoiceKeyed(question.Keyed)],
		})
	}

	return items, nil
}

func loadQuestions(language string) ([]*Question, error) {
	filePath := fmt.Sprintf("internal/app/questions/%s.json", language)
	questions := make([]*Question, 0, DEFAULT_QUESTIONS_LENGTH)

	byteValue, err := readFile(filePath)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(byteValue, &questions)
	if err != nil {
		return nil, err
	}

	return questions, nil
}

func readFile(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return byteValue, nil
}

func loadChoices(language string) (ChoiceList, error) {
	filePath := fmt.Sprintf("internal/app/choices/%s.json", language)
	choices := make(ChoiceList)

	byteValue, err := readFile(filePath)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(byteValue, &choices)
	if err != nil {
		return nil, err
	}

	return choices, nil
}

func (p *svc) SaveTestResults(ctx context.Context, language string, selectedChoices []int) (string, error) {
	items, err := p.GetItems(language)
	if err != nil {
		return "", err
	}

	result := make(map[string]int, DEFAULT_QUESTIONS_LENGTH)

	for i, item := range items {
		choice := item.Choice[selectedChoices[i]]
		savedScore, ok := result[item.Domain]

		if ok {
			result[item.Domain] = choice.Score + savedScore
		} else {
			result[item.Domain] = choice.Score
		}
	}

	var answers []UserAnswers
	for key, value := range result {
		questionsCount := p.questionsPerDomain(key, items)
		answer := NewUserAnswers(
			key, value, questionsCount, p.getResultFromDomain(value, questionsCount),
		)
		answers = append(answers, answer)
	}

	id, err := p.db.SaveTestResults(ctx, answers)
	if err != nil {
		return "", err
	}

	return id, err
}

func (p *svc) getResultFromDomain(score, count int) string {
	avg := float64(score) / float64(count)

	if avg > 3.5 {
		return "high"
	}
	if avg < 2.5 {
		return "low"
	}

	return "neutral"
}

func (p *svc) questionsPerDomain(domain string, items []Items) int {
	count := 0

	for _, item := range items {
		if item.Domain == domain {
			count++
		}
	}

	return count
}

func (p *svc) GetTestResults(ctx context.Context, id string) ([]UserAnswers, error) {
	return p.db.GetTestResults(ctx, id)
}
