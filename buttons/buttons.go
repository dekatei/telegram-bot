package buttons

import (
	"fmt"

	"github.com/dekatei/telegram-bot/base"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const AdminID int = 288848928

// MainMenu возвращает главное меню с дополнительной кнопкой для администратора.
func MainMenu(userID int) tgbotapi.ReplyKeyboardMarkup {
	rows := [][]tgbotapi.KeyboardButton{
		{
			tgbotapi.NewKeyboardButton("📅 Список занятий"),
			tgbotapi.NewKeyboardButton("✅ Записаться"),
		},
		{
			tgbotapi.NewKeyboardButton("❌ Отменить запись"),
			tgbotapi.NewKeyboardButton("👤 Мои занятия"),
		},
	}

	if userID == AdminID {
		rows = append(rows, tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("➕ Добавить занятие"),
			tgbotapi.NewKeyboardButton("➖ Удалить доступное занятие"),
		))
	}
	return tgbotapi.NewReplyKeyboard(rows...)
}

// LessonsListMessage возвращает список доступных занятий текстом.
func LessonsListMessage() (string, error) {
	lessons, err := base.GetAvailableLessons()
	if err != nil {
		return "", err
	}

	if len(lessons) == 0 {
		return "Нет доступных занятий.", nil
	}

	msg := "📅 Доступные занятия:\n"
	for _, l := range lessons {
		msg += fmt.Sprintf("🔹 %s — %s\n", l.Name, l.Date)
	}

	return msg, nil
}

// MyLessonsMessage возвращает список занятий конкретного пользователя.
func MyLessonsMessage(userID int) (string, error) {
	var lessons []base.Lesson
	var err error
	if userID == AdminID {
		lessons, err = base.GetAdminLessons()
	} else {
		lessons, err = base.GetUserLessons(userID)
	}

	if err != nil {
		return "Ошибка при получении ваших занятий.", err
	}
	if len(lessons) == 0 {
		return "У вас нет записей.", nil
	}

	msg := "👤 Ваши записи:\n"
	for _, l := range lessons {
		msg += fmt.Sprintf("🔸 %s — %s\n", l.Name, l.Date)
	}
	return msg, nil
}
