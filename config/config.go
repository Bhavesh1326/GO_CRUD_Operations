package config

import "CRUD_go/models"

func NewConfig() *models.Config {
	return &models.Config{
		DBUser:     "postgres",
		DBPassword: "password",
		DBName:     "items_db",
		DBSSLMode:  "disable",
	}
}
