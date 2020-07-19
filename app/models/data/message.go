package data

import "github.com/f4nt0md3v/tic-tac-go-beeline/app/models/game"

type Request struct {
	Command  string    `json:"command,omitempty"`
	GameInfo game.Game `json:"gameInfo,omitempty"`
	Message  string    `json:"message,omitempty"`
}

type TypeMessage int

const (
	Broadcast TypeMessage = iota
	Single
)

type Response struct {
	Code        int         `json:"code,omitempty"`
	Command     string      `json:"command,omitempty"`
	GameInfo    *game.Game  `json:"gameInfo,omitempty"`
	Message     string      `json:"message,omitempty"`
	MessageType TypeMessage `json:"type,omitempty"`
	Error       string      `json:"error,omitempty"`
}
