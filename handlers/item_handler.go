package handlers

import (
	"CRUD_go/db"
	"CRUD_go/models"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
)

func CreateItem(w http.ResponseWriter, r *http.Request) {
	var item models.Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sqlStatement := `INSERT INTO items (name, phone, email) 
                     VALUES ($1, $2, $3) RETURNING id`
	err := db.DB.QueryRow(sqlStatement, item.Name, item.Phone, item.Email).Scan(&item.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

func GetItems(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT id, name, phone, email FROM items")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []models.Item
	for rows.Next() {
		var item models.Item
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

// GetItem handles fetching a specific item by ID.
func GetItem(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/items/"):]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var item models.Item
	sqlStatement := `SELECT id, name, phone, email FROM items WHERE id=$1`
	err = db.DB.QueryRow(sqlStatement, id).Scan(&item.ID, &item.Name, &item.Phone, &item.Email)
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

// UpdateItem handles updating an existing item.
func UpdateItem(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/items/"):]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var updatedItem models.Item
	if err := json.NewDecoder(r.Body).Decode(&updatedItem); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sqlStatement := `UPDATE items SET name=$1, phone=$2, email=$3 WHERE id=$4`
	_, err = db.DB.Exec(sqlStatement, updatedItem.Name, updatedItem.Phone, updatedItem.Email, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	updatedItem.ID = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedItem)
}

// DeleteItem handles deleting an item by ID.
func DeleteItem(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/items/"):]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	sqlStatement := `DELETE FROM items WHERE id=$1`
	_, err = db.DB.Exec(sqlStatement, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
