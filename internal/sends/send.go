package sends

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type Request struct {
	Text    string
	Name    string
	Chat_id int
	OpenDb  *sql.DB
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

func TranslateWord(r Request) error {
	database, err := sql.Open("sqlite3", "./words.db")
	if err != nil {
		fmt.Print("\033[1;34mError Translate Word 1\033[0m\n")
		return err
	}
	rows, err := database.Query("select word, translate from words where name = ? and word = ?", r.Name, r.Text)
	if err != nil {
		fmt.Print("\033[1;34mError Translate Word 2\033[0m\n")
		return err
	}
	err = SendWords(rows, r.Chat_id)
	if err != nil {
		return err
	}
	fmt.Print("\033[1;34mTranslate Word Ok\033[0m\n")
	return nil
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
	SendMessage(message, chat_id)
	fmt.Print("\033[1;34msendWords Ok\033[0m\n")
	return nil
}

func SetButton(chat_id int) error {
	var m Button
	m.Chat_id = chat_id
	buttonAll := make([][]KeyboardButton, 3)
	buttonOne := make([]KeyboardButton, 2)
	buttonTwo := make([]KeyboardButton, 2)
	buttonThree := make([]KeyboardButton, 2)
	buttonOne[0].Text = "RepeatNew"
	buttonOne[1].Text = "RepeatKnow"
	buttonAll[0] = buttonOne
	buttonTwo[0].Text = "WordNew"
	buttonTwo[1].Text = "WordKnow"
	buttonAll[1] = buttonTwo
	buttonThree[0].Text = "ListNew"
	buttonThree[1].Text = "ListKnow"
	buttonAll[2] = buttonThree
	m.Reply_markup.Keyboard = buttonAll
	m.Reply_markup.Resize_keyboard = true
	m.Reply_markup.One_time_keyboard = true
	m.Reply_markup.Selective = true
	m.Text = "set keyboard"

	buf, err := json.Marshal(m)
	if err != nil {
		fmt.Print("\033[1;34mSet Button\033[0m\n")
		return err
	}
	_, err = http.Post("https://api.telegram.org/bot1060785017:AAG7eJUSygisjIF_g97Dj5TKVzS-ct76su8/sendMessage", "application/json", bytes.NewBuffer(buf))
	if err != nil {
		fmt.Print("\033[1;34mSet Button\033[0m\n")
		return err
	}
	fmt.Print("\033[1;34mSet Button Ok\033[0m\n")
	return nil
}

func SendMessage(message string, chat_id int) error {
	var m BotMessage
	m.Chat_id = chat_id
	m.Text = message

	buf, err := json.Marshal(m)
	if err != nil {
		fmt.Print("\033[1;34mSend Message\033[0m\n")
		return err
	}
	_, err = http.Post("https://api.telegram.org/bot1060785017:AAG7eJUSygisjIF_g97Dj5TKVzS-ct76su8/sendMessage", "application/json", bytes.NewBuffer(buf))
	if err != nil {
		fmt.Print("\033[1;34mSend Message\033[0m\n")
		return err
	}
	fmt.Print("\033[1;34mSend Message Ok\033[0m\n")
	return nil
}
