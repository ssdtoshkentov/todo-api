package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// Todo represents a task item in the list.
type Todo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

// todos stores the list of all todo items (loaded from file)
var todos []Todo

// loadTodosFromFile reads the todos from the specified JSON file.
func loadTodosFromFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &todos)
}

// saveTodosToFile writes the current todos to the specified JSON file.
func saveTodosToFile(filename string) error {
	data, err := json.MarshalIndent(todos, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, data, 0644)
}

// getTodosHandler handles GET /todos requests and returns all todo items.
func getTodosHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(todos)
	if err != nil {
		http.Error(w, "Error: Failed to encode todos to JSON", http.StatusInternalServerError)
	}
}

// createTodoHandler handles POST /todos requests to create a new todo item.
func createTodoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed: Only POST is supported", http.StatusMethodNotAllowed)
		return
	}

	var newTodo Todo
	err := json.NewDecoder(r.Body).Decode(&newTodo)
	if err != nil {
		http.Error(w, "Error: Invalid JSON body", http.StatusBadRequest)
		return
	}

	// Assign unique ID automatically
	newTodo.ID = getNextID()
	todos = append(todos, newTodo)
	saveTodosToFile("todos.json")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newTodo)
}

// updateTodoHandler handles PATCH /todos/{id} requests to update a todo item.
func updateTodoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed: Only PATCH is supported", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/todos/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID in URL", http.StatusBadRequest)
		return
	}

	var updatedTodo Todo
	err = json.NewDecoder(r.Body).Decode(&updatedTodo)
	if err != nil {
		http.Error(w, "Error: Invalid JSON body", http.StatusBadRequest)
		return
	}

	for i, t := range todos {
		if t.ID == id {
			todos[i].Title = updatedTodo.Title
			todos[i].Done = updatedTodo.Done
			saveTodosToFile("todos.json")

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(todos[i])
			return
		}
	}

	http.Error(w, "Not Found: Todo item does not exist", http.StatusNotFound)
}

// deleteTodoHandler handles DELETE /todos/{id} requests to delete a todo item.
func deleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed: Only DELETE is supported", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/todos/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID in URL", http.StatusBadRequest)
		return
	}

	for i, t := range todos {
		if t.ID == id {
			todos = append(todos[:i], todos[i+1:]...)
			saveTodosToFile("todos.json")
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.Error(w, "Not Found: Todo item does not exist", http.StatusNotFound)
}

// getNextID returns a new unique ID for a new todo item.
func getNextID() int {
	maxID := 0
	for _, t := range todos {
		if t.ID > maxID {
			maxID = t.ID
		}
	}
	return maxID + 1
}

// main is the entry point of the application.
func main() {
	err := loadTodosFromFile("todos.json")
	if err != nil {
		fmt.Println("Error: Failed to read todos.json file:", err)
		os.Exit(1)
	}

	// Handle GET and POST on /todos
	http.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			getTodosHandler(w, r)
		} else if r.Method == http.MethodPost {
			createTodoHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Handle PATCH and DELETE on /todos/{id}
	http.HandleFunc("/todos/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPatch {
			updateTodoHandler(w, r)
		} else if r.Method == http.MethodDelete {
			deleteTodoHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("Server is running at http://localhost:8080/todos")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
