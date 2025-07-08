package main

import (
	"log"
	"os"

	"github.com/dekatei/telegram-bot/base"
	"github.com/dekatei/telegram-bot/buttons"
	"github.com/dekatei/telegram-bot/handlers"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
)

const AdminID int = 288848928

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
		if update.CallbackQuery != nil {
			handlers.HandleCallback(bot, update)
			continue
		}
		if update.Message == nil {
			continue
		}
		log.Printf("Сообщение от %s: %s", update.Message.From.UserName, update.Message.Text)
		userID := update.Message.From.ID

		switch update.Message.Text {
		case "/start":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Добро пожаловать!")
			msg.ReplyMarkup = buttons.MainMenu(update.Message.From.ID)
			bot.Send(msg)
			//log.Printf("Ваш Telegram ID: %d", update.Message.From.ID)
		case "➕ Добавить занятие":
			if userID != AdminID {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "⛔️ Только для администратора."))
				break
			}

			handlers.AddState[userID] = map[string]string{}
			rows := [][]tgbotapi.InlineKeyboardButton{
				{
					tgbotapi.NewInlineKeyboardButtonData("Сегодня", "add_date:0"),
					tgbotapi.NewInlineKeyboardButtonData("Завтра", "add_date:1"),
					tgbotapi.NewInlineKeyboardButtonData("Послезавтра", "add_date:2"),
				},
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите дату:")
			msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)
			bot.Send(msg)
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
