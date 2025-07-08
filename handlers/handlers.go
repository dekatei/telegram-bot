package handlers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/dekatei/telegram-bot/base"
	"github.com/dekatei/telegram-bot/buttons"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var AddState = map[int]map[string]string{} // userID -> step data
const adminID = 288848928

func StartBot(bot *tgbotapi.BotAPI) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, _ := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.CallbackQuery != nil {
			HandleCallback(bot, update)
			continue
		}
		if update.Message == nil {
			continue
		}

		userID := update.Message.From.ID
		text := update.Message.Text
		chatID := update.Message.Chat.ID

		switch text {
		case "/start":
			msg := tgbotapi.NewMessage(chatID, "Добро пожаловать!")
			msg.ReplyMarkup = buttons.MainMenu(userID)
			bot.Send(msg)

		case "➕ Добавить занятие":
			if userID != adminID {
				bot.Send(tgbotapi.NewMessage(chatID, "⛔️ Только для администратора."))
				break
			}
			AddState[userID] = map[string]string{}
			msg := tgbotapi.NewMessage(chatID, "Выберите дату:")
			msg.ReplyMarkup = dateKeyboard("add_date")
			bot.Send(msg)

		case "📅 Список занятий":
			text, err := buttons.LessonsListMessage()
			if err != nil {
				text = "Ошибка при получении занятий."
			}
			bot.Send(tgbotapi.NewMessage(chatID, text))

		case "✅ Записаться":
			msg := tgbotapi.NewMessage(chatID, "Выберите день:")
			msg.ReplyMarkup = dateKeyboard("register_date")
			bot.Send(msg)

		case "👤 Мои занятия":
			text, err := buttons.MyLessonsMessage(userID)
			if err != nil {
				text = "Ошибка при получении занятий."
			}
			bot.Send(tgbotapi.NewMessage(chatID, text))

		default:
			bot.Send(tgbotapi.NewMessage(chatID, "Выберите действие из меню."))
		}
	}
}

// dateKeyboard возвращает клавиатуру выбора дня
func dateKeyboard(prefix string) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup([]tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("Сегодня", prefix+":0"),
		tgbotapi.NewInlineKeyboardButtonData("Завтра", prefix+":1"),
		tgbotapi.NewInlineKeyboardButtonData("Послезавтра", prefix+":2"),
	})
}

func HandleCallback(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	cb := update.CallbackQuery
	userID := cb.From.ID
	data := cb.Data

	switch {
	case strings.HasPrefix(data, "add_date:"):
		days, _ := strconv.Atoi(strings.TrimPrefix(data, "add_date:"))
		date := time.Now().AddDate(0, 0, days).Format("2006-01-02")
		AddState[userID] = map[string]string{"date": date}

		var rows [][]tgbotapi.InlineKeyboardButton
		for hour := 10; hour <= 21; hour++ {
			label := fmt.Sprintf("%02d:00", hour)
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(label, "add_time:"+label)))
		}
		msg := tgbotapi.NewMessage(cb.Message.Chat.ID, "Выберите время:")
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)
		bot.Send(msg)

	case strings.HasPrefix(data, "add_time:"):
		timeStr := strings.TrimPrefix(data, "add_time:")
		if AddState[userID] == nil {
			AddState[userID] = map[string]string{}
		}
		AddState[userID]["time"] = timeStr

		rows := [][]tgbotapi.InlineKeyboardButton{
			{
				tgbotapi.NewInlineKeyboardButtonData("Онлайн показ", "add_type:Онлайн показ"),
				tgbotapi.NewInlineKeyboardButtonData("Взвод", "add_type:Взвод"),
				tgbotapi.NewInlineKeyboardButtonData("Любое", "add_type:Любое"),
			},
		}
		msg := tgbotapi.NewMessage(cb.Message.Chat.ID, "Выберите тип занятия:")
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)
		bot.Send(msg)

	case strings.HasPrefix(data, "add_type:"):
		typeStr := strings.TrimPrefix(data, "add_type:")
		info := AddState[userID]
		if info == nil {
			bot.Send(tgbotapi.NewMessage(cb.Message.Chat.ID, "⚠️ Данные не найдены, начните сначала."))
			return
		}
		fullDate := fmt.Sprintf("%s %s", info["date"], info["time"])
		err := base.AddLesson("Занятие", typeStr, fullDate)
		if err != nil {
			bot.Send(tgbotapi.NewMessage(cb.Message.Chat.ID, "❌ Ошибка добавления: "+err.Error()))
		} else {
			bot.Send(tgbotapi.NewMessage(cb.Message.Chat.ID, "✅ Добавлено: "+typeStr+" — "+fullDate))
		}
		delete(AddState, userID)

	case strings.HasPrefix(data, "register_date:"):
		days, _ := strconv.Atoi(strings.TrimPrefix(data, "register_date:"))
		dateStr := time.Now().AddDate(0, 0, days).Format("2006-01-02")

		lessons, err := base.GetLessonsByDate(dateStr)
		if err != nil || len(lessons) == 0 {
			bot.Send(tgbotapi.NewMessage(cb.Message.Chat.ID, "Нет доступных занятий на эту дату."))
			break
		}

		var rows [][]tgbotapi.InlineKeyboardButton
		for _, l := range lessons {
			label := fmt.Sprintf("%s — %s", l.Title, l.Date[11:])
			callbackData := fmt.Sprintf("register:%d", l.ID)
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(label, callbackData)))
		}

		msg := tgbotapi.NewMessage(cb.Message.Chat.ID, fmt.Sprintf("Занятия на %s:", dateStr))
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)
		bot.Send(msg)

	case strings.HasPrefix(data, "register:"):
		id, _ := strconv.Atoi(strings.TrimPrefix(data, "register:"))
		err := base.RegisterUserToLesson(cb.From.ID, id)
		text := "✅ Вы успешно записаны!"
		if err != nil {
			text = "⚠️ Не удалось записаться: " + err.Error()
		}
		bot.Send(tgbotapi.NewMessage(cb.Message.Chat.ID, text))
	}

	bot.AnswerCallbackQuery(tgbotapi.NewCallback(cb.ID, ""))
}
