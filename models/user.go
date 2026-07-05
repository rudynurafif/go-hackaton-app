package models

import "time"

// Role — padanan enum Role di schema.prisma.
type Role string

const (
	RoleParticipant Role = "PARTICIPANT"
	RoleAdmin       Role = "ADMIN"
)

// User — padanan model `user`. Field Password memakai tag `json:"-"`
// sehingga hash bcrypt tidak pernah ikut ter-serialize ke response.
type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Role      Role      `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
