package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"bigfive-web/internal/app"
	"bigfive-web/internal/web/components"
)

type PersonalityTestHandler struct {
	svc *app.PersonalityTestService
}

func NewPersonalityTestHandler(
	svc *app.PersonalityTestService,
) *PersonalityTestHandler {
	return &PersonalityTestHandler{
		svc: svc,
	}
}

func (h *PersonalityTestHandler) GetPersonalityTest(
	w http.ResponseWriter,
	r *http.Request,
) {
	// TODO: get language middleware for multiple supported languagens
	language := "en-us"

	items, err := h.svc.GetItems(language)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.ErrorContext(r.Context(), "GetPersonalityTest - failed to get questions items", "error", err, "items", items)
		return
	}

	c := components.PersonalityTest(items)
	c.Render(r.Context(), w)
}

func (h *PersonalityTestHandler) SaveAnswers(
	w http.ResponseWriter,
	r *http.Request,
) {
	// TODO: get language middleware for multiple supported languagens
	language := "en-us"

	selectedChoices := make([]int, 0, app.DEFAULT_QUESTIONS_LENGTH)

	for i := range app.DEFAULT_QUESTIONS_LENGTH {
		questionNum := i + 1

		selectedChoice := r.FormValue(fmt.Sprintf("question-%v", questionNum))
		selectedChoiceAsInt, _ := strconv.Atoi(selectedChoice)
		selectedChoices = append(selectedChoices, selectedChoiceAsInt)
	}

	id, err := h.svc.SaveTestResults(r.Context(), language, selectedChoices)
	if err != nil {
		slog.ErrorContext(r.Context(), "SaveAnswers - failed to save test results", "error", err, "selectedChoices", selectedChoices)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	slog.Info("saved test results", "id", id)

	w.Header().Set("Location", fmt.Sprintf("/results/%v", id))
	w.WriteHeader(http.StatusSeeOther)
}

func (h *PersonalityTestHandler) getResults(
	id string,
	w http.ResponseWriter,
	r *http.Request,
) {
	results, err := h.svc.GetTestResults(r.Context(), id)
	if err != nil {
		slog.ErrorContext(r.Context(), "failed to get results for id", "error", err, "id", id)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	c := components.PersonalityTestResultPage(results)
	c.Render(r.Context(), w)
}

func (h *PersonalityTestHandler) GetResultsFromForm(
	w http.ResponseWriter,
	r *http.Request,
) {
	id := r.FormValue("id")
	h.getResults(id, w, r)
}

func (h *PersonalityTestHandler) GetResultsFromPostRequest(
	w http.ResponseWriter,
	r *http.Request,
) {
	id := r.PathValue("id")
	h.getResults(id, w, r)
}

func (h *PersonalityTestHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /test", h.GetPersonalityTest)
	mux.HandleFunc("POST /submit", h.SaveAnswers)
	mux.HandleFunc("GET /results/submit", h.GetResultsFromForm)
	mux.HandleFunc("GET /results/{id}", h.GetResultsFromPostRequest)
}
