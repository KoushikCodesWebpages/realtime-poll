package repository

import (
	"context"
	"realtime-poll/internal/db"
	"realtime-poll/internal/models"
	"realtime-poll/internal/dto"
	"time"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
func GetPollsByOwner(ctx context.Context, ownerID string) ([]models.Poll, error) {

	cursor, err := getCollection().Find(
		ctx,
		bson.M{"owner_id": ownerID},
		options.Find().SetSort(bson.M{"meta.created_at": -1}),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var polls []models.Poll
	if err := cursor.All(ctx, &polls); err != nil {
		return nil, err
	}

	return polls, nil
}

func GetPollsByOwnerPaginated(
	ctx context.Context,
	ownerID string,
	f dto.PollFilter,
	cursor *time.Time,
	direction string,
) ([]models.Poll, error) {

	// 1️⃣ Build base filter
	query := buildOwnerFilter(ownerID, f)

	// 2️⃣ Cursor pagination
	if cursor != nil {
		op := "$lt"
		if direction == "prev" {
			op = "$gt"
		}

		query["meta.created_at"] = bson.M{op: *cursor}
	}

	// 3️⃣ Sorting logic (THIS is where it goes)
	sort := bson.M{"meta.created_at": -1}

	switch f.Sort {
	case "oldest":
		sort = bson.M{"meta.created_at": 1}

	case "most_voted":
		sort = bson.M{"meta.total_votes": -1}
	}

	// reverse direction for prev navigation
	if direction == "prev" {
		for k, v := range sort {
			sort[k] = -v.(int)
		}
	}

	opts := options.Find().
		SetSort(sort).
		SetLimit(f.Limit)

	cursorDB, err := getCollection().Find(ctx, query, opts)
	if err != nil {
		return nil, err
	}
	defer cursorDB.Close(ctx)

	var polls []models.Poll
	if err := cursorDB.All(ctx, &polls); err != nil {
		return nil, err
	}

	return polls, nil
}

func buildOwnerFilter(ownerID string, f dto.PollFilter) bson.M {

	filter := bson.M{"owner_id": ownerID}

	now := time.Now()

	// status
	switch f.Status {
	case "active":
		filter["state.is_closed"] = false
		filter["meta.expires_at"] = bson.M{"$gt": now}
	case "scheduled":
		filter["behavior.start_at"] = bson.M{"$gt": now}
	case "closed":
		filter["state.is_closed"] = true
	case "expired":
		filter["meta.expires_at"] = bson.M{"$lt": now}
	}

	// visibility
	if f.Visibility != "" {
		filter["access.visibility"] = f.Visibility
	}

	// search
	if f.Search != "" {
		filter["content.question"] = bson.M{
			"$regex":   f.Search,
			"$options": "i",
		}
	}

	// votes
	if f.HasVotes != nil {
		if *f.HasVotes {
			filter["meta.total_votes"] = bson.M{"$gt": 0}
		} else {
			filter["meta.total_votes"] = 0
		}
	}

	// date range
	if f.DateFrom != nil || f.DateTo != nil {
		dateFilter := bson.M{}
		if f.DateFrom != nil {
			dateFilter["$gte"] = *f.DateFrom
		}
		if f.DateTo != nil {
			dateFilter["$lte"] = *f.DateTo
		}
		filter["meta.created_at"] = dateFilter
	}

	return filter
}
func CountPollsByOwner(ctx context.Context, ownerID string) (int64, error) {
	return getCollection().CountDocuments(ctx, bson.M{"owner_id": ownerID})
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
