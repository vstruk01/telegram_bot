package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"telegram_bot/config"
	"telegram_bot/internal/receiver"
	"telegram_bot/internal/store/postgres"
)

func main() {
	conf, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("load config:", err)
	}

	db, err := postgres.NewPostgres(&postgres.Config{
		Host: conf.DB.Host,
		Name: conf.DB.Name,
		Port: conf.DB.Port,
		AppName: conf.DB.AppName,
		Password: conf.DB.Password,
		User: conf.DB.User,
		SourceFiles: conf.DB.SourceFiles,
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	bot, err := tgbotapi.NewBotAPI(conf.Token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)
	receiver.StartReceiver(bot)
}