package models

type Request struct {
	Command  string `json:"command,omitempty"`
	GameInfo Game   `json:"game_info,omitempty"`
	Message  string `json:"message,omitempty"`
}

type Response struct {
	Code     int    `json:"code,omitempty"`
	Command  string `json:"command,omitempty"`
	GameInfo Game   `json:"game_info,omitempty"`
	Message  string `json:"message,omitempty"`
}

type ErrorResponse struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}
