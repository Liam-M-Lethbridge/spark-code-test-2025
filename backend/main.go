package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// Todo struct to match the OpenAPI schema
type Todo struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// the todos (in memory)
var todos []Todo

func main() {

	http.HandleFunc("/", ToDoListHandler)
	log.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func ToDoListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	// Handle preflight OPTIONS request
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// since we're using a single handler function for the get and post, we have to check if the method is a get or post request
	switch r.Method {

	// if it's a get request, return the todos
	case http.MethodGet:
		getToDos(w)

	case http.MethodPost:
		createToDo(w, r)
	// default case to disallow any other form of http request
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// naming function based on OpenAPI spec
func createToDo(w http.ResponseWriter, r *http.Request) {
	var item Todo

	// Parse JSON body into item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// checks:
	// make sure neither required fields are empty and return status code 400 if invalid inputs
	if item.Title == "" || item.Description == "" {
		http.Error(w, "Missing title or description", http.StatusBadRequest)
		return
	}
	todos = append(todos, item)

	json.NewEncoder(w).Encode(item)
}

func getToDos(w http.ResponseWriter) {
	json.NewEncoder(w).Encode(todos)
}
