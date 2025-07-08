package base

import (
	"log"

	_ "github.com/mattn/go-sqlite3" // или postgres
)

const textAdmin = "Получен запрос от админа "

func AddLesson(name, title, date string) error {
	log.Printf(textAdmin + "на добавление урока")
	_, err := DB.Exec(`INSERT INTO scheduler (name, title, date, state) VALUES (?, ?, ?, ?)`, name, title, date, false)
	return err
}

func DeleteLesson(id int) error {
	log.Printf(textAdmin + "на удаление урока")
	_, err := DB.Exec(`DELETE FROM scheduler WHERE id = ?`, id)
	return err
}

func GetAdminLessons() ([]Lesson, error) {
	rows, err := DB.Query("SELECT id, name, title, date FROM scheduler WHERE state = 1")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lessons []Lesson
	for rows.Next() {
		var l Lesson
		if err := rows.Scan(&l.ID, &l.Name, &l.Title, &l.Date); err != nil {
			return nil, err
		}
		lessons = append(lessons, l)
	}

	return lessons, nil
}
