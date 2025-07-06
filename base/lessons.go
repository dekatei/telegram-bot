package base

import (
	_ "github.com/mattn/go-sqlite3" // или postgres
)

func GetAllLessons() ([]Lesson, error) {
	rows, err := DB.Query("SELECT id, name, title, date FROM scheduler")
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

func RegisterUserToLesson(userID int, lessonID int) error {
	_, err := DB.Exec("INSERT OR IGNORE INTO registrations (user_id, lesson_id) VALUES (?, ?)", userID, lessonID)
	return err
}

func GetUserLessons(userID int) ([]Lesson, error) {
	query := `
	SELECT l.id, l.name, l.title, l.date
	FROM scheduler l
	INNER JOIN registrations r ON r.lesson_id = l.id
	WHERE r.user_id = ?
	`
	rows, err := DB.Query(query, userID)
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
