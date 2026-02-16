package db

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func Connect() {
	uri := os.Getenv("MONGO_URI")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}

	DB = client.Database(os.Getenv("DB_NAME"))
	log.Println("MongoDB connected")
	ensureIndexes()

}
func ensureIndexes() {

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	votes := DB.Collection("votes")

	index := mongo.IndexModel{
		Keys: bson.D{
			{Key: "poll_id", Value: 1},
			{Key: "identity", Value: 1},
		},
		Options: options.Index().
			SetUnique(true).
			SetName("unique_vote_per_identity"),
	}

	_, err := votes.Indexes().CreateOne(ctx, index)
	if err != nil {
		log.Fatal("failed creating vote index:", err)
	}

	log.Println("Vote indexes ensured")
}
