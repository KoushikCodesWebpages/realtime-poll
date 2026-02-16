package repository

import (
	"context"
	"time"

	"realtime-poll/internal/db"
	"realtime-poll/internal/models"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	// "github.com/google/uuid"

)




func sessionCol() *mongo.Collection {
	return db.DB.Collection("sessions")
}


func CreateAnonymousSession(ctx context.Context, ip string) (*models.Session, error) {

	session := &models.Session{
		SessionID: uuid.NewString(),
		UserID:    "", // guest
		IP:        ip,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour), // 30 days
	}

	_, err := sessionCol().InsertOne(ctx, session)
	if err != nil {
		return nil, err
	}

	return session, nil
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

err := sessionCol().FindOne(ctx, bson.M{"session_id": id}).Decode(&s)

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
