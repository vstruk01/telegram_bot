package storage

import (
	"database/sql"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type Sqlite3 struct {
	Connect *sql.DB
}

// connecting with sqlite database and create tables
func (s *Sqlite3) Init(dbName string) error {
	// default dbName ./info/words.db
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return err
	}
	s.Connect = db
	return s.InitTables()
}

// * work with users
func (s *Sqlite3) GetUsersIDs() (*[]int, error) {
	id_s := new([]int)
	rows, err := s.Connect.Query("SELECT chat_id FROM users")
	if err != nil {
		return nil, err
	}

	var id int
	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		*id_s = append(*id_s, id)
	}
	return id_s, nil
}

// to do
func (s *Sqlite3) CheckUser(name string, id int) error {
	var n string

	rows, err := s.Connect.Query("SELECT name FROM users WHERE name = ? AND chat_id = ?", name, id)
	if err != nil {
		return err
	}
	if rows.Next() {
		err = rows.Scan(&n)
		err = rows.Close()
		if err != nil {
			return err
		}
	} else {
		err = s.AddUser(name, id)
		if err != nil {
			return err
		}

		// ! to transfer manager
		// Ch := new(botStruct.Channels)
		// Ch.C = make(chan string, 1)
		// Ch.Done = make(chan bool, 1)
		// master.Routines[id] = Ch
		// // sends.SetButton(id)
	}
	return nil
}

func (s *Sqlite3) AddUser(name string, id int) error {
	statement, err := s.Connect.Prepare("INSERT INTO users (name, chat_id) VALUES(?, ?)")
	if err != nil {
		return err
	}
	_, err = statement.Exec(name, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *Sqlite3) GetWordsNew(name string) (*map[string]string, error) {
	rows, err := s.Connect.Query("SELECT word, translate FROM words WHERE name = ? AND ok = 0", name)
	if err != nil {
		return nil, err
	}
	return s.GetWords(rows), nil
}

func (s *Sqlite3) GetWordsKnow(name string) (*map[string]string, error) {
	rows, err := s.Connect.Query("SELECT word, translate FROM words WHERE name = ? AND ok > 0", name)
	if err != nil {
		return nil, err
	}
	return s.GetWords(rows), nil
}

func (s *Sqlite3) GetWords(rows *sql.Rows) *map[string]string {
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

func (s *Sqlite3) GetTranslate(name, word string) (*string, error) {
	rows, err := s.Connect.Query("SELECT translate FROM words WHERE name = ? AND word = ?", name, strings.TrimSpace(strings.ToLower(word)))
	if err != nil {
		return nil, err
	}
	transaltes := new(string)
	var translate string
	for rows.Next() {
		rows.Scan(&translate)
		*transaltes += " " + translate
	}
	rows.Close()
	if *transaltes == "" {
		return nil, err
	}
	*transaltes += " "
	return transaltes, nil
}

// return error if error exists or error ()
func (s *Sqlite3) CheckWord(word, translate, name string) (bool, error) {
	rows, err := s.Connect.Query("SELECT word translate FROM words WHERE word = ? AND translate = ? AND name = ?", word, translate, name)
	if err != nil {
		return false, err
	}
	if !rows.Next() {
		return false, nil
	}
	rows.Close()
	return true, nil
}

func (s *Sqlite3) DeleteWord(name, word, translate string) error {
	stmt, err := s.Connect.Prepare("DELETE FROM words WHERE name = ? and word = ? AND translate = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(name, word, translate)
	if err != nil {
		return err
	}
	return nil
}

func (s *Sqlite3) AddWord(word, translate, name string) error {
	word = strings.TrimSpace(strings.ToLower(word))
	translate = strings.TrimSpace(strings.ToLower(translate))

	stmt, err := s.Connect.Prepare("INSERT INTO words (name, word, translate, ok) VALUES(?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(name, word, translate, 0)
	if err != nil {
		return err
	}
	return nil
}

func (s *Sqlite3) UpdateWordKnow(name, word, translate string) error {
	_, err := s.Connect.Exec("UPDATE words SET ok = 1 WHERE name = ? AND word = ? AND translate = ?", name, word, translate)
	if err != nil {
		return err
	}
	return nil
}

func (s *Sqlite3) InitTables() error {
	statement, err := s.Connect.Prepare("CREATE TABLE IF NOT EXISTS words (id INTEGER PRIMARY KEY, name TEXT, word TEXT, transcription TEXT, translate TEXT, ok int)")
	if err != nil {
		return err
	}
	_, err = statement.Exec()
	if err != nil {
		return err
	}
	statement, err = s.Connect.Prepare("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, name TEXT, chat_id INT)")
	if err != nil {
		return err
	}
	_, err = statement.Exec()
	if err != nil {
		return err
	}
	return nil
}
