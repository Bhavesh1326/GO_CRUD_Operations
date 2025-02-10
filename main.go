package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"
)

type Item struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Email string `json:"email"`
}

const connStr = "user=postgres password=password dbname=items_db sslmode=disable"

var db *sql.DB

func initDB() {
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error opening database:", err)
	}
	fmt.Println("Connected to PostgreSQL")
}

func createItem(w http.ResponseWriter, r *http.Request) {
	var item Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sqlStatement := `INSERT INTO items (name, phone, email) 
                     VALUES ($1, $2, $3) RETURNING id`
	err := db.QueryRow(sqlStatement, item.Name, item.Phone, item.Email).Scan(&item.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

func getItems(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name, phone, email FROM items")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var item Item
		if err := rows.Scan(&item.ID, &item.Name, &item.Phone, &item.Email); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func getItem(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/items/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var item Item
	sqlStatement := `SELECT id, name, phone, email FROM items WHERE id=$1`
	err = db.QueryRow(sqlStatement, id).Scan(&item.ID, &item.Name, &item.Phone, &item.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Item not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

func updateItem(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/items/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var updatedItem Item
	if err := json.NewDecoder(r.Body).Decode(&updatedItem); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sqlStatement := `UPDATE items SET name=$1, phone=$2, email=$3 WHERE id=$4`
	_, err = db.Exec(sqlStatement, updatedItem.Name, updatedItem.Phone, updatedItem.Email, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	updatedItem.ID = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedItem)
}

func deleteItem(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/items/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	sqlStatement := `DELETE FROM items WHERE id=$1`
	_, err = db.Exec(sqlStatement, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func main() {
	initDB()
	defer db.Close()

	http.HandleFunc("/items", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			getItems(w, r)
		} else if r.Method == http.MethodPost {
			createItem(w, r)
		}
	})

	http.HandleFunc("/items/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getItem(w, r)
		case http.MethodPut:
			updateItem(w, r)
		case http.MethodDelete:
			deleteItem(w, r)
		}
	})

	fmt.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
