package workdb

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	log "github.com/vstruk01/telegram_bot/internal/Logger"
	botStruct "github.com/vstruk01/telegram_bot/internal/struct"
)

func GetUsers() bool {
	return true
}

func GetWordsNew() bool {
	return true
}

func GetWordsKnow() bool {
	return true
}

func GetWords() bool {
	return true
}

func DeleteUser(name string, word string, translate string, db *sql.DB) bool {
	stmt, err := db.Prepare("DELETE FROM words WHERE name = ? and word = ? and translate = ?")
	if err != nil {
		log.Error.Println(err.Error())
		return false
	}
	_, err = stmt.Exec(name, word, translate)
	if err != nil {
		log.Error.Println(err.Error())
		return false
	}
	return true
}

func AddUser() bool {
	return true
}

func CheckUser() bool {
	return true
}

func CheckWord() bool {
	return true
}