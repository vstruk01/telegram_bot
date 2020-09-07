package Init

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/vstruk01/telegram_bot/internal/commands"
	"github.com/vstruk01/telegram_bot/internal/sends"
	botStruct "github.com/vstruk01/telegram_bot/internal/struct"
)

var (
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func initLog() {
	file, err := os.OpenFile("info.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	Info = log.New(file, "\033[1;34mINFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(file, "\033[1;33mWARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(file, "\033[1;32mERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func InitAll() (*botStruct.Master, error) {
	file, err := os.OpenFile("info.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer file.Close()

	Error.SetOutput(file)
	Warning.SetOutput(file)
	Info.SetOutput(file)
	initLog()

	master := new(botStruct.Master)

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
		fmt.Print("\033[1;32mError WordKnow = ", err.Error(), "\033[0m\n")
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
		fmt.Println("id = ", id)
	}

	return master, nil
}

func CheckErr(err error) bool {
	if err != nil {
		Error.Println(err.Error())
		return true
	}
	return false
}

func createDB() (*sql.DB, error) {
	database, err := sql.Open("sqlite3", "./words.db")
	if CheckErr(err) {
		return nil, err
	}
	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS words (id INTEGER PRIMARY KEY, name TEXT, word TEXT, transcription TEXT, translate TEXT, ok int)")
	if CheckErr(err) {
		return nil, err
	}
	_, err = statement.Exec()
	if CheckErr(err) {
		return nil, err
	}
	statement, err = database.Prepare("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, name TEXT, chat_id INT)")
	if CheckErr(err) {
		return nil, err
	}
	_, err = statement.Exec()
	if CheckErr(err) {
		return nil, err
	}
	return database, nil
}
