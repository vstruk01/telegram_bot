package main

import (
	"database/sql"
)

type Request struct {
	Text    string
	Name    string
	Chat_id int
	C       chan string
	ErrC    chan error
	OpenDb   *sql.DB
}

type Master struct {
	Commands map[string]func(Request)
	Rutines  map[int]chan string
	url      string
	offset   int
	errc     chan error
	OpenDb   *sql.DB
}

type Chat struct {
	Id int
}

type User struct {
	Username string `json:"username"`
}

type Message struct {
	Chat Chat
	User User `json:"from"`
	Text string
}

type Update struct {
	Update_id int
	Message   Message `json:"message"`
}

type RestResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type KeyboardButton struct {
	Text string `json:"text"`
}

type ReplyKeyboardMarkup struct {
	Keyboard          [][]KeyboardButton `json:"keyboard"`
	Resize_keyboard   bool               `json:"resize_keyboard"`
	One_time_keyboard bool               `json:"one_time_keyboard"`
	Selective         bool               `json:"selective"`
}

type Button struct {
	Chat_id      int                 `json:"chat_id"`
	Text         string              `json:"text"`
	Reply_markup ReplyKeyboardMarkup `json:"reply_markup"`
}

type BotMessage struct {
	Chat_id int    `json:"chat_id"`
	Text    string `json:"text"`
}
