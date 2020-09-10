package workdb

import (
	"database/sql"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	log "github.com/vstruk01/telegram_bot/internal/Logger"
	sends "github.com/vstruk01/telegram_bot/internal/sends"
	botStruct "github.com/vstruk01/telegram_bot/internal/struct"
)

// * work with users * //
func GetUsersID(db *sql.DB) (*[]int, error) {
	id_s := new([]int)
	rows, err := db.Query("SELECT chat_id FROM users")
	if err != nil {
		return nil, err
	}
	var id int
	for rows.Next() {
		rows.Scan(&id)
		*id_s = append(*id_s, id)
	}
	return id_s, nil
}

func CheckUser(master *botStruct.Master, name string, id int) error {
	var n string

	rows, err := master.OpenDb.Query("SELECT name FROM users WHERE name = ? AND chat_id = ?", name, id)
	if err != nil {
		log.Error.Println(err.Error())
		return err
	}
	if rows.Next() {
		err = rows.Scan(&n)
		err = rows.Close()
		if err != nil {
			log.Error.Println(err.Error())
			return err
		}
	} else {
		err = AddUser(master.OpenDb, name, id)
		if err != nil {
			log.Error.Println(err.Error())
			return err
		}
		Ch := new(botStruct.Channels)
		Ch.C = make(chan string, 1)
		Ch.Done = make(chan bool, 1)
		master.Routines[id] = Ch
		sends.SetButton(id)
	}
	return nil
}

func AddUser(db *sql.DB, name string, id int) error {
	statement, err := db.Prepare("INSERT INTO users (name, chat_id) VALUES(?, ?)")
	if err != nil {
		log.Error.Println(err.Error())
		return err
	}
	_, err = statement.Exec(name, id)
	if err != nil {
		log.Error.Println(err.Error())
		return err
	}
	return nil
}

// * work with words * //

func GetWordsNew(r botStruct.Request) *map[string]string {
	rows, err := r.OpenDb.Query("SELECT word, translate FROM words WHERE name = ? AND ok = 0", r.Name)
	if err != nil {
		log.Error.Println(err.Error())
		return nil
	}
	return GetWords(r, rows)
}

func GetWordsKnow(r botStruct.Request) *map[string]string {
	rows, err := r.OpenDb.Query("SELECT word, translate FROM words WHERE name = ? AND ok > 0", r.Name)
	if err != nil {
		log.Error.Println(err.Error())
		return nil
	}
	return GetWords(r, rows)
}

func MapWordsInStringWords(m_words *map[string]string) *string {
	Word := new(string)

	for k, v := range *m_words {
		*Word += k + " -> " + v + "\n"
	}
	if *Word == "" {
		*Word += "empty :("
	}
	return Word
}

func GetWords(r botStruct.Request, rows *sql.Rows) *map[string]string {
	m_words := new(map[string]string)
	*m_words = make(map[string]string)
	var word, translate string
	for rows.Next() {
		rows.Scan(&word, &translate)
		_, ok := (*m_words)[word]
		if ok {
			(*m_words)[word] += translate + " "
		} else {
			(*m_words)[word] = " " + translate + " "
		}
	}
	rows.Close()
	return m_words
}

func GetTranslate(db botStruct.RequestDb) *string {
	rows, err := db.Db.Query("SELECT translate FROM words WHERE name = ? AND word = ?", db.Name, strings.TrimSpace(strings.ToLower(db.Word)))
	if err != nil {
		log.Error.Println(err.Error())
		return nil
	}
	transaltes := new(string)
	var translate string
	for rows.Next() {
		rows.Scan(&translate)
		*transaltes += " " + translate
	}
	rows.Close()
	if *transaltes == "" {
		return nil
	}
	*transaltes += " "
	return transaltes
}

func CheckWord(db botStruct.RequestDb) bool {
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
	word := strings.TrimSpace(strings.ToLower(db.Word))
	translate := strings.TrimSpace(strings.ToLower(db.Translate))

	stmt, err := db.Db.Prepare("INSERT INTO words (name, word, translate, ok) VALUES(?, ?, ?, ?)")
	if err != nil {
		log.Error.Println(err.Error())
		return false
	}
	_, err = stmt.Exec(db.Name, word, translate, 0)
	if err != nil {
		log.Error.Println(err.Error())
		return false
	}
	return true
}

func UpdateWordKnow(name, word, translate string, db *sql.DB) {
	_, err := db.Exec("UPDATE words SET ok = 1 WHERE name = ? AND word = ? AND translate = ?", name, word, translate)
	if err != nil {
		log.Error.Println(err.Error())
	}
}
