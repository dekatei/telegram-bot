package handlers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/dekatei/telegram-bot/base"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/mattn/go-sqlite3" // или postgres
)

var AddState = map[int]map[string]string{} // userID -> step data

func HandleCallback(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	cb := update.CallbackQuery
	userID := cb.From.ID
	data := cb.Data

	switch {
	case strings.HasPrefix(data, "add_date:"):
		days, _ := strconv.Atoi(strings.TrimPrefix(data, "add_date:"))
		date := time.Now().AddDate(0, 0, days).Format("2006-01-02")
		AddState[userID]["date"] = date

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
