package api

import (
	"main/pkg/db"
	"net/http"
	"strconv"
	"time"
)

// the taskDoneHandler handler makes the task complete
func taskDoneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJson(w, map[string]string{"error": "method not allowed"})
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, map[string]string{"error": "id is required"})
		return
	}

	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		writeJson(w, map[string]string{"error": "invalid id"})
		return
	}

	task, err := db.GetTask(idInt)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	if task.Repeat == "" {
		if err = db.DeleteTask(idInt); err != nil {
			writeJson(w, map[string]string{"error": err.Error()})
			return
		}
		w.WriteHeader(http.StatusOK)
		writeJson(w, map[string]string{})
		return
	}

	now := time.Now().Truncate(24 * time.Hour)
	next, err := NextDate(now, task.Date, task.Repeat)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	if err = db.UpdateDate(next, idInt); err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	writeJson(w, map[string]string{})
}
