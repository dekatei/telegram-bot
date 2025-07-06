package main

import (
	"log"
	"os"

	"github.com/dekatei/telegram-bot/base"
	"github.com/dekatei/telegram-bot/buttons"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
)

const adminID int = 288848928

func main() {

	// Загружаем .env файл
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbFile := os.Getenv("DB_FILE")
	err = base.InitDB(dbFile)
	if err != nil {
		log.Fatal("Ошибка инициализации базы данных:", err)
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
		log.Printf("Сообщение от %s: %s", update.Message.From.UserName, update.Message.Text)
		switch update.Message.Text {
		case "/start":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Добро пожаловать!")
			msg.ReplyMarkup = buttons.MainMenu()
			bot.Send(msg)
			//log.Printf("Ваш Telegram ID: %d", update.Message.From.ID)
		case "Добавить новое время занятий":
			if update.Message.From.ID != adminID {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "⛔️ Доступ запрещён."))
				break
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите название нового занятия:")
			bot.Send(msg)

			// Ждём следующее сообщение от админа
			update2 := <-updates
			name := update2.Message.Text

			msg2 := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите описание занятия:")
			bot.Send(msg2)
			update3 := <-updates
			title := update3.Message.Text

			msg3 := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите дату и время занятия (например: 2025-07-08 15:00):")
			bot.Send(msg3)
			update4 := <-updates
			date := update4.Message.Text

			err := base.AddLesson(name, title, date)
			if err != nil {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка при добавлении занятия."))
			} else {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "✅ Занятие успешно добавлено!"))
			}
		case "📅 Список занятий":
			text, err := buttons.LessonsListMessage()
			if err != nil {
				text = "Ошибка при получении занятий."
			}
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, text))
		case "✅ Записаться":
			text, err := buttons.RegisterMessage(update.Message.From.ID)
			if err != nil {
				text = "Ошибка при записи."
			}
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, text))

		case "👤 Мои занятия":
			text, err := buttons.MyLessonsMessage(update.Message.From.ID)
			if err != nil {
				text = "Ошибка при получении занятий."
			}
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, text))
		default:
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите действие из меню."))
		}
	}
}
