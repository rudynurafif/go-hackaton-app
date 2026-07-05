package services

import (
	"database/sql"
	"errors"
	"fmt"

	"hackaton-management-app/models"
	"hackaton-management-app/utils"
)

// UserService — padanan user.service.ts.
type UserService struct {
	DB *sql.DB
}

// FindAll mengembalikan semua user (tanpa kolom password).
func (s *UserService) FindAll() ([]models.User, error) {
	rows, err := s.DB.Query(
		`SELECT id, name, email, role, created_at, updated_at
		 FROM users ORDER BY created_at`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Diinisialisasi sebagai slice kosong (bukan nil) agar ter-serialize
	// sebagai [] di JSON, bukan null — sama seperti findMany Prisma.
	users := []models.User{}
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Role, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

// FindOne mengembalikan satu user berdasarkan id, atau 404 jika tidak ada.
func (s *UserService) FindOne(id string) (*models.User, error) {
	var u models.User
	err := s.DB.QueryRow(
		`SELECT id, name, email, role, created_at, updated_at
		 FROM users WHERE id = $1`,
		id,
	).Scan(&u.ID, &u.Name, &u.Email, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, utils.NewNotFound(fmt.Sprintf("User with id %q not found", id))
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}
