package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

var schema = `CREATE TABLE scheduler (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  date CHAR(8) NOT NULL DEFAULT "",
  title VARCHAR(100) NOT NULL DEFAULT "",
  comment VARCHAR(256) NOT NULL DEFAULT "",
  repeat VARCHAR(128) NOT NULL DEFAULT ""
);

CREATE INDEX date_task ON scheduler (date);`

var db *sql.DB

//Init initializes the connection to the SQLite database
func Init(dbFile string) error {
	_, err := os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}

	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	if install {
		_, err = db.Exec(schema)
		if err != nil {
			return err
		}
	}
	return nil
}
