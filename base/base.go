package base

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3" // или postgres
)

var DB *sql.DB

type Lesson struct {
	ID   int
	Name string
	Date string
}

func InitDB(path string) {
	var err error
	DB, err = sql.Open("sqlite3", path)
	if err != nil {
		log.Fatal("DB open error:", err)
	}

	createTable := `
    CREATE TABLE IF NOT EXISTS lessons (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT,
        date TEXT
    );`
	_, err = DB.Exec(createTable)
	if err != nil {
		log.Fatal("Create table error:", err)
	}
}

func GetAllLessons() ([]Lesson, error) {
	rows, err := DB.Query("SELECT id, name, date FROM lessons")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lessons []Lesson
	for rows.Next() {
		var l Lesson
		if err := rows.Scan(&l.ID, &l.Name, &l.Date); err != nil {
			return nil, err
		}
		lessons = append(lessons, l)
	}

	return lessons, nil
}
