package dto

import "time"

// CreateHackathonDTO — padanan CreateHackathonDto (class-validator).
// "future" adalah custom validator (didaftarkan di routes) yang memeriksa
// tanggal harus setelah waktu request — padanan @MinDate(() => new Date()).
// Field tanggal memakai pointer supaya `required` bisa membedakan
// "tidak dikirim" dari zero value.
type CreateHackathonDTO struct {
	Name        string     `json:"name" binding:"required,min=3"`
	Description *string    `json:"description" binding:"omitempty,min=10,max=1000"`
	StartsAt    *time.Time `json:"startsAt" binding:"required,future"`
	EndsAt      *time.Time `json:"endsAt" binding:"required,future"`
	IsActive    *bool      `json:"isActive"`
}

// UpdateHackathonDTO — padanan PartialType(CreateHackathonDto):
// semua field opsional (pointer + omitempty), tetapi aturan validasi yang
// sama tetap berlaku untuk field yang dikirim (semantik PATCH).
type UpdateHackathonDTO struct {
	Name        *string    `json:"name" binding:"omitempty,min=3"`
	Description *string    `json:"description" binding:"omitempty,min=10,max=1000"`
	StartsAt    *time.Time `json:"startsAt" binding:"omitempty,future"`
	EndsAt      *time.Time `json:"endsAt" binding:"omitempty,future"`
	IsActive    *bool      `json:"isActive"`
}
