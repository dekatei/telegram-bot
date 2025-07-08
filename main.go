package main

import (
	"log"
	"os"

	"github.com/dekatei/telegram-bot/base"
	"github.com/dekatei/telegram-bot/handlers"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
)

const AdminID int = 288848928

func main() {
	_ = godotenv.Load()

	err := base.InitDB(os.Getenv("DB_FILE"))
	if err != nil {
		log.Fatal("Ошибка инициализации БД:", err)
	}

	token := os.Getenv("TELEGRAM_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_TOKEN не найден в .env")
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal("Ошибка запуска бота:", err)
	}

	bot.Debug = true
	log.Printf("✅ Бот авторизован как @%s", bot.Self.UserName)

	handlers.StartBot(bot)
}
