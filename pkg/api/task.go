package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"main/pkg/db"
	"net/http"
	"strconv"
)

type Task struct {
	ID      string `json:"id,omitempty"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
}

// taskHandler handler for working with a single task
func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// TasksResp - a structure from the task list
type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

// tasksHandler — a handler for receiving a list of tasks -----------------------------------
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	search := r.URL.Query().Get("search")

	tasks, err := db.SearchTasks(search, 50)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusFailedDependency)
		return
	}
	
	apiTasks := make([]Task, len(tasks))
	for i, t := range tasks {
		apiTasks[i] = Task{
			ID:      fmt.Sprintf("%d", t.ID),
			Date:    t.Date,
			Title:   t.Title,
			Comment: t.Comment,
			Repeat:  t.Repeat,
		}
	}
	writeJson(w, map[string][]Task{"tasks": apiTasks}, http.StatusOK)
}
//___________________________________________________________________________________________

// getTaskHandler — a handler for getting a single task
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, map[string]string{"error": "id is required"}, http.StatusBadRequest)
		return
	}

	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		writeJson(w, map[string]string{"error": "invalid id"}, http.StatusBadRequest)
		return
	}

	t, err := db.GetTask(idInt)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusFailedDependency)
		return
	}

	apiTask := Task{
		ID:      fmt.Sprintf("%d", t.ID),
		Date:    t.Date,
		Title:   t.Title,
		Comment: t.Comment,
		Repeat:  t.Repeat,
	}

	writeJson(w, apiTask, http.StatusOK)
}

// updateTaskHandler — a handler for updating an existing task
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		writeJson(w, map[string]string{"error": "title is required"}, http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(task.ID, 10, 64)
	if err != nil {
		writeJson(w, map[string]string{"error": "invalid id"}, http.StatusBadRequest)
		return
	}

	dbTask := &db.Task{
		ID:      id,
		Date:    task.Date,
		Title:   task.Title,
		Comment: task.Comment,
		Repeat:  task.Repeat,
	}

	if err = checkDate(dbTask); err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	if err = db.UpdateTask(dbTask); err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusFailedDependency)
		return
	}

	writeJson(w, map[string]string{"status": "ok"}, http.StatusOK)
}

// deleteTaskHandler — handler for deleting a task by ID
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, map[string]string{"error": "id is required"}, http.StatusBadRequest)
		return
	}

	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		writeJson(w, map[string]string{"error": "invalid id"}, http.StatusBadRequest)
		return
	}

	if err = db.DeleteTask(idInt); err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusFailedDependency)
		return
	}

	writeJson(w, map[string]any{}, http.StatusOK)
}
