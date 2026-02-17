package ws

// Called when poll expires
var OnPollExpired func(pollID string)

// Called when a client joins a room
var BuildPollSnapshot func(sess *Session) (any, error)
