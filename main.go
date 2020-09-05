package main

import (
	"database/sql"
	"fmt"
	_ "github.com/vstruk01/telegram_bot/tree/master/internal/manager"
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
