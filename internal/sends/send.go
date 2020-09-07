package sends

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	log "github.com/vstruk01/telegram_bot/internal/Logger"
	botStruct "github.com/vstruk01/telegram_bot/internal/struct"
)

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

func SendWords(rows *sql.Rows, chat_id int) error {
	var message string

	var word, translate string
	for rows.Next() {
		rows.Scan(&word, &translate)
		translates := strings.Split(translate, ",")
		message += word + "  ->  "
		for i := 0; i < len(translates); i++ {
			message += translates[i] + "   "
		}
		message += "\n"
	}
	if message == "" {
		message += "Hmmmmm I think that you wrong"
	}
	err := SendMessage(message, chat_id)
	if log.CheckErr(err) {
		return err
	}
	return nil
}

func SetButton(chat_id int) error {
	var m Button
	m.Chat_id = chat_id
	buttonAll := make([][]KeyboardButton, 4)
	buttonOne := make([]KeyboardButton, 2)
	buttonTwo := make([]KeyboardButton, 2)
	buttonThree := make([]KeyboardButton, 2)
	buttonFour := make([]KeyboardButton, 2)
	buttonOne[0].Text = "RepeatNew"
	buttonOne[1].Text = "RepeatKnow"
	buttonAll[0] = buttonOne
	buttonTwo[0].Text = "WordNew"
	buttonTwo[1].Text = "WordKnow"
	buttonAll[1] = buttonTwo
	buttonThree[0].Text = "ListNew"
	buttonThree[1].Text = "ListKnow"
	buttonAll[2] = buttonThree
	buttonFour[0].Text = "DeleteWord"
	buttonFour[1].Text = "AddWord"
	buttonAll[3] = buttonFour
	m.Reply_markup.Keyboard = buttonAll
	m.Reply_markup.Resize_keyboard = true
	m.Reply_markup.One_time_keyboard = true
	m.Reply_markup.Selective = true
	m.Text = "set keyboard"

	buf, err := json.Marshal(m)
	if log.CheckErr(err) {
		return err
	}
	_, err = http.Post(botStruct.Url+botStruct.Token+"/sendMessage", "application/json", bytes.NewBuffer(buf))
	if log.CheckErr(err) {
		return err
	}
	return nil
}

func SendMessage(message string, chat_id int) error {
	var m BotMessage
	m.Chat_id = chat_id
	m.Text = message

	buf, err := json.Marshal(m)
	if log.CheckErr(err) {
		return err
	}
	_, err = http.Post(botStruct.Url+botStruct.Token+"/sendMessage", "application/json", bytes.NewBuffer(buf))
	if log.CheckErr(err) {
		return err
	}
	return nil
}
