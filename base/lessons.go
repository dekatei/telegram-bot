package base

import (
	"log"

	_ "github.com/mattn/go-sqlite3" // или postgres
)

// SQL запросы от учеников
// получаем доступные для записи уроки
func GetAvailableLessons() ([]Lesson, error) {
	rows, err := DB.Query("SELECT id, name, title, date FROM scheduler WHERE state = 0")
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

// добавляем ученика в таблицу с зарегистрированным уроком и устанавливаем статус недоступный урок в общей таблице
func RegisterUserToLesson(userID int, lessonID int) error {
	_, err := DB.Exec("INSERT OR IGNORE INTO registrations (user_id, lesson_id) VALUES (?, ?)", userID, lessonID)
	if err != nil {
		log.Print("Не удалось добавить запись")
	}
	_, err = DB.Exec("UPDATE scheduler SET state = 1 WHERE id = ?", lessonID)
	return err
}

// получаем зарегистрированные уроки ученика
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

func GetLessonsByDate(date string) ([]Lesson, error) {
	rows, err := DB.Query(`
		SELECT id, name, title, date 
		FROM scheduler 
		WHERE date LIKE ? AND state = 0
		ORDER BY date ASC
	`, date+"%")
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
