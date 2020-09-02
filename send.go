package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func sendWord(text string, name string, chat_id int) error {
	database, err := sql.Open("sqlite3", "./words.db")
	if err != nil {
		fmt.Print("\033[1;34mError sendWord 1\033[0m\n")
		return err
	}
	rows, err := database.Query("select translate from words where name = ? and word = ?", name, text)
	if err != nil {
		fmt.Print("\033[1;34mError sendWord 2\033[0m\n")
		return err
	}
	var translate string

	rows.Next()
	err = rows.Scan(&translate)
	if err != nil {
		err = sendMessage("Hmmmmm I think that you wrong", chat_id)
		if err != nil {
			return err
		}
	} else {
		err = sendMessage(translate, chat_id)
		if err != nil {
			return err
		}
	}
	err = rows.Close()
	if err != nil {
		fmt.Print("\033[1;34mError sendWord 3\033[0m\n")
		return err
	}
	fmt.Print("\033[1;34msendWord Ok\033[0m\n")
	return nil
}

func sendWords(rows *sql.Rows, chat_id int) error {
	var message string

	var word, translate string
	for rows.Next() {
		rows.Scan(&word, &translate)
		translates := strings.Split(translate, ",")
		message += word + " ->   "
		for i := 0; i < len(translates); i++ {
			message += translates[i] + "   "
		}
		message += "\n"
	}
	if message == "" {
		message += "there are not words"
	}
	sendMessage(message, chat_id)
	fmt.Print("\033[1;34msendWords Ok\033[0m\n")
	return nil
}

func setButton(chat_id int) error {
	var m Button
	m.Chat_id = chat_id
	buttonAll := make([][]KeyboardButton, 3)
	buttonOne := make([]KeyboardButton, 2)
	buttonTwo := make([]KeyboardButton, 2)
	buttonThree := make([]KeyboardButton, 2)
	buttonOne[0].Text = "repeatNew"
	buttonOne[1].Text = "repeatKnow"
	buttonAll[0] = buttonOne
	buttonTwo[0].Text = "Know"
	buttonTwo[1].Text = "NotKnow"
	buttonAll[1] = buttonTwo
	buttonThree[0].Text = "listKnow"
	buttonThree[1].Text = "listNew"
	buttonAll[2] = buttonThree
	m.Reply_markup.Keyboard = buttonAll
	m.Reply_markup.Resize_keyboard = true
	m.Reply_markup.One_time_keyboard = true
	m.Reply_markup.Selective = true
	m.Text = "set keyboard"

	buf, err := json.Marshal(m)
	if err != nil {
		fmt.Print("\033[1;34mset Button\033[0m\n")
		return err
	}
	_, err = http.Post("https://api.telegram.org/bot1060785017:AAG7eJUSygisjIF_g97Dj5TKVzS-ct76su8/sendMessage", "application/json", bytes.NewBuffer(buf))
	if err != nil {
		fmt.Print("\033[1;34mset Button\033[0m\n")
		return err
	}
	fmt.Print("\033[1;34msetButton Ok\033[0m\n")
	return nil
}

func sendMessage(message string, chat_id int) error {
	var m BotMessage
	m.Chat_id = chat_id
	m.Text = message

	buf, err := json.Marshal(m)
	if err != nil {
		fmt.Print("\033[1;34msendMessage\033[0m\n")
		return err
	}
	_, err = http.Post("https://api.telegram.org/bot1060785017:AAG7eJUSygisjIF_g97Dj5TKVzS-ct76su8/sendMessage", "application/json", bytes.NewBuffer(buf))
	if err != nil {
		fmt.Print("\033[1;34msendMessage\033[0m\n")
		return err
	}
	fmt.Print("\033[1;34msendMessage Ok\033[0m\n")
	return nil
}
