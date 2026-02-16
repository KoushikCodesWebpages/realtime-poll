package apperror

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
