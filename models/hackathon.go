package models

import "time"

// Hackathon — padanan model `Hackathon` di schema.prisma.
// Description bertipe pointer karena kolomnya nullable; nilai nil akan
// ter-serialize sebagai null di JSON, sama seperti output Prisma.
type Hackathon struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	StartDate   time.Time `json:"startDate"`
	EndDate     time.Time `json:"endDate"`
	IsActive    bool      `json:"isActive"`
	AuthorID    string    `json:"authorId"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
