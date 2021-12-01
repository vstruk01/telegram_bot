package receiver

import (
	"database/sql"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Receiver struct {
	DB *sql.DB
	Bot *tgbotapi.BotAPI
}

func New() *Receiver {
	return &Receiver{

	}
}

func (r *Receiver)StartReceiver() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := r.Bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		log.Println("is command", update.Message.IsCommand())

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		_, err := r.Bot.Send(msg)
		if err != nil {
			log.Println(err.Error())
		}
	}
}
