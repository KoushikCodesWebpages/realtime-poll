package apperror

func TokenInvalid() *AppError {
	return &AppError{
		Code:    TOKEN_INVALID,
		Message: "Invalid or corrupted link",
	}
}

func TokenExpired() *AppError {
	return &AppError{
		Code:    TOKEN_EXPIRED,
		Message: "This link has expired",
	}
}


func Unauthorized() *AppError {
	return &AppError{
		Code:    AUTH_UNAUTHORIZED,
		Message: "You are not authorized",
		Action:  "LOGOUT",
	}
}

func InvalidCredentials() *AppError {
	return &AppError{
		Code:    AUTH_INVALID_CREDENTIALS,
		Message: "Invalid username or password",
	}
}

func PollVotingStarted() *AppError {
	return &AppError{
		Code:    POLL_VOTING_STARTED,
		Message: "Poll already started",
	}
}

func Internal() *AppError {
	return &AppError{
		Code:    INTERNAL_ERROR,
		Message: "Something went wrong",
	}
}

func PollNotStarted() *AppError {
	return &AppError{
		Code:    POLL_NOT_STARTED,
		Message: "Voting has not started yet",
	}
}


func PollNotFound() *AppError {
	return &AppError{
		Code:    POLL_NOT_FOUND,
		Message: "Voting has not started yet",
	}
}

func PollEnded() *AppError {
	return &AppError{
		Code:    POLL_ENDED,
		Message: "Voting has ended",
	}
}

func PollClosed() *AppError {
	return &AppError{
		Code:    POLL_CLOSED,
		Message: "Poll is closed",
	}
}

func AlreadyVoted() *AppError {
	return &AppError{
		Code:    POLL_ALREADY_VOTED,
		Message: "You already voted",
	}
}
func BadRequest(msg string) *AppError {
	return &AppError{
		Code:    VALIDATION_FAILED,
		Message: msg,
	}
}

func Validation(message string) *AppError {
	return &AppError{
		Code:    VALIDATION_FAILED,
		Message: message,
	}
}

