package main

import (
	"CRUD_go/config"
	"CRUD_go/db"
	"CRUD_go/routes"
	"fmt"
	"log"
	"net/http"
)

func main() {

	cfg := config.NewConfig()

	db.InitDB(cfg)
	defer db.CloseDB()

	routes.InitRoutes()

	fmt.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
