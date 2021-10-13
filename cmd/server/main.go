package main

import (
	"fmt"
	"telegram_bot/config"
)

func main() {
	conf := config.GetConfig("config.json")

	fmt.Println(conf)
}
