package repository

import (
	"context"
	"realtime-poll/internal/db"
	"realtime-poll/internal/models"
	"time"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

func IsExpired(p models.Poll) bool {
	if p.Behavior.EndAt == nil {
		return false
	}
	return time.Now().After(*p.Behavior.EndAt)
}

func getCollection() *mongo.Collection {
	return db.DB.Collection("polls")
}

func CreatePoll(ctx context.Context, poll *models.Poll) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := getCollection().InsertOne(ctx, poll)
	return err
}

func GetPoll(id string) (*models.Poll, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var poll models.Poll
	err := getCollection().FindOne(ctx, bson.M{"_id": id}).Decode(&poll)
	return &poll, err
}

func IncrementVote(pollID, optionID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := getCollection().UpdateOne(ctx,
		bson.M{"_id": pollID, "options.id": optionID},
		bson.M{"$inc": bson.M{"options.$.votes": 1}},
	)
	return err
}
