# 📝 Todo API

A simple yet functional RESTful API for managing a Todo List, built in Go (Golang).  
This project demonstrates basic CRUD operations, JSON handling, and HTTP routing — ideal for learning and portfolio use.

## 🔧 Features

- Create new todo items
- Read and list all todos
- Update a specific todo
- Delete a todo
- JSON-based API
- Clean and simple structure

## 📁 Project Structure

todo-api/
├── main.go # Entry point
├── models/ # Data structures
├── handlers/ # HTTP handler functions
├── routes/ # Routing logic
└── utils/ # Helper functions (if any)

bash
Копировать
Редактировать

## 🚀 How to Run

### 1. Clone the repository

```bash
git clone https://github.com/ssdtoshkentov/todo-api.git
cd todo-api
2. Run the server
bash
Копировать
Редактировать
go run main.go
Server will start at http://localhost:8080

📬 API Endpoints
Method	Endpoint	Description
GET	/todos	Get all todos
POST	/todos	Create a new todo
PATCH	/todos/{id}	Update a todo by ID
DELETE	/todos/{id}	Delete a todo by ID

All requests and responses use application/json.

📌 Example Request: POST /todos
json
Копировать
Редактировать
{
  "title": "Learn Golang",
  "completed": false
}
📌 Example Response
json
Копировать
Редактировать
{
  "id": 1,
  "title": "Learn Golang",
  "completed": false
}
📖 Technologies Used
Go (Golang)

net/http

JSON encoding/decoding

RESTful API design

✅ Future Improvements
Connect to PostgreSQL or SQLite database

Add JWT Authentication

Use Clean Architecture

Write unit tests

Add Docker support

👨‍💻 Author
Created by Abduvali Toshkentov
Open to feedback and contributions!

📄 License
This project is licensed under the MIT License.

yaml
Копировать
Редактировать
