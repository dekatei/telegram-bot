package base

import (
	_ "github.com/mattn/go-sqlite3" // или postgres
)

func AddLesson(name, title, date string) error {
	_, err := DB.Exec(`INSERT INTO scheduler (name, title, date) VALUES (?, ?, ?)`, name, title, date)
	return err
}
