package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

var schema = `CREATE TABLE IF NOT EXISTS scheduler (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  date CHAR(8) NOT NULL DEFAULT "",
  title VARCHAR(100) NOT NULL DEFAULT "",
  comment VARCHAR(256) NOT NULL DEFAULT "",
  repeat VARCHAR(128) NOT NULL DEFAULT ""
);

CREATE INDEX date_task ON scheduler (date);`

var DB *sql.DB

//Init initializes the connection to the SQLite database
func Init(dbFile string) error {
	_, err := os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}

	DB, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	if install {
		_, err = DB.Exec(schema)
		if err != nil {
			return err
		}
	}
	return nil
}
