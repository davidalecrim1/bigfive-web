package web

import (
	"net/http"

	"bigfive-web/internal/web/components"

	"github.com/a-h/templ"
)

func RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("GET /", templ.Handler(components.HomePage()))
	mux.Handle("GET /results", templ.Handler(components.PersonalityTestGetResultPage()))
}
