package handler

import (
	"log/slog"
	"net/http"

	"bigfive-web/internal/view"

	"github.com/alexedwards/scs/v2"
)

type GlobalState struct {
	Count int
}

var (
	globalState    GlobalState
	SessionManager *scs.SessionManager
)

func GetCountHandler(w http.ResponseWriter, r *http.Request) {
	userCount := SessionManager.GetInt(r.Context(), "count")

	c := view.CountPage(globalState.Count, userCount)
	c.Render(r.Context(), w)
}

func PostCountHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.ErrorContext(r.Context(), "failed to parse form", "error", err)
		return
	}

	if r.Form.Has("global") {
		globalState.Count++
	}

	if r.Form.Has("user") {
		currentCount := SessionManager.GetInt(r.Context(), "count")
		SessionManager.Put(r.Context(), "count", currentCount+1)
	}

	GetCountHandler(w, r)
}
