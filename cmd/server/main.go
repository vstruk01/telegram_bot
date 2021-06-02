package main

import (
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	Init "github.com/vstruk01/telegram_bot/init"
	log "github.com/vstruk01/telegram_bot/internal/Logger"
	"github.com/vstruk01/telegram_bot/internal/commands"
	"github.com/vstruk01/telegram_bot/internal/manager"

	// "github.com/vstruk01/telegram_bot/internal/sends"
	botStruct "github.com/vstruk01/telegram_bot/internal/struct"
	db "github.com/vstruk01/telegram_bot/internal/workdb"
)

var (
	Url   string = "https://api.telegram.org/bot"                   // * url of telegram
	Token string = "1060785017:AAG7eJUSygisjIF_g97Dj5TKVzS-ct76su8" // * your token of telegram bot
)

func init() {
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
		var ch m.Channels
		ch.C = make(chan string, 1)
		ch.Done = make(chan bool, 1)
		master.Routines[v] = &ch
	}
	return master, nil
}

func main() {
	master, err := Init.InitAll()

	if err != nil {
		fmt.Println("\033[1;32mError = ", err.Error(), "\033[0m")
		return
	}
	manager.GetUpdate(master)
}
