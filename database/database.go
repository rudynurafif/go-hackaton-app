// Package database mengelola koneksi PostgreSQL dan migrasi schema.
// Padanan dari PrismaService + `prisma migrate` di proyek NestJS.
package database

import (
	"database/sql"
	"embed"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib" // register driver "pgx" untuk database/sql
)

//go:embed migrations/*.sql
var migrations embed.FS

// Connect membuka pool koneksi ke PostgreSQL dan memastikan server terjangkau.
func Connect(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return db, nil
}

// Migrate menjalankan file migrasi SQL. Semua statement bersifat idempotent
// (CREATE TABLE IF NOT EXISTS, dst.) sehingga aman dijalankan setiap startup.
func Migrate(db *sql.DB) error {
	entries, err := migrations.ReadDir("migrations")
	if err != nil {
		return fmt.Errorf("read migrations dir: %w", err)
	}

	// fs.ReadDir sudah mengurutkan berdasarkan nama file (001_, 002_, ...).
	for _, entry := range entries {
		content, err := migrations.ReadFile("migrations/" + entry.Name())
		if err != nil {
			return fmt.Errorf("read migration %s: %w", entry.Name(), err)
		}
		if _, err := db.Exec(string(content)); err != nil {
			return fmt.Errorf("run migration %s: %w", entry.Name(), err)
		}
	}

	return nil
}
