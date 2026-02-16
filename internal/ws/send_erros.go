package ws

import (
	"encoding/json"
	"realtime-poll/internal/apperror"
)

type WSResponse struct {
	Type    string      `json:"type"`
	Success bool        `json:"success"`
	Code    string      `json:"code,omitempty"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func SendError(conn *Client, err error) {

	var res WSResponse

	if appErr, ok := err.(*apperror.AppError); ok {
		res = WSResponse{
			Type:    "error",
			Success: false,
			Code:    appErr.Code,
			Message: appErr.Message,
		}
	} else {
		res = WSResponse{
			Type:    "error",
			Success: false,
			Code:    apperror.INTERNAL_ERROR,
			Message: "Unexpected error",
		}
	}

	payload, _ := json.Marshal(res)
	conn.Send(payload)
}
