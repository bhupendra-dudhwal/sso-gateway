package models

type Response struct {
	StatusCode int    `json:"status_code"`
	Status     bool   `json:"status"`
	RequestID  string `json:"request_id"`
	Message    string `json:"message"`
	Error      *Error `json:"error,omitempty"`
	Payload    any    `json:"payload,omitempty"`
}

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Detail  any    `json:"detail,omitempty"`
}
