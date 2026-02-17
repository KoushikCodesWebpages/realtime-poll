package ws

import (
	"context"
	"time"

	"realtime-poll/internal/repository"
	// "realtime-poll/internal/services"
	
)

func startExpiryTimer(r *Room) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	poll, err := repository.GetPollByID(ctx, r.PollID)
	if err != nil || poll == nil {
		return
	}

	// No expiry → nothing to schedule
	if poll.Meta.ExpiresAt == nil {
		return
	}

	expireAt := *poll.Meta.ExpiresAt
	now := time.Now()

	// Already expired (server restart case)
	if now.After(expireAt) {
		closePoll(r)
		return
	}

	wait := time.Until(expireAt)

	timer := time.NewTimer(wait)

	go func() {
		<-timer.C
		closePoll(r)
	}()
}

func closePoll(r *Room) {

	if r.closed {
		return
	}
	r.closed = true

	// notify application layer
	if OnPollExpired != nil {
		OnPollExpired(r.PollID)
	}
}
