package app

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	// Change working directory to the project root.
	// Assuming your tests execute from "bigfive-web/internal/app",
	// go up two directory to get to "bigfive-web" root folder
	// This is needed because the questions and choices files are in the root directory
	if err := os.Chdir("../.."); err != nil {
		panic(fmt.Sprintf("failed to change working directory: %s", err))
	}

	os.Exit(m.Run())
}

func TestPersonalityTestService_GetItems(t *testing.T) {
	type args struct {
		language string
	}
	tests := []struct {
		name        string
		args        args
		want        []Items
		expectedErr error
	}{
		{
			name: "get questions for personality test with valid language",
			args: args{
				language: "en-us",
			},
			want:        []Items{},
			expectedErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := NewFakePersonalityTestRepository()
			svc := NewPersonalityTestService(db)
			got, err := svc.GetItems(tt.args.language)

			require.Equal(t, err, tt.expectedErr)
			require.IsType(t, []Items{}, got, "Expected result to be a slice")
			require.Greater(t, len(got), 0, "Expected slice length to be greater than 0")
		})
	}
}

type FakePersonalityTestDatabase struct {
	data   map[string][]UserAnswers
	nextID int // Simulate auto-incrementing IDs for inserted documents.
}

func NewFakePersonalityTestRepository() *FakePersonalityTestDatabase {
	return &FakePersonalityTestDatabase{
		data:   make(map[string][]UserAnswers),
		nextID: 1,
	}
}

func (f *FakePersonalityTestDatabase) SaveTestResults(ctx context.Context, answers []UserAnswers) (string, error) {
	id := fmt.Sprintf("%d", f.nextID)
	f.nextID++
	f.data[id] = answers

	return id, nil
}

func (f *FakePersonalityTestDatabase) GetTestResults(ctx context.Context, id string) ([]UserAnswers, error) {
	answers, exists := f.data[id]
	if !exists {
		return nil, fmt.Errorf("document with ID %s not found", id)
	}

	return answers, nil
}
