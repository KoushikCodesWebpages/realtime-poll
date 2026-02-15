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



func CreateUser(u models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := userCol().InsertOne(ctx, u)
	return err
}

func FindUserByEmail(email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	err := userCol().FindOne(ctx, bson.M{"email": email}).Decode(&user)

	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func FindUserByUsername(username string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	err := userCol().FindOne(ctx, bson.M{"username": username}).Decode(&user)
		if err == mongo.ErrNoDocuments {
		return nil, nil // IMPORTANT: user truly not found
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func FindUserByIdentifier(identifier string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User

	filter := bson.M{
		"$or": []bson.M{
			{"username": identifier},
			{"email": identifier},
		},
	}

	err := userCol().FindOne(ctx, filter).Decode(&user)
	return &user, err
}
func FindUserByAuthID(id string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var user models.User
	err := userCol().FindOne(ctx, bson.M{"auth_user_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}


func FindUserByID(id string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	err := userCol().FindOne(ctx, bson.M{"auth_user_id": id}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &user, err
}

