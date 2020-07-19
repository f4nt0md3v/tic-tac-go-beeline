package models

type Request struct {
	Command  string `json:"command,omitempty"`
	GameInfo Game   `json:"gameInfo,omitempty"`
	Message  string `json:"message,omitempty"`
}

type Response struct {
	Code     int    `json:"code,omitempty"`
	Command  string `json:"command,omitempty"`
	GameInfo Game   `json:"gameInfo,omitempty"`
	Message  string `json:"message,omitempty"`
	Type     int    `json:"type,omitempty"`
	Error    string `json:"error,omitempty"`
}
