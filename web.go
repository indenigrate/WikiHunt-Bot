package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/google/uuid"
)

// Task represents a pathfinding task.
type Task struct {
	ID    string   `json:"id"`
	Path  []string `json:"path"`
	Done  bool     `json:"done"`
	Error string   `json:"error,omitempty"`
}

var (
	tasks      = make(map[string]*Task)
	tasksMutex = &sync.Mutex{}
)

// PathRequest models the expected JSON request body for initiating a pathfinding task.
type PathRequest struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

// PathResponse models the JSON response after initiating a task.
type PathResponse struct {
	TaskID string `json:"task_id"`
}

// pathHandler creates and starts a new pathfinding task.
func pathHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var req PathRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	taskID := uuid.New().String()
	task := &Task{
		ID:   taskID,
		Path: []string{req.Start},
		Done: false,
	}

	tasksMutex.Lock()
	tasks[taskID] = task
	tasksMutex.Unlock()

	go wikiHuntAsync(task, req.Start, req.End, false)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(PathResponse{TaskID: taskID})
}

// statusHandler returns the status of a pathfinding task.
func statusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	taskID := strings.TrimPrefix(r.URL.Path, "/path/")
	tasksMutex.Lock()
	task, ok := tasks[taskID]
	tasksMutex.Unlock()

	if !ok {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// wikiHuntAsync performs the pathfinding asynchronously, updating the task object.
func wikiHuntAsync(task *Task, start string, end string, backlinks bool) {
	if backlinks {
		start, end = end, start
	}

	current := start
	traversed := make(map[string]bool)

	for len(task.Path) < 50 {
		if current == end {
			break // Success
		}

		traversed[current] = true
		log.Println("Current:", current)

		links, err := fetchWikiLinks(current, backlinks)
		if err != nil {
			task.Error = fmt.Sprintf("failed to fetch links for '%s': %v", current, err)
			break
		}

		topChoices, err := checkSimilarity(end, links, traversed)
		if err != nil {
			task.Error = fmt.Sprintf("failed to check similarity for '%s': %v", current, err)
			break
		}

		if len(topChoices) == 0 {
			task.Error = fmt.Sprintf("stuck at '%s', no further non-traversed links found", current)
			break
		}

		current = topChoices[0].Choice
		tasksMutex.Lock()
		task.Path = append(task.Path, current)
		tasksMutex.Unlock()
	}

	if task.Error == "" && current != end {
		task.Error = "path not found within 50 steps"
	}

	tasksMutex.Lock()
	task.Done = true
	tasksMutex.Unlock()
}

// StartServer initializes and starts the HTTP server.
func StartServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/path", pathHandler)
	mux.HandleFunc("/path/", statusHandler)

	// Add a simple CORS middleware
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		mux.ServeHTTP(w, r)
	})

	log.Println("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("Error starting server: %s\n", err)
	}
}
