package db

import (
	"database/sql"
	"fmt"
	"log"

	// "CRUD_go/config"
	"CRUD_go/models"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB(cfg *models.Config) {
	var err error
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode)

	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error opening database:", err)
	}
	log.Println("Connected to PostgreSQL")
}

func CloseDB() {
	if err := DB.Close(); err != nil {
		log.Fatal("Error closing database:", err)
	}
}
