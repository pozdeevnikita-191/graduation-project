package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"main/pkg/db"
	"net/http"
	"time"
)

//the writeJson function encodes data into json format
func writeJson(w http.ResponseWriter, data any, code int) {
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    if err := json.NewEncoder(w).Encode(data); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
		w.WriteHeader(code)
}

// The addTaskHandler handler adds a new task
func addTaskHandler(w http.ResponseWriter, r *http.Request) {
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
	if task.Date != "" {
		if _, err := time.Parse(fdate, task.Date); err != nil {
			writeJson(w, map[string]string{"error": "invalid date format, expected YYYYMMDD"}, http.StatusBadRequest)
			return
		}
	}

	dbTask := &db.Task{
		Date:    task.Date,
		Title:   task.Title,
		Comment: task.Comment,
		Repeat:  task.Repeat,
	}

	err = checkDate(dbTask)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	id, err := db.AddTask(dbTask)
	if err != nil {
		writeJson(w, map[string]string{"error": "failed to add task"}, http.StatusBadRequest)
		return
	}

	writeJson(w, map[string]string{"id": id}, http.StatusCreated)
}

// the checkDate function checks the date for relevance
func checkDate(task *db.Task) error {
	now := time.Now().Truncate(24 * time.Hour)

	if task.Date == "" {
		task.Date = now.Format(fdate)
	}

	t, err := time.Parse(fdate, task.Date)
	if err != nil {
		return fmt.Errorf("invalid date format, expected YYYYMMDD: %w", err)
	}

	var next string
	if task.Repeat != "" {
		next, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
	}

	if t.Before(now) {
		if task.Repeat == "" {
			task.Date = now.Format(fdate)
		} else {
			task.Date = next
		}
	}
	return nil
}


