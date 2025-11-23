package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

// Task represents a task in the scheduler system.
type Task struct {
	ID      int64  `json:"id,omitempty"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
}

// AddTask adds a new task to the database
func AddTask(task *Task) (string, error) {
	var id int64

	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return "", err
	}

	id, err = res.LastInsertId()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%d", id), nil
}

// Tasks returns a list of tasks from the database
func Tasks(limit int) ([]*Task, error) {
	rows, err := DB.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?", limit)
	if err != nil {
		return []*Task{}, err
	}

	defer rows.Close()

	res := []*Task{}
	for rows.Next() {
		t := Task{}
		err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			return []*Task{}, err
		}

		res = append(res, &t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

// GetTask returns a task by its ID.
func GetTask(id int64) (*Task, error) {
	t := &Task{}
	row := DB.QueryRow(`SELECT id, date, title, comment, repeat 
	FROM scheduler WHERE id = ?`, id)
	err := row.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		return nil, fmt.Errorf("task not found")
	}

	return t, nil
}

// UpdateTask updates an existing task in the database
func UpdateTask(task *Task) error {

	query := `UPDATE scheduler SET 
		date = ?, title = ?, comment = ?, repeat = ? 
		WHERE id = ?`
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}
	return nil
}

// DeleteTask deletes a task by its ID
func DeleteTask(id int64) error {
	_, err := DB.Exec("DELETE FROM scheduler WHERE id = ?", id)
	if err != nil {
		return err
	}
	return nil
}

// UpdateDate only updates the Date field of a task with a specified ID
func UpdateDate(next string, id int64) error {
	query := `UPDATE scheduler SET 
		date = ? WHERE id = ?`
	res, err := DB.Exec(query, next, id)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)

	}

	return nil
}

//SEARCH

// Tasks returns a list of tasks from the database
func SearchTasks(search string, limit int) ([]*Task, error) {
	var rows *sql.Rows
	var err error

	date, err := time.Parse("02.01.2006", search)
	if err == nil{
		search = date.Format("20060102")
	}

	like := "%" + search + "%"

	if search == "" {
		rows, err = DB.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?", limit)
		if err != nil {
			return []*Task{}, err
		}
	} else {
		rows, err = DB.Query("SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE ? OR comment LIKE ? OR date LIKE ? ORDER BY date LIMIT ?", like, like, like, limit)
		if err != nil {
			return []*Task{}, err
		}
	}

	defer rows.Close()

	res := []*Task{}
	for rows.Next() {
		t := Task{}
		err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			return []*Task{}, err
		}

		res = append(res, &t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}
