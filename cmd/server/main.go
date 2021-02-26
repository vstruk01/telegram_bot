package main

import (
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	Init "github.com/vstruk01/telegram_bot/init"
	"github.com/vstruk01/telegram_bot/internal/manager"
)

var (
	Url   string = "https://api.telegram.org/bot"                   // * url of telegram
	Token string = "1060785017:AAG7eJUSygisjIF_g97Dj5TKVzS-ct76su8" // * your token of telegram bot
)

func main() {
	master, err := Init.InitAll()

	if err != nil {
		fmt.Println("\033[1;32mError = ", err.Error(), "\033[0m")
		return
	}
	manager.GetUpdate(master)
}
