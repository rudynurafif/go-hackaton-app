// Package services berisi business logic + akses database.
// Padanan dari file *.service.ts di NestJS; SQL mentah menggantikan
// pemanggilan Prisma Client.
package services

import (
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"hackaton-management-app/config"
	"hackaton-management-app/dto"
	"hackaton-management-app/models"
	"hackaton-management-app/utils"
)

// AuthService menggantikan Better Auth: registrasi dengan hash bcrypt,
// login yang menerbitkan JWT.
type AuthService struct {
	DB  *sql.DB
	Cfg *config.Config
}

// AuthResult — bentuk data response register/login: user + access token.
type AuthResult struct {
	User  models.User `json:"user"`
	Token string      `json:"token"`
}

// Register membuat user baru. Role tidak pernah diterima dari input —
// selalu jatuh ke default kolom (PARTICIPANT), sama seperti `input: false`
// di konfigurasi Better Auth.
func (s *AuthService) Register(input dto.RegisterDTO) (*AuthResult, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	var user models.User
	err = s.DB.QueryRow(
		`INSERT INTO users (name, email, password)
		 VALUES ($1, $2, $3)
		 RETURNING id, name, email, role, created_at, updated_at`,
		input.Name, input.Email, string(hash),
	).Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		// Constraint UNIQUE pada kolom email menolak pendaftaran ganda —
		// race-safe, tidak perlu check-then-insert.
		if isUniqueViolation(err) {
			return nil, utils.NewBadRequest("User with this email already exists")
		}
		return nil, err
	}

	return s.buildAuthResult(user)
}

// Login memverifikasi kredensial dan menerbitkan JWT. Pesan error sengaja
// sama untuk "email tidak terdaftar" dan "password salah" agar tidak
// membocorkan email mana yang terdaftar.
func (s *AuthService) Login(input dto.LoginDTO) (*AuthResult, error) {
	var user models.User
	err := s.DB.QueryRow(
		`SELECT id, name, email, password, role, created_at, updated_at
		 FROM users WHERE email = $1`,
		input.Email,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, utils.NewUnauthorized("Invalid email or password")
	}
	if err != nil {
		return nil, err
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)) != nil {
		return nil, utils.NewUnauthorized("Invalid email or password")
	}

	return s.buildAuthResult(user)
}

func (s *AuthService) buildAuthResult(user models.User) (*AuthResult, error) {
	token, err := utils.GenerateToken(user.ID, string(user.Role), s.Cfg.JWTSecret, s.Cfg.JWTExpires)
	if err != nil {
		return nil, err
	}
	return &AuthResult{User: user, Token: token}, nil
}
