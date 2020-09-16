package Init

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	log "github.com/vstruk01/telegram_bot/internal/Logger"
	"github.com/vstruk01/telegram_bot/internal/commands"
	// "github.com/vstruk01/telegram_bot/internal/sends"
	db "github.com/vstruk01/telegram_bot/internal/workdb"
	botStruct "github.com/vstruk01/telegram_bot/internal/struct"
)

func InitAll() (*botStruct.Master, error) {
	log.InitLog()

	master := new(botStruct.Master)
	var err error

	// * create map of handlers
	master.Commands = make(map[string]func(botStruct.Request))
	master.HandlerFunc("/start", commands.CommandStart)
	master.HandlerFunc("/help", commands.CommandHelp)
	master.HandlerFunc("/add_word", commands.CommandAddWord)
	master.HandlerFunc("/repeat_know", commands.CommandRepeatKnow)
	master.HandlerFunc("/list_know", commands.CommandListKnow)
	master.HandlerFunc("/word_know", commands.CommandWordKnow)
	master.HandlerFunc("/list_new", commands.CommandListNew)
	master.HandlerFunc("/repeat_new", commands.CommandRepeatNew)
	master.HandlerFunc("/delete_word", commands.CommandDeleteWord)
	// * handlers for buttom
	master.HandlerFunc("Repeat Know", commands.CommandRepeatKnow)
	master.HandlerFunc("List Know", commands.CommandListKnow)
	master.HandlerFunc("Word Know", commands.CommandWordKnow)
	master.HandlerFunc("List New", commands.CommandListNew)
	master.HandlerFunc("Repeat New", commands.CommandRepeatNew)
	master.HandlerFunc("Delete Word", commands.CommandDeleteWord)
	master.HandlerFunc("Add Word", commands.CommandAddWord)

	// * initialization other veriables
	master.Offset = 0
	master.OpenDb, err = createDB()
	if err != nil {
		return nil, err
	}

	// * create map of chans for goroutines
	master.Routines = make(map[int]*botStruct.Channels)
	users_id, err := db.GetUsersID(master.OpenDb)
	if err != nil {
		return nil, err
	}
	for _, v := range *users_id {
		var ch botStruct.Channels
		ch.C = make(chan string, 1)
		ch.Done = make(chan bool, 1)
		master.Routines[v] = &ch
	}
	return master, nil
}

func createDB() (*sql.DB, error) {
	database, err := sql.Open("sqlite3", "./info/words.db")
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
