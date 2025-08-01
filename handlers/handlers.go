package handlers

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/dekatei/telegram-bot/base"
	"github.com/dekatei/telegram-bot/buttons"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var AddState = map[int]map[string]interface{}{} // userID -> step data
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

		case "Добавить занятие":
			if userID != adminID {
				bot.Send(tgbotapi.NewMessage(chatID, "⛔️ Только для администратора."))
				break
			}
			AddState[userID] = map[string]interface{}{}
			msg := tgbotapi.NewMessage(chatID, "Выберите дату:")
			msg.ReplyMarkup = dateKeyboardForAdd("add_date")
			bot.Send(msg)
		case "Удалить доступное занятие":
			if userID != adminID {
				bot.Send(tgbotapi.NewMessage(chatID, "⛔️ Только для администратора."))
				break
			}
			lessons, err := base.GetAvailableLessons()
			if err != nil || len(lessons) == 0 {
				bot.Send(tgbotapi.NewMessage(chatID, "У вас нет записей."))
				break
			}

			var rows [][]tgbotapi.InlineKeyboardButton
			for _, l := range lessons {
				label := fmt.Sprintf("%s — %s", l.Title, l.Date[11:])
				callbackData := fmt.Sprintf("delete_lesson:%d", l.ID)
				rows = append(rows, tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData(label, callbackData)))
			}

			msg := tgbotapi.NewMessage(chatID, "Выберите занятие для удаления:")
			msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)
			bot.Send(msg)

		case "📅 Свободные занятия":
			text, err := buttons.LessonsListMessage()
			if err != nil {
				text = "Ошибка при получении занятий."
			}
			bot.Send(tgbotapi.NewMessage(chatID, text))

		case "✅ Записаться":
			msg := tgbotapi.NewMessage(chatID, "Выберите день:")
			msg.ReplyMarkup = dateKeyboardForRegistration("register_date")
			bot.Send(msg)

		case "👤 Мои занятия":
			text, err := buttons.MyLessonsMessage(userID)
			if err != nil {
				text = "Ошибка при получении занятий."
			}
			bot.Send(tgbotapi.NewMessage(chatID, text))
		case "❌ Отменить запись":
			var lessons []base.Lesson
			var err error
			if userID == adminID {
				lessons, err = base.GetAdminLessons()
			} else {
				lessons, err = base.GetUserLessons(userID)
			}

			if err != nil || len(lessons) == 0 {
				bot.Send(tgbotapi.NewMessage(chatID, "У вас нет записей."))
				break
			}

			var rows [][]tgbotapi.InlineKeyboardButton
			for _, l := range lessons {
				label := fmt.Sprintf("%s — %s", l.Title, l.Date[11:])
				callbackData := fmt.Sprintf("cancel_lesson:%d", l.ID)
				rows = append(rows, tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData(label, callbackData)))
			}

			msg := tgbotapi.NewMessage(chatID, "Выберите занятие для удаления:")
			msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)
			bot.Send(msg)

		default:
			bot.Send(tgbotapi.NewMessage(chatID, "Выберите действие из меню."))
		}
	}
}

func formatDate(d time.Time) string {
	months := [...]string{
		"января", "февраля", "марта", "апреля", "мая", "июня",
		"июля", "августа", "сентября", "октября", "ноября", "декабря",
	}
	return fmt.Sprintf("%02d %s", d.Day(), months[d.Month()-1])
}

// dateKeyboard возвращает клавиатуру выбора дня с доступными уроками
func dateKeyboardForRegistration(prefix string) tgbotapi.InlineKeyboardMarkup {
	dates, err := base.GetDatesWithAvailableLessons()
	if err != nil {
		log.Println("Ошибка получения доступных дат:", err)
		return tgbotapi.NewInlineKeyboardMarkup()
	}

	var rows [][]tgbotapi.InlineKeyboardButton
	for _, d := range dates {
		data := fmt.Sprintf("%s:%s", prefix, d) // Передаём дату строкой
		parsed, _ := time.Parse("2006-01-02", d)
		label := parsed.Format("02.01") // или "02 января"
		btn := tgbotapi.NewInlineKeyboardButtonData(label, data)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn))
	}

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func dateKeyboardForAdd(prefix string) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton
	today := time.Now()
	for i := 0; i < 7; i++ {
		date := today.AddDate(0, 0, i)
		formatted := formatDate(date) // "09 июля"
		data := fmt.Sprintf("%s:%d", prefix, i)
		btn := tgbotapi.NewInlineKeyboardButtonData(formatted, data)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn))
	}
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func HandleCallback(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	cb := update.CallbackQuery
	userID := cb.From.ID
	data := cb.Data

	switch {
	case strings.HasPrefix(data, "add_date:"):
		days, _ := strconv.Atoi(strings.TrimPrefix(data, "add_date:"))
		selectedDate := time.Now().AddDate(0, 0, days)
		AddState[userID] = map[string]interface{}{
			"date":  selectedDate,
			"times": []string{},
		}

		rows := [][]tgbotapi.InlineKeyboardButton{}
		for hour := 10; hour <= 21; hour++ {
			label := fmt.Sprintf("%02d:00", hour)
			callback := "add_time_multi:" + label
			btn := tgbotapi.NewInlineKeyboardButtonData(label, callback)

			if hour%3 == 1 {
				rows = append(rows, []tgbotapi.InlineKeyboardButton{})
			}
			rows[len(rows)-1] = append(rows[len(rows)-1], btn)
		}
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Готово", "add_time_done"),
		))

		msg := tgbotapi.NewMessage(cb.Message.Chat.ID, "Выберите одно или несколько времен:")
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)
		bot.Send(msg)
	case strings.HasPrefix(data, "add_time_multi:"):
		timeStr := strings.TrimPrefix(data, "add_time_multi:")
		info := AddState[userID]
		if info == nil {
			bot.Send(tgbotapi.NewMessage(cb.Message.Chat.ID, "⚠️ Начните с выбора даты."))
			return
		}

		times := info["times"].([]string)
		found := false
		for i, t := range times {
			if t == timeStr {
				// Удаляем
				info["times"] = append(times[:i], times[i+1:]...)
				found = true
				break
			}
		}
		if !found {
			info["times"] = append(times, timeStr)
		}

		// Перерисовываем клавиатуру с отметками
		var rows [][]tgbotapi.InlineKeyboardButton
		times = info["times"].([]string)

		for hour := 10; hour <= 21; hour++ {
			label := fmt.Sprintf("%02d:00", hour)
			display := label
			for _, t := range times {
				if t == label {
					display = "✅ " + label
					break
				}
			}
			btn := tgbotapi.NewInlineKeyboardButtonData(display, "add_time_multi:"+label)

			if hour%3 == 1 {
				rows = append(rows, []tgbotapi.InlineKeyboardButton{})
			}
			rows[len(rows)-1] = append(rows[len(rows)-1], btn)
		}
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Готово", "add_time_done"),
		))

		edit := tgbotapi.NewEditMessageReplyMarkup(
			cb.Message.Chat.ID,
			cb.Message.MessageID,
			tgbotapi.NewInlineKeyboardMarkup(rows...),
		)
		bot.Send(edit)

	case data == "add_time_done":
		info := AddState[userID]
		if info == nil || len(info["times"].([]string)) == 0 {
			bot.Send(tgbotapi.NewMessage(cb.Message.Chat.ID, "❗ Сначала выберите хотя бы одно время."))
			return
		}

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

	case strings.HasPrefix(data, "add_time:"):
		timeStr := strings.TrimPrefix(data, "add_time:")
		if AddState[userID] == nil {
			AddState[userID] = map[string]interface{}{}
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

		times := info["times"].([]string)
		var messages []string

		date := info["date"].(time.Time)
		for _, t := range times {
			fullDateTimeStr := fmt.Sprintf("%s %s", date.Format("2006-01-02"), t)
			fullDateTime, err := time.Parse("2006-01-02 15:04", fullDateTimeStr)
			if err != nil {
				messages = append(messages, fmt.Sprintf("❌ %s — неверный формат даты", t))
				continue
			}
			err = base.AddLesson("Занятие", typeStr, fullDateTime.Format("2006-01-02 15:04:05"))
			if err != nil {
				messages = append(messages, fmt.Sprintf("❌ %s — ошибка: %v", t, err))
			} else {
				messages = append(messages, fmt.Sprintf("✅ %s", t))
			}
		}
		result := "Добавлены занятия на " + date.Format("02.01.2006") + " в :\n" + strings.Join(messages, "\n")
		bot.Send(tgbotapi.NewMessage(cb.Message.Chat.ID, result))
		delete(AddState, userID)

	case strings.HasPrefix(data, "register_date:"):
		dateStr := strings.TrimPrefix(data, "register_date:")

		lessons, err := base.GetLessonsByDate(dateStr)
		if err != nil || len(lessons) == 0 {
			bot.Send(tgbotapi.NewMessage(cb.Message.Chat.ID, "Нет доступных занятий на эту дату."))
			break
		}

		var rows [][]tgbotapi.InlineKeyboardButton
		for _, l := range lessons {
			timeOnly := l.Date[11:] // "18:00" если формат "2006-01-02 15:04"
			label := fmt.Sprintf("%s — %s", l.Title, timeOnly)
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
	case strings.HasPrefix(data, "delete_lesson:"):
		idStr := strings.TrimPrefix(data, "delete_lesson:")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			bot.Send(tgbotapi.NewMessage(cb.Message.Chat.ID, "❌ Неверный ID."))
			return
		}

		err = base.DeleteLesson(id)
		if err != nil {
			bot.Send(tgbotapi.NewMessage(cb.Message.Chat.ID, "❌ Ошибка при удалении: "+err.Error()))
		} else {
			bot.Send(tgbotapi.NewMessage(cb.Message.Chat.ID, "✅ Урок успешно удалён."))
		}
	case strings.HasPrefix(data, "cancel_lesson:"):
		idStr := strings.TrimPrefix(data, "cancel_lesson:")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			bot.Send(tgbotapi.NewMessage(cb.Message.Chat.ID, "❌ Неверный ID."))
			return
		}

		err = base.CancelUserRegistration(userID, id)
		if err != nil {
			bot.Send(tgbotapi.NewMessage(cb.Message.Chat.ID, "❌ Ошибка при отмене: "+err.Error()))
		} else {
			bot.Send(tgbotapi.NewMessage(cb.Message.Chat.ID, "✅ Урок успешно отменён."))
		}

	}

	bot.AnswerCallbackQuery(tgbotapi.NewCallback(cb.ID, ""))

}
