package models

type Config struct {
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
}

type Item struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Email string `json:"email"`
}
