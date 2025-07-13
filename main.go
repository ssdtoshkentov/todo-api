package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Todo model
type Todo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

// Foydalanuvchi korsatish uchun "yodda saqlanadigon" todo royxati
var todos = []Todo{
	{ID: 1, Title: "Kitob o‘qish", Done: false},
	{ID: 2, Title: "Go dasturlash tilini o‘rganish", Done: true},
}

// Handler: Get /Todos
func getTodosHandler(w http.ResponseWriter, r *http.Request) {
	//Ruxsat beramiz json farmati
	w.Header().Set("Content-Type", "application/json")

	// Todos ni Json ga, o'girib Foydalanuvchi ga yuboramiz
	err := json.NewEncoder(w).Encode(todos)
	if err != nil {
		http.Error(w, "Xatolik: JSON chiqish bulmadi", http.StatusInternalServerError)
	}
}

// Post /todos
func createTodoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Faqat Post method ga ruxsat", http.StatusMethodNotAllowed)
		return
	}

	var newTodo Todo
	err := json.NewDecoder(r.Body).Decode(&newTodo)
	if err != nil {
		http.Error(w, "Xato no'tog'ri JSON", http.StatusBadRequest)
		return
	}

	newTodo.ID = len(todos) + 1
	todos = append(todos, newTodo)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newTodo)
}

// main function
func main() {
	http.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			getTodosHandler(w, r)
		} else if r.Method == http.MethodPost {
			createTodoHandler(w, r)
		} else {
			http.Error(w, "No'tog'ri Method", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("Server ishga tushdi: http://localhost:8080/todos")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Server da xatolik:", err)
	}
}
