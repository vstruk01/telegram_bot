package main

import (
	"database/sql"

	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	err := createDB()
	if err != nil {
		return
	}

	functions := make(map[string]func(Request) error)
	functions["/start"] = CommandStart
	functions["RepeatKnow"] = CommandStart
	functions["ListKnow"] = CommandStart
	functions["WordKnow"] = CommandStart
	functions["WordNew"] = CommandStart
	functions["ListNew"] = CommandStart
	functions["RepeatNew"] = CommandStart

	stmt_list := make(map[string]sql.Rows)
	actions_list := make(map[string]int)

	url := "https://api.telegram.org/bot1060785017:AAG7eJUSygisjIF_g97Dj5TKVzS-ct76su8/"
	offset := 0

	for {
		err = getUpdate(url, &offset, &stmt_list, &actions_list, functions)
		if err != nil {
			fmt.Println("\033[1;32mError = ", err.Error(), "\033[0m")
		}
	}
}

func createDB() error {
	database, err := sql.Open("sqlite3", "./words.db")
	if err != nil {
		fmt.Print("\033[1;34createDB open\033[0m\n")
		return err
	}
	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS words (id INTEGER PRIMARY KEY, name TEXT, word TEXT, transcription TEXT, translate TEXT, ok int)")
	if err != nil {
		fmt.Print("\033[1;34createDB words\033[0m\n")
		return err
	}
	_, err = statement.Exec()
	if err != nil {
		fmt.Print("\033[1;34createDB open\033[0m\n")
		return err
	}
	statement, err = database.Prepare("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, name TEXT, chat_id INT)")
	if err != nil {
		fmt.Print("\033[1;34createDB users\033[0m\n")
		return err
	}
	_, err = statement.Exec()
	if err != nil {
		fmt.Print("\033[1;34createDB users\033[0m\n")
		return err
	}
	return nil
}
