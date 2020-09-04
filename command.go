package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)


func CommandWordKnow(r Request) {
	err := SendMessage("Enter Word Please", r.Chat_id)
	word := <- r.C

	
	if err != nil {
		r.ErrC <- err
	}
	r.ErrC <- nil
}

func CommandListNew(r Request) {
	database, err := sql.Open("sqlite3", "./words.db")
	if err != nil {
		fmt.Print("\033[1;34mlistNew\033[0m\n")
		r.ErrC <- err
	}
	rows, err := database.Query("select word, translate from words where name = ? and ok = 0", r.Name)
	if err != nil {
		fmt.Print("\033[1;34mlistNew\033[0m\n")
		r.ErrC <- err
	}
	SendWords(rows, r.Chat_id)
	fmt.Print("\033[1;34mlistNew Ok\033[0m\n")
	r.ErrC <- nil
}

func CommandListKnow(r Request) {
	database, err := sql.Open("sqlite3", "./words.db")
	if err != nil {
		fmt.Print("\033[1;34mlistKnow\033[0m\n")
		r.ErrC <- err
	}
	rows, err := database.Query("select word, translate from words where name = ? and ok > 0", r.Name)
	if err != nil {
		fmt.Print("\033[1;34mlistKnow\033[0m\n")
		r.ErrC <- err
	}
	SendWords(rows, r.Chat_id)
	fmt.Print("\033[1;34mlistKnow Ok\033[0m\n")
	r.ErrC <- nil
}

func InsertWord(name string, words []string) error {
	var old_words string

	database, err := sql.Open("sqlite3", "./words.db")
	if err != nil {
		fmt.Print("\033[1;34mInsert Word\033[0m\n")
		return err
	}
	rows, err := database.Query("select translate from words WHERE name = ? and word = ?", name, words[0])
	defer rows.Close()
	if err != nil {
		fmt.Print("\033[1;34mInsert Word\033[0m\n")
		return err
	}

	rows.Next()
	err = rows.Scan(&old_words)
	if err != nil {
		statement, _ := database.Prepare("insert into words (name, word, translate, ok)values(?, ?, ?, ?)")
		statement.Exec(name, words[0], ","+words[1]+",", 0)
	} else {
		new_word := old_words + words[1] + ","
		_, err := database.Exec("update words set translate = ?1 where name = ?2 and word = ?3", new_word, name, words[0])
		if err != nil {
			fmt.Print("\033[1;34mInsert Word\033[0m\n")
			return err
		}
	}
	fmt.Print("\033[1;34mInsert Word Ok\033[0m\n")
	return nil
}
