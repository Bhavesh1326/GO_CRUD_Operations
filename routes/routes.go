package routes

import (
	"CRUD_go/handlers"
	"net/http"
)

func InitRoutes() {
	http.HandleFunc("/items", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handlers.GetItems(w, r)
		} else if r.Method == http.MethodPost {
			handlers.CreateItem(w, r)
		}
	})

	http.HandleFunc("/items/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.GetItem(w, r)
		case http.MethodPut:
			handlers.UpdateItem(w, r)
		case http.MethodDelete:
			handlers.DeleteItem(w, r)
		}
	})
}
