package apperror

type Code string

const (
	TokenMalformed Code = "TOKEN_MALFORMED"
	TokenInvalid   Code = "TOKEN_INVALID"
	TokenExpired   Code = "TOKEN_EXPIRED"
)

type Error struct {
	Code    Code   `json:"code"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return string(e.Code)
}

func New(code Code, msg string) *Error {
	return &Error{Code: code, Message: msg}
}
