package repository

import (
	"context"
	"time"

	"realtime-poll/internal/db"
	"realtime-poll/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func userCol() *mongo.Collection {
	return db.DB.Collection("users")
}

func sessionCol() *mongo.Collection {
	return db.DB.Collection("sessions")
}

func CreateUser(u models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := userCol().InsertOne(ctx, u)
	return err
}

func FindUserByUsername(username string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	err := userCol().FindOne(ctx, bson.M{"username": username}).Decode(&user)
	return &user, err
}

func CreateSession(s models.Session) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := sessionCol().InsertOne(ctx, s)
	return err
}

func GetSession(id string) (*models.Session, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var s models.Session
	err := sessionCol().FindOne(ctx, bson.M{"_id": id}).Decode(&s)
	return &s, err
}

func DeleteSession(id string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sessionCol().DeleteOne(ctx, bson.M{"_id": id})
}
