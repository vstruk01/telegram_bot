package vocabulary

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Vocabulary interface {
	Add()
	Delete()
	Repeat()
	Execute(update tgbotapi.Update) (tgbotapi.MessageConfig, error)
}

type Repository interface {
	AddPhrase(phrase string, translate string) error
	AddWord(word, translate string) error
	DeleteWord(word string) error
}

type Performer struct {
	Store  Repository
	States map[int64]*State
}

func New(store Repository, states map[int64]*State) Vocabulary {
	if states == nil {
		states = make(map[int64]*State)
	}

	return &Performer{
		Store:  store,
		States: states,
	}
}

func (p Performer) Add() {
	panic("implement me")
}

func (p Performer) Delete() {
	panic("implement me")
}

func (p Performer) Repeat() {
	panic("implement me")
}

func (p Performer) Execute(update tgbotapi.Update) (tgbotapi.MessageConfig, error) {
	var message = update.Message.Text

	if update.Message.IsCommand() {
		message = fmt.Sprintf("message (%s) is command", message)
	}

	return tgbotapi.NewMessage(update.Message.Chat.ID, message), nil
}
