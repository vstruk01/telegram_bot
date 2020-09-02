package main

import (
	"database/sql"

	// "fmt"
	"log"

	// "strconv"

	// "os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	database, _ := sql.Open("sqlite3", "./words.db")
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS words (id INTEGER PRIMARY KEY, name TEXT, word TEXT, transcription TEXT, translate TEXT, ok int)")
	statement.Exec()
	stmt_list := make(map[string]sql.Rows)
	actions_list := make(map[string]int)

	url := "https://api.telegram.org/bot1060785017:AAG7eJUSygisjIF_g97Dj5TKVzS-ct76su8/"

	offset := 0
	for {
		err := getUpdate(url, &offset, &stmt_list, &actions_list)
		if err != nil {
			log.Println("Error update", err.Error())
		}
	}
}
