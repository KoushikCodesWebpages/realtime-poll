package ws

import "realtime-poll/internal/models"

func assignRole(s *Session, ownerID string) {

	if s.UserID != nil && *s.UserID == ownerID {
		s.Role = RoleOwner
		return
	}

	if s.UserID != nil {
		s.Role = RoleVoter
		return
	}

	s.Role = RoleViewer
}

func assignVisibility(s *Session, poll *models.Poll) {

	// poll closed overrides everything
	if poll.State.IsClosed {
		s.Visibility = ViewFinal
		return
	}

	// owner always live
	if s.Role == RoleOwner {
		s.Visibility = ViewOwner
		return
	}

	hide := poll.Vote.HideResults
	live := poll.Behavior.ShowLiveResults

	switch {
	case hide && !live:
		s.Visibility = ViewHidden

	case !hide && live:
		s.Visibility = ViewLive

	default:
		s.Visibility = ViewHidden
	}
}
