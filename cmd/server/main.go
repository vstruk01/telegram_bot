package main

import (
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	Init "github.com/vstruk01/telegram_bot/init"
	"github.com/vstruk01/telegram_bot/internal/manager"
)

func main() {
	master, err := Init.InitAll()
	if err != nil {
		fmt.Println("\033[1;32mError = ", err.Error(), "\033[0m")
		return
	}
	manager.GetUpdate(master)
}
