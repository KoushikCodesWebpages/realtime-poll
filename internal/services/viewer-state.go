package services

import (
	"context"
	"time"

	"realtime-poll/internal/models"
	"realtime-poll/internal/repository"
)

type ViewerState struct {
    AlreadyVoted   bool
    CanVote        bool
    CanViewResults bool
    Changed        bool
    VotedOptionID  string
    Started        bool
    Ended          bool
}

func BuildViewerState(
    ctx context.Context,
    poll *models.Poll,
    userID string,
    sessionID string,
) (*ViewerState, error) {

    now := time.Now()

    started := poll.Behavior.StartAt == nil || now.After(*poll.Behavior.StartAt)
    ended   := poll.State.IsClosed || (poll.Behavior.EndAt != nil && now.After(*poll.Behavior.EndAt))

    viewer := &ViewerState{
        Started: started,
        Ended:   ended,
    }

    vote, _ := repository.FindVote(ctx, poll.PollID, userID, sessionID)

    if vote != nil {
        viewer.AlreadyVoted  = true
        viewer.VotedOptionID = vote.OptionID
    }

    // Can view results
    viewer.CanViewResults =
        poll.Behavior.ShowLiveResults ||
        ended ||
        viewer.AlreadyVoted

    // Can vote logic
    if !started || ended {
        viewer.CanVote = false
        return viewer, nil
    }

    if viewer.AlreadyVoted {
        if poll.Vote.AllowChangeVote {
            viewer.CanVote = true
            viewer.Changed = true
        } else {
            viewer.CanVote = false
        }
        return viewer, nil
    }

    viewer.CanVote = true
    return viewer, nil
}

