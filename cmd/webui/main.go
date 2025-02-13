package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"bigfive-web/internal/app"
	"bigfive-web/internal/infra"
	"bigfive-web/internal/web"
	"bigfive-web/internal/web/handler"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		slog.Error("failed to load .env file", "error", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	endpoint := os.Getenv("MONGODB_ENDPOINT")
	if endpoint == "" {
		slog.Error("MONGODB_ENDPOINT environment variable is not set")
		panic("MONGODB_ENDPOINT environment variable is required")
	}

	opts := options.Client().
		ApplyURI(endpoint).
		SetServerSelectionTimeout(5 * time.Second)

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		slog.Error("failed to create MongoDB client", "error", err)
		panic(err)
	}
	defer client.Disconnect(ctx)

	err = client.Ping(ctx, nil)
	if err != nil {
		slog.Error("failed to ping MongoDB", "error", err)
		panic(err)
	}

	slog.Info("successfully connected to MongoDB")

	database := client.Database(os.Getenv("MONGODB_DATABASE"))
	collection := database.Collection("answers")

	db := infra.NewPersonalityTestMongoDB(database, collection)
	svc := app.NewPersonalityTestService(db)
	handler := handler.NewPersonalityTestHandler(svc)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)
	web.RegisterRoutes(mux)

	err = http.ListenAndServe(":"+os.Getenv("PORT"), mux)
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
