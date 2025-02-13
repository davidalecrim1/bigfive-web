package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"bigfive-web/internal/handler"
	"bigfive-web/internal/view"

	"github.com/a-h/templ"
	"github.com/alexedwards/scs/v2"
	"github.com/joho/godotenv"
)

func main() {
	sessionManager := scs.New()
	sessionManager.Lifetime = time.Hour * 24
	handler.SessionManager = sessionManager

	err := godotenv.Load(".env")
	if err != nil {
		slog.Error("failed to load .env file", "error", err)
	}

	mux := http.NewServeMux()

	mux.Handle("GET /", templ.Handler(view.Home()))
	mux.HandleFunc("GET /count", handler.GetCountHandler)
	mux.HandleFunc("POST /count", handler.PostCountHandler)

	muxWithSessionMiddleware := sessionManager.LoadAndSave(mux)

	err = http.ListenAndServe(":"+os.Getenv("PORT"), muxWithSessionMiddleware)
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
