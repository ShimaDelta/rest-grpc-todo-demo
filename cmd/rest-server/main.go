package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type Task struct {
	ID    int32  `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

type TaskCreate struct {
	Title string `json:"title"`
}

type TaskUpdate struct {
	Done bool `json:"done"`
}

var (
	mu     sync.Mutex
	tasks  []Task
	nextID int32 = 1
)

func main() {
	http.HandleFunc("/tasks", handleTasks)
	http.HandleFunc("/tasks/", handleTaskByID)

	log.Println("REST server listening on http://localhost:8000 ...")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}

// GET /tasks, POST /tasks
func handleTasks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		listTasks(w, r)
	case http.MethodPost:
		createTask(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func listTasks(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	writeJSON(w, tasks)
}

func createTask(w http.ResponseWriter, r *http.Request) {
	var body TaskCreate
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(body.Title) == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	t := Task{
		ID:    nextID,
		Title: body.Title,
		Done:  false,
	}
	nextID++
	tasks = append(tasks, t)

	w.WriteHeader(http.StatusCreated)
	writeJSON(w, t)
}

// PATCH /tasks/{id}, DELETE /tasks/{id}
func handleTaskByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/tasks/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodPatch:
		updateTask(w, r, int32(id))
	case http.MethodDelete:
		deleteTask(w, r, int32(id))
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func updateTask(w http.ResponseWriter, r *http.Request, id int32) {
	var body TaskUpdate
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Done = body.Done
			writeJSON(w, tasks[i])
			return
		}
	}
	http.Error(w, "not found", http.StatusNotFound)
}

func deleteTask(w http.ResponseWriter, r *http.Request, id int32) {
	mu.Lock()
	defer mu.Unlock()

	newTasks := make([]Task, 0, len(tasks))
	for _, t := range tasks {
		if t.ID != id {
			newTasks = append(newTasks, t)
		}
	}
	tasks = newTasks
	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Println("failed to write json:", err)
	}
}
