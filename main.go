package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	master, err := InitAll()
	if err != nil {
		fmt.Println("\033[1;32mError = ", err.Error(), "\033[0m")
		return
	}

	for {
		err = GetUpdate(master)
		if err != nil {
			fmt.Println("\033[1;32mError = ", err.Error(), "\033[0m")
		}
	}
}

func InitAll() (*Master, error) {
	var err error
	master := new(Master)

	master.HandeFunc("/start", CommandStart)
	master.HandeFunc("RepeatKnow", CommandRepeatKnow)
	master.HandeFunc("ListKnow", CommandListKnow)
	master.HandeFunc("WordKnow", CommandWordKnow)
	master.HandeFunc("WordNew", CommandWordNew)
	master.HandeFunc("ListNew", CommandListNew)
	master.HandeFunc("RepeatNew", CommandRepeatNew)
	master.HandeFunc("DeleteWord", CommandDeleteWord)
	master.Url = "https://api.telegram.org/bot1060785017:AAG7eJUSygisjIF_g97Dj5TKVzS-ct76su8/"
	master.Offset = 0
	master.OpenDb, err = createDB()
	if err != nil {
		return nil, err
	}
	return master, nil
}

func createDB() (*sql.DB, error) {
	database, err := sql.Open("sqlite3", "./words.db")
	if err != nil {
		fmt.Print("\033[1;34m CreateDB open\033[0m\n")
		return nil, err
	}
	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS words (id INTEGER PRIMARY KEY, name TEXT, word TEXT, transcription TEXT, translate TEXT, ok int)")
	if err != nil {
		fmt.Print("\033[1;34createDB words\033[0m\n")
		return nil, err
	}
	_, err = statement.Exec()
	if err != nil {
		fmt.Print("\033[1;34createDB open\033[0m\n")
		return nil, err
	}
	statement, err = database.Prepare("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, name TEXT, chat_id INT)")
	if err != nil {
		fmt.Print("\033[1;34mCreate DB Users\033[0m\n")
		return nil, err
	}
	_, err = statement.Exec()
	if err != nil {
		fmt.Print("\033[1;34mCreate DB Users\033[0m\n")
		return nil, err
	}
	return database, nil
}
