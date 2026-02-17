package services

import (
	"context"
	"time"

	"realtime-poll/internal/repository"
	"realtime-poll/internal/ws"
)

func HandlePollExpired(pollID string) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// mark closed
	_ = repository.ClosePoll(ctx, pollID)

	// build final results
	svc := VoteService{}
	results, err := svc.GetResults(ctx, pollID)
	if err != nil {
		return
	}

	// broadcast final projection
	ws.EmitClosed(pollID, results)
}
