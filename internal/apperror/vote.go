package apperror

const (
	// voting state
	VoteAlreadyCast    Code = "VOTE_ALREADY_CAST"
	VoteNotAllowed     Code = "VOTE_NOT_ALLOWED"
	VoteInvalidOption  Code = "VOTE_INVALID_OPTION"

	// lifecycle
	PollNotStarted Code = "POLL_NOT_STARTED"
	PollEnded      Code = "POLL_ENDED"
	PollClosed     Code = "POLL_CLOSED"

	// permissions
	LoginRequired Code = "LOGIN_REQUIRED"
	NotWhitelisted Code = "NOT_WHITELISTED"
)

func VoteAlready() *Error {
	return New(VoteAlreadyCast, "you have already voted")
}

func VoteDenied(msg string) *Error {
	return New(VoteNotAllowed, msg)
}

func PollStartedErr() *Error {
	return New(PollNotStarted, "poll has not started yet")
}

func PollEndedErr() *Error {
	return New(PollEnded, "poll has ended")
}

func PollClosedErr() *Error {
	return New(PollClosed, "poll is closed")
}
