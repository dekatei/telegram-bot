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
}

/*
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
*/
func InitDB(envDBFILE string) error {
	var err error
	var appPath string
	if envDBFILE != "" {
		appPath = envDBFILE
	} else {
		appPath, err = os.Getwd() //не смогла реализовать через os.Executable()
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
        	date TEXT NOT NULL DEFAULT "");`)
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
