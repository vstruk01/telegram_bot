package receiver

import (
	"fmt"
	"log"
	"telegram_bot/internal/vocabulary"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	ErrMessage = "something wend wrong: %s"
)

type Receiver struct {
	Vocabulary vocabulary.Vocabulary
	Bot        *tgbotapi.BotAPI
}

func New(vocabulary vocabulary.Vocabulary, bot *tgbotapi.BotAPI) *Receiver {
	return &Receiver{
		Vocabulary: vocabulary,
		Bot:        bot,
	}
}

func (r *Receiver) StartReceiver() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := r.Bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		message, err := r.Vocabulary.Execute(update)
		if err != nil {
			message = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf(ErrMessage, err.Error()))
		}

		_, err = r.Bot.Send(message)
		if err != nil {
			log.Println(err.Error())
		}
	}
}
