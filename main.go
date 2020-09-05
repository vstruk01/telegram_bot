package main

import (
	"fmt"
	"github.com/vstruk01/telegram_bot/internal/manager"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	master, err := InitAll()
	if err != nil {
		fmt.Println("\033[1;32mError = ", err.Error(), "\033[0m")
		return
	}

	for {
		err = GetUpdate(master)
		if err != nil {
			fmt.Println("\033[1;32mError = ", err.Error(), "\033[0m")
		}
	}
}
