package main

import (
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	hi "github.com/vstruk01/telegram_bot/internal/manager"
)

func main() {
	master, err := hi.InitAll()
	if err != nil {
		fmt.Println("\033[1;32mError = ", err.Error(), "\033[0m")
		return
	}

	for {
		err = hi.GetUpdate(master)
		if err != nil {
			fmt.Println("\033[1;32mError = ", err.Error(), "\033[0m")
		}
	}
}
