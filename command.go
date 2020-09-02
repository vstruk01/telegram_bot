package main

import (
	"database/sql"

	// "fmt"

	// "strconv"

	// "os"

	_ "github.com/mattn/go-sqlite3"
)

func listNew(name string, chat_id int) error {
	database, err := sql.Open("sqlite3", "./words.db")
	if err != nil {
		return err
	}
	rows, err := database.Query("select word, translate from words where name = ? and ok = 0", name)
	if err != nil {
		return err
	}
	sendWords(rows, chat_id)
	return nil
}

func listKnow(name string, chat_id int) error {
	database, err := sql.Open("sqlite3", "./words.db")
	if err != nil {
		return err
	}
	rows, err := database.Query("select word, translate from words where name = ? and ok > 0", name)
	if err != nil {
		return err
	}
	sendWords(rows, chat_id)
	return nil
}
