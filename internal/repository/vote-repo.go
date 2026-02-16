package repository

import (
	"context"
	"time"
	"errors"

	"realtime-poll/internal/db"
	"realtime-poll/internal/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func voteCollection() *mongo.Collection {
	return db.DB.Collection("votes")
}

func GetVoteForViewer(ctx context.Context, pollID, userID, sessionID string) (*models.VoteRecord, error) {

	filter := bson.M{"poll_id": pollID}

	if userID != "" {
		filter["user_id"] = userID
	} else {
		filter["session_id"] = sessionID
	}

	var vote models.VoteRecord
	err := voteCollection().FindOne(ctx, filter).Decode(&vote)
	if err != nil {
		return nil, err
	}

	return &vote, nil
}


func FindExistingVote(ctx context.Context, pollID, userID, sessionID string) (*models.Vote, error) {

	filter := bson.M{
		"poll_id": pollID,
		"$or": []bson.M{
			{"user_id": userID},
			{"session_id": sessionID},
		},
	}

	var vote models.Vote
	err := voteCollection().FindOne(ctx, filter).Decode(&vote)
	if err != nil {
		return nil, nil
	}
	return &vote, nil
}

func InsertVote(ctx context.Context, vote models.Vote) error {
	vote.CreatedAt = time.Now()
	_, err := voteCollection().InsertOne(ctx, vote)
	return err
}

func UpdateVote(ctx context.Context, id primitive.ObjectID, optionID string) error {
	_, err := voteCollection().UpdateByID(ctx, id, bson.M{
		"$set": bson.M{"option_id": optionID},
	})
	return err
}

func CountVotesByUser(ctx context.Context, pollID, userID, sessionID, ip string) (int64, error) {

	filter := bson.M{"poll_id": pollID}

	conditions := []bson.M{}

	if userID != "" {
		conditions = append(conditions, bson.M{"user_id": userID})
	}

	if sessionID != "" {
		conditions = append(conditions, bson.M{"session_id": sessionID})
	}

	if ip != "" {
		conditions = append(conditions, bson.M{"ip": ip})
	}

	if len(conditions) > 0 {
		filter["$or"] = conditions
	}

	return voteCollection().CountDocuments(ctx, filter)
}

func IncrementOptionVote(ctx context.Context, pollID, optionID string, delta int) error {

	filter := bson.M{
		"poll_id": pollID,
		"content.options.option_id": optionID,
	}

	update := bson.M{
		"$inc": bson.M{
			"content.options.$.votes": delta,
			"meta.total_votes":        delta,
		},
	}

	res, err := getCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return errors.New("option not found")
	}

	return nil
}



func InsertVoteAtomic(ctx context.Context, vote models.VoteRecord) error {
	_, err := db.DB.Collection("votes").InsertOne(ctx, vote)
	return err
}

func IsDuplicateKey(err error) bool {
	if err == nil {
		return false
	}
	return mongo.IsDuplicateKeyError(err)
}

func UpdateVoteOption(ctx context.Context, voteID, optionID string) error {
	_, err := db.DB.Collection("votes").UpdateOne(ctx,
		bson.M{"vote_id": voteID},
		bson.M{"$set": bson.M{
			"option_id":  optionID,
			"updated_at": time.Now(),
		}},
	)
	return err
}

func GetVoteByIdentity(ctx context.Context, pollID, identity string) (*models.VoteRecord, error) {

	var vote models.VoteRecord

	err := db.DB.Collection("votes").FindOne(ctx, bson.M{
		"poll_id":  pollID,
		"identity": identity,
	}).Decode(&vote)

	if err != nil {
		return nil, err
	}
	return &vote, nil
}
