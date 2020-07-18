package models

import "time"

type Game struct {
	ID             int        `json:"id,omitempty"`
	GameId         string     `json:"game_id,omitempty"`
	FirstUserId    string     `json:"first_user_id,omitempty"`
	SecondUserId   string     `json:"second_user_id,omitempty"`
	State          string     `json:"move,omitempty"`
	LastMoveUserId string     `json:"last_move_user_id"`
	CreatedAt      *time.Time `json:"created_at,omitempty"`
	LastModifiedAt *time.Time `json:"last_modified_at"`
}
