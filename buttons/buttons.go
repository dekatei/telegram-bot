package buttons

import (
	"fmt"

	"github.com/dekatei/telegram-bot/base"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func MainMenu() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("📅 Список занятий"),
			tgbotapi.NewKeyboardButton("✅ Записаться"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("❌ Отменить запись"),
			tgbotapi.NewKeyboardButton("👤 Мои занятия"),
		),
	)
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
