package base

import (
	"log"

	_ "github.com/mattn/go-sqlite3" // или postgres
)

// SQL запросы от админа
const textAdmin = "Получен запрос от админа "

func AddLesson(name, title, date string) error {
	log.Printf(textAdmin + "на добавление урока")
	_, err := DB.Exec(`INSERT INTO scheduler (name, title, date) VALUES (?, ?, ?)`, name, title, date)
	return err
}

func DeleteLesson(id int) error {
	log.Printf(textAdmin + "на удаление урока")
	_, err := DB.Exec(`DELETE FROM scheduler WHERE id = ?`, id)
	return err
}

func GetAdminLessons() ([]Lesson, error) {
	query := `
	SELECT s.id, s.name, s.title, s.date
	FROM scheduler s
	INNER JOIN registrations r ON s.id = r.lesson_id
	ORDER BY s.date ASC
	`
	rows, err := DB.Query(query)
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
