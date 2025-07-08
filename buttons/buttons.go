package buttons

import (
	"fmt"

	"github.com/dekatei/telegram-bot/base"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const adminID int = 288848928

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

	if userID == adminID {
		rows = append(rows, tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("➕ Добавить занятие"),
		))
	}
	return tgbotapi.NewReplyKeyboard(rows...)
}

func LessonsListMessage() (string, error) {
	lessons, err := base.GetAllLessons()
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

func RegisterMessage(userID int) (string, error) {
	lessons, err := base.GetAllLessons()
	if err != nil || len(lessons) == 0 {
		return "Нет доступных занятий.", err
	}

	// для простоты: записываем на первое занятие
	selected := lessons[0]

	err = base.RegisterUserToLesson(userID, selected.ID)
	if err != nil {
		return "Ошибка при записи на занятие.", err
	}

	return fmt.Sprintf("✅ Вы записаны на занятие: %s — %s", selected.Name, selected.Date), nil
}

func MyLessonsMessage(userID int) (string, error) {
	lessons, err := base.GetUserLessons(userID)
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
