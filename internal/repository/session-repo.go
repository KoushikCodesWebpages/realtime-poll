package repository

import (
	"context"
	"time"

	"realtime-poll/internal/db"
	"realtime-poll/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	// "github.com/google/uuid"

)

func sessionCol() *mongo.Collection {
	return db.DB.Collection("sessions")
}


func CreateSession(s models.Session) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := sessionCol().InsertOne(ctx, bson.M{
		"_id":        s.SessionID,
		"user_id":    s.UserID,
		"created_at": s.CreatedAt,
		"expires_at": s.ExpiresAt,
	})

	return err
}


func GetSession(id string) (*models.Session, error) {
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

var s models.Session

err := sessionCol().FindOne(ctx, bson.M{"_id": id}).Decode(&s)

if err == mongo.ErrNoDocuments {
	return nil, nil
}
if err != nil {
	return nil, err
}

// expiry check (UTC safe)
if time.Now().UTC().After(s.ExpiresAt.UTC()) {
	return nil, nil
}

return &s, nil


}

func InsertSession(session *models.Session) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := sessionCol().InsertOne(ctx, session)
	return err
}

func DeleteSession(id string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sessionCol().DeleteOne(ctx, bson.M{"_id": id})
}
