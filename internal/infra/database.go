package infra

import (
	"context"
	"fmt"
	"time"

	"bigfive-web/internal/app"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PersonalityTestMongoDB struct {
	db         *mongo.Database
	collection *mongo.Collection
}

func NewPersonalityTestMongoDB(db *mongo.Database, collection *mongo.Collection) *PersonalityTestMongoDB {
	return &PersonalityTestMongoDB{
		db:         db,
		collection: collection,
	}
}

func (p *PersonalityTestMongoDB) SaveTestResults(ctx context.Context, answers []app.UserAnswers) (string, error) {
	document := bson.M{
		"answers":    answers,
		"created_at": time.Now(),
	}

	insertResult, err := p.collection.InsertOne(ctx, document)
	if err != nil {
		return "", fmt.Errorf("failed to insert document: %w", err)
	}

	insertedID, ok := insertResult.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("failed to convert InsertedID to ObjectID")
	}

	return insertedID.Hex(), nil
}

func (p *PersonalityTestMongoDB) GetTestResults(ctx context.Context, id string) ([]app.UserAnswers, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ObjectID: %w", err)
	}

	var result struct {
		Answers []app.UserAnswers `bson:"answers"`
	}
	err = p.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("failed to find document: %w", err)
	}

	return result.Answers, nil
}
