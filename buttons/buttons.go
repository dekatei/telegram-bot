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
