// Entry point aplikasi — padanan dari main.ts (bootstrap) di NestJS.
package main

import (
	"log"

	"hackaton-management-app/config"
	"hackaton-management-app/database"
	"hackaton-management-app/routes"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := database.Migrate(db); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	router := routes.Setup(db, cfg)

	log.Printf("server listening on http://localhost:%s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
