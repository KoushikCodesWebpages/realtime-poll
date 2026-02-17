package ws

type Role int

const (
	RoleViewer Role = iota
	RoleVoter
	RoleOwner
)

type Visibility int

const (
	ViewHidden Visibility = iota
	ViewLive
	ViewOwner
	ViewFinal
)

type Session struct {
	ConnID string

	UserID    *string
	Anonymous string
	TokenSub  string // <--- NEW (JWT subject)

	PollID string

	Role       Role
	Visibility Visibility
	IsOwner bool 
}

func (s *Session) IdentityKey() string {

	// strongest priority: websocket token
	if s.TokenSub != "" {
		return "t:" + s.TokenSub
	}

	if s.UserID != nil {
		return "u:" + *s.UserID
	}

	if s.Anonymous != "" {
		return "a:" + s.Anonymous
	}

	return s.ConnID
}

