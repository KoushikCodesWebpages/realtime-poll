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
	ConnID     string
	UserID     *string
	Anonymous  string
	PollID     string
	Role       Role
	Visibility Visibility
}
