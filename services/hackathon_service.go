package services

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"

	"hackaton-management-app/dto"
	"hackaton-management-app/models"
	"hackaton-management-app/utils"
)

const hackathonColumns = `id, name, description, start_date, end_date, is_active, author_id, created_at, updated_at`

// HackathonService — padanan hackaton.service.ts.
type HackathonService struct {
	DB *sql.DB
}

func scanHackathon(row interface{ Scan(...any) error }) (*models.Hackathon, error) {
	var h models.Hackathon
	err := row.Scan(
		&h.ID, &h.Name, &h.Description, &h.StartDate, &h.EndDate,
		&h.IsActive, &h.AuthorID, &h.CreatedAt, &h.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &h, nil
}

// Create membuat hackathon milik author yang diberikan (admin yang login).
func (s *HackathonService) Create(input dto.CreateHackathonDTO, authorID string) (*models.Hackathon, error) {
	// isActive opsional; jika tidak dikirim, jatuh ke default kolom (false).
	isActive := false
	if input.IsActive != nil {
		isActive = *input.IsActive
	}

	return scanHackathon(s.DB.QueryRow(
		`INSERT INTO hackathons (name, description, start_date, end_date, is_active, author_id)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING `+hackathonColumns,
		input.Name, input.Description, input.StartsAt, input.EndsAt, isActive, authorID,
	))
}

func (s *HackathonService) FindAll() ([]models.Hackathon, error) {
	rows, err := s.DB.Query(`SELECT ` + hackathonColumns + ` FROM hackathons ORDER BY created_at`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	hackathons := []models.Hackathon{}
	for rows.Next() {
		h, err := scanHackathon(rows)
		if err != nil {
			return nil, err
		}
		hackathons = append(hackathons, *h)
	}
	return hackathons, rows.Err()
}

func (s *HackathonService) FindOne(id string) (*models.Hackathon, error) {
	h, err := scanHackathon(s.DB.QueryRow(
		`SELECT `+hackathonColumns+` FROM hackathons WHERE id = $1`, id,
	))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, utils.NewNotFound(fmt.Sprintf("Hackathon with id %q not found", id))
	}
	if err != nil {
		return nil, err
	}
	return h, nil
}

// Update melakukan partial update: hanya field non-nil yang masuk ke klausa
// SET — padanan "undefined fields are ignored by Prisma" di service NestJS.
func (s *HackathonService) Update(id string, input dto.UpdateHackathonDTO) (*models.Hackathon, error) {
	if _, err := s.FindOne(id); err != nil {
		return nil, err // 404 jika tidak ada
	}

	sets := []string{}
	args := []any{}
	addSet := func(column string, value any) {
		args = append(args, value)
		sets = append(sets, fmt.Sprintf("%s = $%d", column, len(args)))
	}

	if input.Name != nil {
		addSet("name", *input.Name)
	}
	if input.Description != nil {
		addSet("description", *input.Description)
	}
	if input.StartsAt != nil {
		addSet("start_date", *input.StartsAt)
	}
	if input.EndsAt != nil {
		addSet("end_date", *input.EndsAt)
	}
	if input.IsActive != nil {
		addSet("is_active", *input.IsActive)
	}

	if len(sets) == 0 {
		// Tidak ada field yang dikirim — tidak ada yang diubah.
		return s.FindOne(id)
	}

	sets = append(sets, "updated_at = now()")
	args = append(args, id)

	query := fmt.Sprintf(
		`UPDATE hackathons SET %s WHERE id = $%d RETURNING %s`,
		strings.Join(sets, ", "), len(args), hackathonColumns,
	)
	return scanHackathon(s.DB.QueryRow(query, args...))
}

func (s *HackathonService) Remove(id string) (*models.Hackathon, error) {
	if _, err := s.FindOne(id); err != nil {
		return nil, err // 404 jika tidak ada
	}

	// Peserta ikut terhapus lewat ON DELETE CASCADE.
	return scanHackathon(s.DB.QueryRow(
		`DELETE FROM hackathons WHERE id = $1 RETURNING `+hackathonColumns, id,
	))
}

// Join mendaftarkan user sebagai peserta hackathon yang sedang dibuka.
func (s *HackathonService) Join(id, userID string) (*models.HackathonParticipant, error) {
	hackathon, err := s.FindOne(id) // 404 jika tidak ada
	if err != nil {
		return nil, err
	}

	if !hackathon.IsActive {
		return nil, utils.NewBadRequest("This hackathon is not active")
	}

	if !hackathon.EndDate.After(time.Now()) {
		return nil, utils.NewBadRequest("This hackathon has already ended")
	}

	var p models.HackathonParticipant
	err = s.DB.QueryRow(
		`INSERT INTO hackathon_participants (hackathon_id, user_id)
		 VALUES ($1, $2)
		 RETURNING id, hackathon_id, user_id, joined_at`,
		id, userID,
	).Scan(&p.ID, &p.HackathonID, &p.UserID, &p.JoinedAt)
	if err != nil {
		// Andalkan constraint UNIQUE (hackathon_id, user_id) untuk menolak
		// join ganda — race-safe. Padanan penanganan error P2002 Prisma.
		if isUniqueViolation(err) {
			return nil, utils.NewBadRequest("You have already joined this hackathon")
		}
		return nil, err
	}
	return &p, nil
}

// isUniqueViolation memeriksa apakah error berasal dari pelanggaran
// constraint UNIQUE PostgreSQL (SQLSTATE 23505) — padanan kode error
// P2002 pada Prisma.
func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}
