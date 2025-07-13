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

type Todo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

// Yodda saqlanadigan todo ro'yxati (fayldan yuklanadi)
var todos []Todo

// Fayldan o‘qish funksiyasi
func loadTodosFromFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &todos)
}

// Faylga yozish funksiyasi
func saveTodosToFile(filename string) error {
	data, err := json.MarshalIndent(todos, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, data, 0644)
}

// GET /todos - Barcha todo’larni olish
func getTodosHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(todos)
	if err != nil {
		http.Error(w, "Xatolik: JSON chiqish bo‘lmadi", http.StatusInternalServerError)
	}
}

// POST /todos - Yangi todo yaratish
func createTodoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Faqat POST methodga ruxsat", http.StatusMethodNotAllowed)
		return
	}

	var newTodo Todo
	err := json.NewDecoder(r.Body).Decode(&newTodo)
	if err != nil {
		http.Error(w, "Xato: noto‘g‘ri JSON", http.StatusBadRequest)
		return
	}

	// ID avtomatik
	newTodo.ID = getNextID()
	todos = append(todos, newTodo)
	saveTodosToFile("todos.json")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newTodo)
}

// PATCH /todos/{id} - Mavjud todo’ni yangilash
func updateTodoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Faqat PATCH methodga ruxsat", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/todos/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Noto‘g‘ri ID", http.StatusBadRequest)
		return
	}

	var updatedTodo Todo
	err = json.NewDecoder(r.Body).Decode(&updatedTodo)
	if err != nil {
		http.Error(w, "Xato: noto‘g‘ri JSON", http.StatusBadRequest)
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

	http.Error(w, "Topilmadi: todo mavjud emas", http.StatusNotFound)
}

// DELETE /todos/{id} - Todo’ni o‘chirish
func deleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Faqat DELETE methodga ruxsat", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/todos/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Noto‘g‘ri ID", http.StatusBadRequest)
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

	http.Error(w, "Topilmadi: todo mavjud emas", http.StatusNotFound)
}

// IDni topuvchi yordamchi funksiya
func getNextID() int {
	maxID := 0
	for _, t := range todos {
		if t.ID > maxID {
			maxID = t.ID
		}
	}
	return maxID + 1
}

// main function
func main() {
	err := loadTodosFromFile("todos.json")
	if err != nil {
		fmt.Println("Xatolik: todos.json fayl o‘qilmadi:", err)
		os.Exit(1)
	}

	// GET va POST uchun
	http.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			getTodosHandler(w, r)
		} else if r.Method == http.MethodPost {
			createTodoHandler(w, r)
		} else {
			http.Error(w, "Noto‘g‘ri Method", http.StatusMethodNotAllowed)
		}
	})

	// PATCH va DELETE uchun
	http.HandleFunc("/todos/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPatch {
			updateTodoHandler(w, r)
		} else if r.Method == http.MethodDelete {
			deleteTodoHandler(w, r)
		} else {
			http.Error(w, "Noto‘g‘ri method", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("Server ishga tushdi: http://localhost:8080/todos")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Xatolik:", err)
	}
}
