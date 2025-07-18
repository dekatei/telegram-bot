package base

import (
	"log"

	_ "github.com/mattn/go-sqlite3" // или postgres
)

// SQL запросы от учеников
// получаем доступные для записи уроки
func GetAvailableLessons() ([]Lesson, error) {
	query := `
	SELECT s.id, s.name, s.title, s.date
	FROM scheduler s
	LEFT JOIN registrations r ON s.id = r.lesson_id
	WHERE r.lesson_id IS NULL
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

// получаем список дат, в которых есть свободные уроки
func GetDatesWithAvailableLessons() ([]string, error) {
	rows, err := DB.Query(`
		SELECT DISTINCT DATE(date) as day
		FROM scheduler
		WHERE id IN (
			SELECT id FROM scheduler
			EXCEPT
			SELECT lesson_id FROM registrations
		)
		ORDER BY day ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dates []string
	for rows.Next() {
		var date string
		if err := rows.Scan(&date); err != nil {
			return nil, err
		}
		// Опционально: форматировать дату, если нужно отображение "DD.MM"
		dates = append(dates, date)
	}
	return dates, nil
}

// записываем пользователя на урок, если урок ещё свободен
func RegisterUserToLesson(userID int, lessonID int) error {
	_, err := DB.Exec(`
		INSERT INTO registrations (user_id, lesson_id)
		SELECT ?, ?
		WHERE NOT EXISTS (
			SELECT 1 FROM registrations WHERE lesson_id = ?
		)
	`, userID, lessonID, lessonID)

	if err != nil {
		log.Println("Ошибка при регистрации пользователя:", err)
	}

	return err
}

// получаем уроки, на которые записан пользователь
func GetUserLessons(userID int) ([]Lesson, error) {
	query := `
	SELECT s.id, s.name, s.title, s.date
	FROM scheduler s
	INNER JOIN registrations r ON r.lesson_id = s.id
	WHERE r.user_id = ?
	ORDER BY s.date ASC
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

// получаем свободные уроки на определенную дату
func GetLessonsByDate(date string) ([]Lesson, error) {
	query := `
	SELECT s.id, s.name, s.title, s.date
	FROM scheduler s
	LEFT JOIN registrations r ON s.id = r.lesson_id
	WHERE r.lesson_id IS NULL AND s.date LIKE ?
	ORDER BY s.date ASC
	`
	rows, err := DB.Query(query, date+"%")
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

// отменяем регистрацию пользователя на урок
func CancelUserRegistration(userID, lessonID int) error {
	_, err := DB.Exec(`
		DELETE FROM registrations 
		WHERE user_id = ? AND lesson_id = ?
	`, userID, lessonID)

	if err != nil {
		log.Print("Не удалось удалить запись:", err)
	}
	return err
}
