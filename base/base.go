package base

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3" // или postgres
)

var DB *sql.DB

type Lesson struct {
	ID    int
	Name  string
	Title string
	Date  string
	State string
}

// инициализируем таблицу с уроками
func InitDB(envDBFILE string) error {
	var err error
	var appPath string
	if envDBFILE != "" {
		appPath = envDBFILE
	} else {
		appPath, err = os.Getwd()
		if err != nil {
			log.Fatal(err)
			return err
		}
	}

	var install bool
	dbFile := filepath.Join(appPath, "scheduler.db")
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		install = true
		fmt.Println("db не найдена, создаём новую")
	}

	DB, err = sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatal(err)
		return err
	}

	if install {

		// Создание базы данных
		_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS scheduler(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
       		name TEXT NOT NULL DEFAULT "",
			title TEXT NOT NULL DEFAULT "",
        	date TEXT NOT NULL DEFAULT "",
			state BOOLEAN DEFAULT 0);`)
		if err != nil {
			log.Fatal(err)
			return err
		}
		_, err = DB.Exec(`CREATE INDEX IF NOT EXISTS scheduler_date ON scheduler (date);`)
		if err != nil {
			log.Fatal(err)
			return err
		}

		fmt.Println("База данных успешно создана!")
	}
	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS registrations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		lesson_id INTEGER,
		UNIQUE(user_id, lesson_id)
	);`)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
