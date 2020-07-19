package models

import "time"

type Game struct {
	ID             int        `json:"id,omitempty"`
	GameId         string     `json:"gameId,omitempty"`
	FirstUserId    string     `json:"firstUserId,omitempty"`
	SecondUserId   string     `json:"secondUserId,omitempty"`
	State          string     `json:"state,omitempty"`
	LastMoveUserId string     `json:"lastMoveUserId"`
	CreatedAt      *time.Time `json:"createdAt,omitempty"`
	LastModifiedAt *time.Time `json:"lastModifiedAt,omitempty"`
}
