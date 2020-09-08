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

func GetTranslate(db botStruct.RequestDb) *string {
	rows, err := db.Db.Query("select translate from words where name = ? and word = ?", db.Name, db.Word)
	if err != nil {
		log.Error.Println(err.Error())
		return nil
	}
	transaltes := new(string)
	var translate string;
	for rows.Next() {
		rows.Scan(&translate)
		*transaltes += " " + translate
	}
	if *transaltes == "" {
		return nil
	}
	*transaltes += " "
	return transaltes
}

func GetWord(db botStruct.RequestDb) bool {
	rows, err := db.Db.Query("SELECT word translate FROM words WHERE word = ? AND translate = ? AND name = ?", db.Word, db.Translate, db.Name)
	if err != nil {
		log.Error.Println(err.Error())
		return false
	}
	if rows.Next() {
		return false
	}
	rows.Close()
	return true
}

func DeleteWord(name string, word string, translate string, db *sql.DB) bool {
	stmt, err := db.Prepare("DELETE FROM words WHERE name = ? and word = ? AND translate = ?")
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

func AddWord(db botStruct.RequestDb) bool {
	stmt, err := db.Db.Prepare("INSERT INTO words (name, word, translate, ok) VALUES(?, ?, ?, ?)")
	if err != nil {
		log.Error.Println(err.Error())
		return false
	}
	_, err = stmt.Exec(db.Name, db.Word, db.Translate, 0)
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
