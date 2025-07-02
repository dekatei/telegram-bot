package main

import (
	"log"
	"os"

	"github.com/dekatei/telegram-bot/buttons"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
)

func main() {
	// Загружаем .env файл
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	//получаем токен из .env файла
	token := os.Getenv("TELEGRAM_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_TOKEN not set in .env")
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		switch update.Message.Text {
		case "/start":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Добро пожаловать!")
			msg.ReplyMarkup = buttons.MainMenu()
			bot.Send(msg)

		case "📅 Список занятий":
			text, err := buttons.LessonsListMessage()
			if err != nil {
				text = "Ошибка при получении занятий."
			}
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, text))

		default:
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите действие из меню."))
		}
	}
}
