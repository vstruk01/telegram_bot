package Init

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	log "github.com/vstruk01/telegram_bot/internal/Logger"
	"github.com/vstruk01/telegram_bot/internal/commands"
	"github.com/vstruk01/telegram_bot/internal/sends"
	botStruct "github.com/vstruk01/telegram_bot/internal/struct"
)

func InitAll() (*botStruct.Master, error) {
	master := new(botStruct.Master)
	var err error

	// * create map of handlers
	master.Commands = make(map[string]func(botStruct.Request))
	master.HandeFunc("/start", commands.CommandStart)
	master.HandeFunc("RepeatKnow", commands.CommandRepeatKnow)
	master.HandeFunc("ListKnow", commands.CommandListKnow)
	master.HandeFunc("WordKnow", commands.CommandWordKnow)
	master.HandeFunc("WordNew", commands.CommandWordNew)
	master.HandeFunc("ListNew", commands.CommandListNew)
	master.HandeFunc("RepeatNew", commands.CommandRepeatNew)
	master.HandeFunc("DeleteWord", commands.CommandDeleteWord)
	master.HandeFunc("AddWord", commands.CommandDeleteWord)

	// * initialization other veriables
	master.Offset = 0
	master.OpenDb, err = createDB()
	if err != nil {
		return nil, err
	}

	// * create map of chans for goroutines
	master.Routines = make(map[int]botStruct.Channels)
	rows, err := master.OpenDb.Query("select chat_id from users")
	if err != nil {
		return nil, err
	}
	var id int
	for rows.Next() {
		var ch botStruct.Channels
		rows.Scan(&id)
		sends.SetButton(id)
		ch.C = make(chan string, 1)
		ch.Done = make(chan bool, 1)
		master.Routines[id] = ch
	}

	return master, nil
}

func createDB() (*sql.DB, error) {
	database, err := sql.Open("sqlite3", "./words.db")
	if log.CheckErr(err) {
		return nil, err
	}
	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS words (id INTEGER PRIMARY KEY, name TEXT, word TEXT, transcription TEXT, translate TEXT, ok int)")
	if log.CheckErr(err) {
		return nil, err
	}
	_, err = statement.Exec()
	if log.CheckErr(err) {
		return nil, err
	}
	statement, err = database.Prepare("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, name TEXT, chat_id INT)")
	if log.CheckErr(err) {
		return nil, err
	}
	_, err = statement.Exec()
	if log.CheckErr(err) {
		return nil, err
	}
	return database, nil
}
