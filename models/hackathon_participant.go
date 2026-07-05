package models

import "time"

// HackathonParticipant — padanan model `HackathonParticipant`.
// Constraint UNIQUE (hackathon_id, user_id) di database mencegah
// satu user join hackathon yang sama dua kali.
type HackathonParticipant struct {
	ID          string    `json:"id"`
	HackathonID string    `json:"hackathonId"`
	UserID      string    `json:"userId"`
	JoinedAt    time.Time `json:"joinedAt"`
}
