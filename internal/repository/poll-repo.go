package repository

import (
	"context"
	"errors"
	"time"
	"strings"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"

	"realtime-poll/internal/db"
	"realtime-poll/internal/models"
	"realtime-poll/internal/dto"
)

func getCollection() *mongo.Collection {
	return db.DB.Collection("polls")
}

func PollHasVotes(ctx context.Context, pollID string) (bool, error) {

	count, err := getCollection().CountDocuments(ctx, bson.M{
		"poll_id": pollID,
	})

	if err != nil {
		return false, err
	}

	return count > 0, nil
}


func GetPollByIDRaw(ctx context.Context, pollID string) (*models.Poll, error) {

	var poll models.Poll

	err := getCollection().
		FindOne(ctx, bson.M{"poll_id": pollID}).
		Decode(&poll)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &poll, nil
}

func GetPollByID(ctx context.Context, pollID string) (*models.Poll, error) {

	var poll models.Poll

	err := getCollection().FindOne(
		ctx,
		bson.M{
			"poll_id": pollID,
			"state.is_deleted": bson.M{"$ne": true},
		},
	).Decode(&poll)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &poll, nil
}
func HardDeletePoll(ctx context.Context, pollID string) error {

	res, err := getCollection().DeleteOne(
		ctx,
		bson.M{
			"poll_id": pollID,
			"state.is_deleted": true,
		},
	)

	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("poll must be soft deleted first")
	}

	return nil
}

func SoftDeletePoll(ctx context.Context, pollID string) error {

	res, err := getCollection().UpdateOne(
		ctx,
		bson.M{
			"poll_id": pollID,
			"state.is_deleted": false,
		},
		bson.M{
			"$set": bson.M{
				"state.is_deleted": true,
				"meta.updated_at": time.Now(),
			},
		},
	)

	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}
func UpdatePollFields(ctx context.Context, pollID string, update bson.M) error {

	res, err := getCollection().UpdateOne(
		ctx,
		bson.M{"poll_id": pollID, "state.is_deleted": false},
		bson.M{"$set": update, "$inc": bson.M{"state.version": 1}},
	)

	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func IsExpired(p models.Poll) bool {
	if p.Behavior.EndAt == nil {
		return false
	}
	return time.Now().After(*p.Behavior.EndAt)
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
		bson.M{
			"owner_id": ownerID,
			"state.is_deleted": bson.M{"$ne": true},
		},
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

	filter := bson.M{
	"owner_id": ownerID,
	"state.is_deleted": bson.M{"$ne": true},
	}

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

func IncrementVote(pollID, optionID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := getCollection().UpdateOne(ctx,
		bson.M{"_id": pollID, "options.id": optionID},
		bson.M{"$inc": bson.M{"options.$.votes": 1}},
	)
	return err
}


func ClosePoll(ctx context.Context, pollID string) error {

	_, err := getCollection().UpdateOne(
		ctx,
		bson.M{"poll_id": pollID},
		bson.M{
			"$set": bson.M{
				"state.is_closed": true,
			},
		},
	)

	return err
}

func UpdatePollTotalVotes(ctx context.Context, pollID string, total int64) error {
	_, err := getCollection().UpdateOne(ctx,
		bson.M{"_id": pollID},
		bson.M{"$set": bson.M{"meta.total_votes": total}},
	)
	return err
}

func UpdatePollsTotalVotesBulk(ctx context.Context, updates map[string]int64) error {

	if len(updates) == 0 {
		return nil
	}

	models := make([]mongo.WriteModel, 0, len(updates))

	for id, total := range updates {
		models = append(models, mongo.NewUpdateOneModel().
			SetFilter(bson.M{"_id": id}).
			SetUpdate(bson.M{"$set": bson.M{"meta.total_votes": total}}),
		)
	}

	_, err := getCollection().BulkWrite(ctx, models)
	return err
}


func ReplacePoll(ctx context.Context, pollID string, poll *models.Poll) error {
	_, err := getCollection().ReplaceOne(
		ctx,
		bson.M{"poll_id": pollID},
		poll,
	)
	return err
}


func ApplyEdit(poll *models.Poll, req dto.EditPollReq) bool {

	changed := false

	// =====================================================
	// CONTENT
	// =====================================================
	if req.Content != nil {

		if req.Content.Question != nil &&
			strings.TrimSpace(*req.Content.Question) != poll.Content.Question {

			poll.Content.Question = strings.TrimSpace(*req.Content.Question)
			changed = true
		}

		if req.Content.Description != nil &&
			strings.TrimSpace(*req.Content.Description) != poll.Content.Description {

			poll.Content.Description = strings.TrimSpace(*req.Content.Description)
			changed = true
		}

		// ----- OPTIONS (CRITICAL PART) -----
		if req.Content.Options != nil {

			oldOptions := poll.Content.Options
			newOptions := []models.Option{}

			// map existing options by id
			oldMap := map[string]models.Option{}
			for _, o := range oldOptions {
				oldMap[o.OptionID] = o
			}

			for _, incoming := range req.Content.Options {

				text := strings.TrimSpace(incoming.Text)
				if text == "" {
					continue
				}

				// existing option (keep votes)
				if incoming.OptionID != "" {
					if old, ok := oldMap[incoming.OptionID]; ok {

						if old.Text != text {
							old.Text = text
							changed = true
						}

						newOptions = append(newOptions, old)
						continue
					}
				}

				// new option
				newOptions = append(newOptions, models.Option{
					OptionID: uuid.NewString(),
					Text:     text,
					Votes:    0,
				})
				changed = true
			}

			// detect removed options
			if len(newOptions) != len(oldOptions) {
				changed = true
			}

			poll.Content.Options = newOptions
		}

		if req.Content.AllowCustomOption != nil &&
			*req.Content.AllowCustomOption != poll.Content.AllowCustom {

			poll.Content.AllowCustom = *req.Content.AllowCustomOption
			changed = true
		}

		if req.Content.RandomizeOptions != nil &&
			*req.Content.RandomizeOptions != poll.Content.Randomize {

			poll.Content.Randomize = *req.Content.RandomizeOptions
			changed = true
		}
	}

	// =====================================================
	// ACCESS
	// =====================================================
	if req.Access != nil {

		if req.Access.Visibility != nil &&
			*req.Access.Visibility != poll.Access.Visibility {

			poll.Access.Visibility = *req.Access.Visibility
			poll.Access.RequireLogin = *req.Access.Visibility != "public"
			changed = true
		}

		if req.Access.AllowedEmails != nil {
			poll.Access.AllowedEmails = req.Access.AllowedEmails
			changed = true
		}
	}

	// =====================================================
	// VOTE
	// =====================================================
	if req.Vote != nil {

		setBool := func(ptr *bool, target *bool) {
			if ptr != nil && *ptr != *target {
				*target = *ptr
				changed = true
			}
		}

		setBool(req.Vote.AllowChangeVote, &poll.Vote.AllowChangeVote)
		setBool(req.Vote.AnonymousVote, &poll.Vote.AnonymousVote)
		setBool(req.Vote.HideResultsUntilEnd, &poll.Vote.HideResults)
		setBool(req.Vote.ShowVoters, &poll.Vote.ShowVoters)
		setBool(req.Vote.UniqueIP, &poll.Vote.UniqueIP)
		setBool(req.Vote.UniqueSession, &poll.Vote.UniqueSession)
	}

	// =====================================================
	// BEHAVIOR
	// =====================================================
	if req.Behavior != nil {

		if req.Behavior.StartAt != nil {
			poll.Behavior.StartAt = req.Behavior.StartAt
			changed = true
		}

		if req.Behavior.EndAt != nil {
			poll.Behavior.EndAt = req.Behavior.EndAt
			poll.Meta.ExpiresAt = req.Behavior.EndAt
			changed = true
		}

		if req.Behavior.AutoClose != nil &&
			*req.Behavior.AutoClose != poll.Behavior.AutoClose {

			poll.Behavior.AutoClose = *req.Behavior.AutoClose
			changed = true
		}

		if req.Behavior.ShowLiveResults != nil &&
			*req.Behavior.ShowLiveResults != poll.Behavior.ShowLiveResults {

			poll.Behavior.ShowLiveResults = *req.Behavior.ShowLiveResults
			changed = true
		}
	}

	return changed
}

