package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	// "strconv"

	// "os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	database, _ := sql.Open("sqlite3", "./words.db")
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS words (id INTEGER PRIMARY KEY, name TEXT, word TEXT, transcription TEXT, translate TEXT, ok int, notok int)")
	statement.Exec()
	// statement, _ = database.Prepare("CREATE TABLE IF NOT EXISTS people (id INTEGER PRIMARY KEY, firstname TEXT, lastname TEXT)")
	// statement.Exec()
	// statement, _ = database.Prepare("INSERT INTO people (firstname, lastname) VALUES (?, ?)")
	// statement.Exec("Nic", "Raboy")
	// rows, e := database.Query("SELECT id, firstname, lastname FROM people where id = ?", 5)
	// if e != nil {
	// 	log.Println("Error rows")
	// 	os.Exit(1)
	// }
	// var id int
	// var firstname string
	// var lastname string
	// for rows.Next() {
	// 	rows.Scan(&id, &firstname, &lastname)
	// 	fmt.Println(strconv.Itoa(id) + ": " + firstname + " " + lastname)
	// }

	url := "https://api.telegram.org/bot1060785017:AAG7eJUSygisjIF_g97Dj5TKVzS-ct76su8/"

	offset := 0
	for {
		err := getUpdate(url, &offset)
		if err != nil {
			log.Println("Error update", err.Error())
		}
	}
}

func respondPrintWords(name string, chat_id int) {
	var message string

	database, _ := sql.Open("sqlite3", "./words.db")
	rows, _ := database.Query("select word, translate from words where name = ?", name)
	var word, translate string
	for rows.Next() {
		rows.Scan(&word, &translate)
		message += word + "---" + translate
	}
	if message == "" {
		message += "there are not words"
	}
	sendMessage(message, chat_id)
}

// func respondWords(url string, eWords map[string]string, char_id int) error {
// 	var message BotMessage
// 	message.Chat_id = char_id
// 	for k, v := range eWords {
// 		message.Text += k + " - " + v + "\n"
// 	}
// 	respond(url, message)
// 	return nil
// }

// func respondEcho(url string, update Update) error {
// 	var message BotMessage
// 	message.Chat_id = update.Message.Chat.Id

// 	message.Text = update.Message.Text
// 	respond(url, message)
// 	return nil
// }

func sendMessage(message string, chat_id int) error {
	var m BotMessage
	m.Chat_id = chat_id
	m.Text = message

	buf, err := json.Marshal(m)
	if err != nil {
		return err
	}
	_, err = http.Post("https://api.telegram.org/bot1060785017:AAG7eJUSygisjIF_g97Dj5TKVzS-ct76su8/sendMessage", "application/json", bytes.NewBuffer(buf))
	if err != nil {
		return err
	}
	return nil
}
